package repositories

import (
	"context"
	"fmt"
	"golang/internal/infrastructure/database/models"

	"github.com/jackc/pgx/v5/pgxpool"
)


type UserRepository struct {
	DB *pgxpool.Pool
}


func (repo *UserRepository) GetUserById(ctx context.Context, userId int) (*models.User, error) {
	userRow := repo.DB.QueryRow(ctx, "SELECT * FROM user WHERE user.id = $1", userId)
	
	var user models.User
	
	err := userRow.Scan(&user)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	return &user, nil
}