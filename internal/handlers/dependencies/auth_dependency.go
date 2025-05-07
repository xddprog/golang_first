package deps

import (
	"golang/internal/core/services"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"net/http"
	"strings"
)


type AuthenticatedHandlerFunc func(w http.ResponseWriter, r *http.Request, user *models.UserModel)




type AuthDependency struct {
    Service *services.AuthService
}


func NewAuthDependency(authService *services.AuthService) *AuthDependency {
    return &AuthDependency{
        Service: authService,
    }
}


func (d *AuthDependency) Protected(handler AuthenticatedHandlerFunc) http.HandlerFunc {
    return func(response http.ResponseWriter, request *http.Request) {
        authHeader := request.Header.Get("Authorization")
        if authHeader == "" {
            apierrors.WriteHTTPError(response, &apierrors.ErrInvalidToken)
            return
        }
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == "" {
            apierrors.WriteHTTPError(response, &apierrors.ErrInvalidToken)
            return
        }

        user, err := d.Service.ValidateToken(request.Context(), tokenString)
        if err != nil {
            apierrors.WriteHTTPError(response, err)
            return
        }
        handler(response, request, user)
    }
}