package utils

import (
	"encoding/json"
	"fmt"
	apierrors "golang/internal/infrastructure/errors"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx"
)


func GetSetParams(form any) (string, []any) {
	var args []interface{}
	var clauses []string
	paramIndex := 1

	values := reflect.ValueOf(form)
    types := values.Type()

    for i := 0; i < values.NumField(); i++ {
		value := values.Field(i)

		if !value.IsNil() {
			field := types.Field(i).Tag.Get("db")
			safeField := pgx.Identifier{field}.Sanitize()
			clauses = append(clauses, fmt.Sprintf("%s = $%d", safeField, paramIndex))
			args = append(args, value.Interface())
			paramIndex++
		}
    }
	clausesStr := strings.Join(clauses, ", ")
	return clausesStr, args
}


func ValidateForm(form any) *apierrors.APIError {
	validate := validator.New()
    if err := validate.Struct(form); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            return apierrors.NewValidationError(validationErrors)
        }
        return &apierrors.ErrValidationError
    }
    return nil
}


func RandSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}


func GetLimitAndOffset(request *http.Request) (int, int) {
	limit := 10
	offset := 0

	if limitStr := request.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := request.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}

	return limit, offset
}


func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if data == nil {
        return nil
    }
    if err := json.NewEncoder(w).Encode(data); err != nil {
        apierrors.WriteHTTPError(w, apierrors.ErrEncodingError)
        return err
    }
    return nil
}