package apierrors

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)


type APIError struct {
	Code    int
	Message any
}


func (e *APIError) Error() string {
    switch msg := e.Message.(type) {
    case string:
        return msg
    case error:
        return msg.Error()
    default:
        return fmt.Sprintf("%v", msg)
    }
}


var (
	ErrUserNotFound = APIError{Code: http.StatusNotFound, Message: "user not found"}
	ErrUserAlreadyExist = APIError{Code: http.StatusConflict, Message: "user already exists"}
	ErrInternalServerError = APIError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrInvalidRequestBody = APIError{Code: http.StatusBadRequest, Message: "invalid request body"}
	ErrEncodingError = APIError{Code: http.StatusInternalServerError, Message: "encoding error"}
	ErrValidationError = APIError{Code: http.StatusBadRequest, Message: "validation error"}
)



func NewValidationError(errs validator.ValidationErrors) *APIError {
	validationErrors := make(map[string]string)
	for _, err := range errs {
		validationErrors[err.Field()] = err.Tag()
	}
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: validationErrors,
	}
}