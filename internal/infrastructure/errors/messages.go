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
	ErrItemNotFound = func (itemName string) *APIError {
		return &APIError{Code: http.StatusNotFound, Message: fmt.Sprintf("%s not found", itemName)}
	}
	ErrUserAlreadyExist = APIError{Code: http.StatusConflict, Message: "user already exists"}
	ErrInternalServerError = APIError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrInvalidRequestBody = APIError{Code: http.StatusBadRequest, Message: "invalid request body"}
	ErrEncodingError = APIError{Code: http.StatusInternalServerError, Message: "encoding error"}
	ErrValidationError = APIError{Code: http.StatusBadRequest, Message: "validation error"}
	ErrInvalidToken = APIError{Code: http.StatusUnauthorized, Message: "invalid token"}
	ErrInvaliLoginData = APIError{Code: http.StatusUnauthorized, Message: "invalid login data"}
	ErrDocumentAccessDenied = APIError{Code: http.StatusForbidden, Message: "access to document denied"}
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