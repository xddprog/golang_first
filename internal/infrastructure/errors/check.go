package apierrors

import "github.com/jackc/pgx/v5"

func CheckDBError(err error) *APIError {
	switch err {
	case pgx.ErrNoRows:
		return &ErrUserNotFound
	default:
		return &ErrInternalServerError
	}
}