package apierrors

import (
	"log"

	"github.com/jackc/pgx/v5"
)

func CheckDBError(err error, itemName string) *APIError {
	switch err {
	case pgx.ErrNoRows:
		return ErrItemNotFound(itemName)
	default:
		log.Printf("Internal Server Error: %v", err)
		return &ErrInternalServerError
	}
}