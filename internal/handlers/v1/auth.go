package handlers

import (
	"encoding/json"
	"golang/internal/core/services"
	deps "golang/internal/handlers/dependencies"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"net/http"
	"strings"
)


type AuthHandler struct {
	Service *services.AuthService
}


func (handler *AuthHandler) RegisterUser(response http.ResponseWriter, request *http.Request) {
	user, serviceErr := handler.Service.RegisterUser(request.Context(), request.Body)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(response).Encode(user); err != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrEncodingError)
	}
}


func (handler *AuthHandler) GetCurrentUser(response http.ResponseWriter, request *http.Request) {
	authHeader := request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := handler.Service.ValidateToken(request.Context(), tokenString)

	response.Header().Set("Content-Type", "application/json")

	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}
	
	if err := json.NewEncoder(response).Encode(user); err != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrEncodingError)
	}
}


func (handler *AuthHandler) RefreshToken(response http.ResponseWriter, request *http.Request) {
	refreshToken := request.URL.Query().Get("refresh_token")
	user, err := handler.Service.RefreshToken(request.Context(), refreshToken)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusAccepted)

	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return 
	}

	if err := json.NewEncoder(response).Encode(user); err != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrEncodingError)
	}
}


func (handler *AuthHandler) LoginUser(response http.ResponseWriter, request *http.Request) {
	var userForm models.LoginUserModel
	err := json.NewDecoder(request.Body).Decode(&userForm)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	user, serviceErr := handler.Service.LoginUser(request.Context(), userForm)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(response).Encode(user); err != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrEncodingError)
	}
}


func (handler *AuthHandler) SetupRoutes(server *http.ServeMux, baseUrl string, protected *deps.AuthDependency) {
	server.HandleFunc(baseUrl+"/auth/register", handler.RegisterUser)
	server.HandleFunc(baseUrl+"/auth/login", handler.LoginUser)
	server.HandleFunc(baseUrl+"/auth/refresh", handler.RefreshToken)
	server.HandleFunc(baseUrl+"/auth/current", handler.GetCurrentUser)
}