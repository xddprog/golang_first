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
	userRow := repo.DB.QueryRow(ctx, "SELECT * FROM user WHERE email = $2", value)

	var user models.UserModel

	err := userRow.Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetUserById(ctx context.Context, userId int) (*models.UserModel, error) {
	userRow := repo.DB.QueryRow(ctx, "SELECT * FROM user WHERE id = $1", userId)

	var user models.UserModel

	err := userRow.Scan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, userForm models.RegisterUserModel) (*models.UserModel, error) {
	userRow := repo.DB.QueryRow(
		ctx,
		"INSERT INTO user (username, email, password) RETURNING username, email",
		userForm.Username, userForm.Email, userForm.Password,
	)

	var user models.UserModel

	err := userRow.Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func (repo *UserRepository) UpdateUser(ctx context.Context, userId int, userForm models.UpdateUserModel) (*models.UserModel, error) {
	userRow := repo.DB.QueryRow(
		ctx,
		"UPDATE user SET username = $1, email = $2 WHERE id = $3 RETURNING id, username, email",
		userForm.Username, userForm.Email, userId,
	)

	var user models.UserModel

	err := userRow.Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func (repo *UserRepository) DeleteUser(ctx context.Context, userId int) error {
	_, err := repo.DB.Exec(ctx, "DELETE FROM user WHERE id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}