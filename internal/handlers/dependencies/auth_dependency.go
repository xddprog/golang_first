package deps

import (
	"context"
	"golang/internal/core/services"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"net/http"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)


type AuthenticatedHandlerFunc func(w http.ResponseWriter, r *http.Request, user *models.UserModel)


type AuthenticatedSocketHandlerFunc func(s socketio.Conn) error


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


type (
	SocketHandler          		func(s socketio.Conn) error
	SocketEventHandler     		func(s socketio.Conn, data string)
	AuthedSocketHandler    		func(s socketio.Conn, user *models.BaseUserModel) error
	AuthedSocketEventHandler 	func(s socketio.Conn, data string, user *models.BaseUserModel)
)


func (d *AuthDependency) ProtectConnect(handler SocketHandler) SocketHandler {
	return func(s socketio.Conn) error {
		_, err := d.authenticateSocket(s)
		if err != nil {
			return err
		}
		return handler(s)
	}
}


func (d *AuthDependency) ProtectEvent(handler AuthedSocketEventHandler) SocketEventHandler {
	return func(s socketio.Conn, data string) {
		user, ok := s.Context().(*models.BaseUserModel)
		if !ok {
			s.Emit("error", "authentication required")
			return
		}
		handler(s, data, user)
	}
}


func (d *AuthDependency) authenticateSocket(s socketio.Conn) (*models.UserModel, *apierrors.APIError) {
	tokenString := strings.TrimPrefix(
		s.RemoteHeader().Get("Authorization"),
		"Bearer ",
	)

	if tokenString == "" {
		s.Emit("error", "missing auth token")
		return nil, &apierrors.ErrInvalidToken
	}

	user, err := d.Service.ValidateToken(context.Background(), tokenString)
	if err != nil {
		s.Emit("error", "invalid token")
		return nil, err
	}

	s.SetContext(user)
	return user, nil
}