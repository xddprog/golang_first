package services

import (
	"context"
	"golang/internal/core/repositories"
	"golang/internal/infrastructure/config"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"golang/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)


type AuthService struct {
	Config *config.JwtConfig
	Repository *repositories.UserRepository
}


func (s *AuthService) HashPassword(password string) (string, *apierrors.APIError) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", &apierrors.ErrInternalServerError
    }
    return string(hashedPassword), nil
}


func (s *AuthService) CheckPassword(password, hashedPassword string) *apierrors.APIError {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return &apierrors.ErrInvaliLoginData
	}
	return nil
}


func (s *AuthService) RegisterUser(
	ctx context.Context, 
	userForm models.CreateUserModel,
) (*models.AuthResponseModel, *apierrors.APIError) {
	checkExist, err := s.Repository.GetUserByEmail(ctx, userForm.Email)
	if err != nil && err != pgx.ErrNoRows {
		return nil, apierrors.CheckDBError(err)
	}

	if checkExist != nil {
		return nil, &apierrors.ErrUserAlreadyExist
	}

	validate := validator.New()
	if err := validate.Struct(userForm); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return nil, apierrors.NewValidationError(validationErrors)
		}
		return nil, &apierrors.ErrValidationError
	}

	hashedPassword, hashErr := s.HashPassword(userForm.Password)
	if err != nil {
		return nil, hashErr
	}

	userForm.Password = hashedPassword

	user, err := s.Repository.CreateUser(ctx, userForm)
	if err != nil {
		return nil, apierrors.CheckDBError(err)
	}

	tokenPair, err := s.createTokenPair(user.Id)
	return &models.AuthResponseModel{
        TokenPair: *tokenPair,
        User:      *user,
    }, nil
}


func (s *AuthService) createTokenPair(userId int) (*models.TokenPair, *apierrors.APIError) {
	accessToken, err := s.createToken(userId, utils.AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.createToken(userId, utils.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) createToken(userId int, tokenType string) (string, *apierrors.APIError) {
	var expiresAt time.Duration

	switch tokenType {
	case utils.AccessToken:
		expiresAt = s.Config.AccessTokenTime 
	case utils.RefreshToken:
		expiresAt = s.Config.RefreshTokenTime
	default:
		return "", &apierrors.ErrInternalServerError
	}

	token := jwt.NewWithClaims(s.Config.SigningMethod, jwt.MapClaims{
			"sub": userId,
			"exp": time.Now().Add(expiresAt).Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(s.Config.Secret))
	if err != nil {
		return err.Error(), &apierrors.ErrInternalServerError
	}
	return tokenString, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*models.UserModel, *apierrors.APIError) {
	if tokenString != "" {
		return nil, &apierrors.ErrInvalidToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != s.Config.SigningMethod {
			return nil, &apierrors.ErrInvalidToken
		}

		return s.Config.Secret, nil
	})

	if err != nil {
		return nil, &apierrors.ErrInternalServerError
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(float64)
		if !ok {
			return nil, &apierrors.ErrInvalidToken
		}

		user, err := s.Repository.GetUserById(ctx, int(userID))
		if err != nil {
			return nil, apierrors.CheckDBError(err)
		}
		return user, nil
	}
	return nil, &apierrors.ErrInternalServerError
}


func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponseModel, *apierrors.APIError) {
	user, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.createToken(user.Id, utils.AccessToken)
	if err != nil {
		return nil, err
	}
	return &models.AuthResponseModel{
        TokenPair: models.TokenPair{
			AccessToken: accessToken, 
			RefreshToken: refreshToken,
		},
        User:      *user,
    }, nil
}


func (s *AuthService) LoginUser(ctx context.Context, userForm models.LoginUserModel) (*models.AuthResponseModel, *apierrors.APIError) {
	user, err := s.Repository.GetUserByEmail(ctx, userForm.Email)
	if err != nil {
		return nil, apierrors.CheckDBError(err)
	}

	passErr := s.CheckPassword(userForm.Password, user.Password)
	if passErr != nil {
		return nil, passErr
	}

	tokenPair, tokenPairErr := s.createTokenPair(user.Id)
	if tokenPairErr != nil {
		return nil , tokenPairErr
	}
	return &models.AuthResponseModel{
		TokenPair: *tokenPair,
		User: *user,
	}, nil
}