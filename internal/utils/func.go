package utils

import (
	"fmt"
	apierrors "golang/internal/infrastructure/errors"
	"reflect"
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