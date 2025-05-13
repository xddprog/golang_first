package repositories

import (
	"context"
	"golang/internal/infrastructure/database/models"

	"github.com/jackc/pgx/v5/pgxpool"
)


type UserRepository struct {
	DB *pgxpool.Pool
}


func (repo *UserRepository) GetUserByEmail(ctx context.Context, value string) (*models.UserModel, error) {
	var user models.UserModel

	err := repo.DB.QueryRow(ctx, "SELECT id, username, email, password FROM users WHERE email = $1", value).Scan(
		&user.Id, &user.Username, &user.Email, &user.Password,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetUserById(ctx context.Context, userId int) (*models.BaseUserModel, error) {
	var user models.BaseUserModel

	err := repo.DB.QueryRow(ctx, "SELECT id, username, email FROM users WHERE id = $1", userId).Scan(
		&user.Id, &user.Username, &user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, userForm models.RegisterUserModel) (*models.BaseUserModel, error) {
	var user models.BaseUserModel
	
	err := repo.DB.QueryRow(
		ctx,
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email",
		userForm.Username, userForm.Email, userForm.Password,
	).Scan(&user.Id, &user.Username, &user.Email)

	if err != nil {
		return nil, err
	}
	return &user, nil
}


func (repo *UserRepository) UpdateUser(ctx context.Context, userId int, userForm models.UpdateUserModel) (*models.BaseUserModel, error) {
	var user models.BaseUserModel

	err := repo.DB.QueryRow(
		ctx,
		"UPDATE users SET username = $1, email = $2 WHERE id = $3 RETURNING id, username, email",
		userForm.Username, userForm.Email, userId,
	).Scan(&user.Id, &user.Username, &user.Email)

	if err != nil {
		return nil, err
	}
	return &user, nil
}


func (repo *UserRepository) DeleteUser(ctx context.Context, userId int) error {
	_, err := repo.DB.Exec(ctx, "DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}