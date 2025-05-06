package services

import (
	"context"
	"golang/internal/core/repositories"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
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
