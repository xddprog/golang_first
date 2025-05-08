package services

import (
	"context"
	"encoding/json"
	"golang/internal/core/repositories"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"golang/internal/utils"
	"io"
)


type UserService struct {
	Repository *repositories.UserRepository
}



func (s *UserService) GetUserById(ctx context.Context, userId int) (*models.UserModel, *apierrors.APIError) {
	user, err := s.Repository.GetUserById(ctx, userId)
	if err != nil {
		return nil, apierrors.CheckDBError(err) 
	}
	return user, nil
}


func (s *UserService) UpdateUser(ctx context.Context, userId int, userForm io.ReadCloser) (*models.UserModel, *apierrors.APIError) {
	var userFormEncoded models.UpdateUserModel

	if err := json.NewDecoder(userForm).Decode(&userFormEncoded); err != nil {
		return nil, &apierrors.ErrInvalidRequestBody
	}

	if err := utils.ValidateForm(userFormEncoded); err != nil {
		return nil, err
	}
	
	user, err := s.Repository.UpdateUser(ctx, userId, userFormEncoded)
	if err != nil {
		return nil, apierrors.CheckDBError(err) 
	}
	return user, nil
}


func (s *UserService) DeleteUser(ctx context.Context, userId int) *apierrors.APIError {
	err := s.Repository.DeleteUser(ctx, userId)
	if err != nil {
		return apierrors.CheckDBError(err) 
	}
	return nil
}