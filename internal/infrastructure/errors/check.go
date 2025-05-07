package apierrors

import (
	"log"

	"github.com/jackc/pgx/v5"
)

func CheckDBError(err error) *APIError {
	switch err {
	case pgx.ErrNoRows:
		return &ErrUserNotFound
	default:
		log.Printf("Internal Server Error: %v", err)
	}
	return nil
}