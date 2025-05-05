package services

import (
	"context"
	"fmt"
	"golang/internal/core/repositories"
	"golang/internal/infrastructure/database/models"
)


type UserService struct {
	Repository *repositories.UserRepository
}


func (s *UserService) GetUserById(ctx context.Context, userId int) (*models.User, error) {
	user, err := s.Repository.GetUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("")
	}
	return user, nil
}