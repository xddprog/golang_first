package services

import (
	"context"
	"golang/internal/core/repositories"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)


type UserService struct {
	Repository *repositories.UserRepository
}


func checkDBError(err error) *apierrors.APIError {
	switch err {
	case pgx.ErrNoRows:
		return &apierrors.ErrUserNotFound
	default:
		return &apierrors.ErrInternalServerError
	}
}


func (s *UserService) GetUserById(ctx context.Context, userId int) (*models.UserModel, *apierrors.APIError) {
	user, err := s.Repository.GetUserById(ctx, userId)
	if err != nil {
		return nil, checkDBError(err) 
	}
	return user, nil
}


func (s *UserService) CreateUser(ctx context.Context, userForm models.CreateUserModel) (*models.UserModel, *apierrors.APIError) {
	checkExist, err := s.Repository.GetUserByEmail(ctx, userForm.Email)
	if err != nil && err != pgx.ErrNoRows {
		return nil, checkDBError(err)
	}

	if checkExist != nil {
		return nil, &apierrors.ErrUserAlreadyExist
	}

	validate := validator.New()
	if err :=  validate.Struct(userForm); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return nil, apierrors.NewValidationError(validationErrors)
		}
		return nil, &apierrors.ErrValidationError
	}

	user, err := s.Repository.CreateUser(ctx, userForm)
	if err != nil {
		return nil, checkDBError(err)
	}
	return user, nil
	
}