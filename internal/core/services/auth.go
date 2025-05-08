package services

import (
    "context"
    "encoding/json"
    "golang/internal/core/repositories"
    "golang/internal/infrastructure/config"
    "golang/internal/infrastructure/database/models"
    "golang/internal/infrastructure/errors"
    "golang/internal/utils"
    "io"
    "log"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/jackc/pgx/v5"
    "golang.org/x/crypto/bcrypt"
)


type AuthService struct {
    Config     *config.JwtConfig
    Repository *repositories.UserRepository
}


func (s *AuthService) HashPassword(password string) (string, *apierrors.APIError) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("[INTERNAL] Failed to generate password hash: %v", err)
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
    userForm io.ReadCloser,
) (*models.AuthResponseModel, *apierrors.APIError) {
    var userFormEncoded models.RegisterUserModel
    
    if err := json.NewDecoder(userForm).Decode(&userFormEncoded); err != nil {
        return nil, &apierrors.ErrInvalidRequestBody
    }

    checkExist, err := s.Repository.GetUserByEmail(ctx, userFormEncoded.Email)
    if err != nil && err != pgx.ErrNoRows {
        return nil, apierrors.CheckDBError(err)
    }
    if checkExist != nil {
        return nil, &apierrors.ErrUserAlreadyExist
    }

    if err := utils.ValidateForm(userFormEncoded); err != nil {
        return nil, err
    }

    hashedPassword, hashErr := s.HashPassword(userFormEncoded.Password)
    if hashErr != nil {
        return nil, hashErr
    }

    userFormEncoded.Password = hashedPassword

    user, err := s.Repository.CreateUser(ctx, userFormEncoded)
    if err != nil {
        return nil, apierrors.CheckDBError(err)
    }

    tokenPair, tokenPairErr := s.createTokenPair(user.Id)
    if tokenPairErr != nil {
        return nil, tokenPairErr
    }

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
    }

    token := jwt.NewWithClaims(s.Config.SigningMethod, jwt.MapClaims{
        "sub": userId,
        "exp": time.Now().Add(expiresAt).Unix(),
    })

    tokenString, err := token.SignedString([]byte(s.Config.Secret))
    if err != nil {
        log.Printf("[INTERNAL] Failed to sign JWT token: %v", err)
        return "", &apierrors.ErrInternalServerError
    }
    return tokenString, nil
}


func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*models.UserModel, *apierrors.APIError) {
    if tokenString == "" {
        return nil, &apierrors.ErrInvalidToken
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
        if token.Method != s.Config.SigningMethod {
            return nil, &apierrors.ErrInvalidToken
        }
        return []byte(s.Config.Secret), nil
    })

    if err != nil {
        log.Printf("Failed to parse token: %v", err)
        return nil, &apierrors.ErrInvalidToken
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

    return nil, &apierrors.ErrInvalidToken
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
            AccessToken:  accessToken,
            RefreshToken: refreshToken,
        },
        User: *user,
    }, nil
}


func (s *AuthService) LoginUser(ctx context.Context, userForm io.ReadCloser) (*models.AuthResponseModel, *apierrors.APIError) {
    var userFormEncoded models.LoginUserModel
    err := json.NewDecoder(userForm).Decode(&userFormEncoded)
    if err != nil {
        return nil, &apierrors.ErrInvalidRequestBody
    }

    if err := utils.ValidateForm(userFormEncoded); err != nil {
        return nil, err
    }

    user, err := s.Repository.GetUserByEmail(ctx, userFormEncoded.Email)
    if err != nil { 
        return nil, apierrors.CheckDBError(err)
    }

    passErr := s.CheckPassword(userFormEncoded.Password, user.Password)
    if passErr != nil {
        return nil, passErr
    }

    tokenPair, tokenPairErr := s.createTokenPair(user.Id)
    if tokenPairErr != nil {
        return nil, tokenPairErr
    }

    return &models.AuthResponseModel{
        TokenPair: *tokenPair,
        User:      *user,
    }, nil
}