package handlers

import (
	"golang/internal/core/services"
	"golang/internal/handlers/dependencies"
	"golang/internal/infrastructure/errors"
	"golang/internal/utils"
	"net/http"
	"strings"
)


type AuthHandler struct {
	Service *services.AuthService
}


func (handler *AuthHandler) RegisterUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	user, serviceErr := handler.Service.RegisterUser(request.Context(), request.Body)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	utils.WriteJSONResponse(response, http.StatusCreated, user)
}


func (handler *AuthHandler) GetCurrentUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	authHeader := request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := handler.Service.ValidateToken(request.Context(), tokenString)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}
	
	utils.WriteJSONResponse(response, http.StatusOK, user)
}


func (handler *AuthHandler) RefreshToken(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	refreshToken := request.URL.Query().Get("refresh_token")
	user, err := handler.Service.RefreshToken(request.Context(), refreshToken)

	response.WriteHeader(http.StatusAccepted)

	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return 
	}

	utils.WriteJSONResponse(response, http.StatusOK, user)
}


func (handler *AuthHandler) LoginUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	user, serviceErr := handler.Service.LoginUser(request.Context(), request.Body)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	utils.WriteJSONResponse(response, http.StatusOK, user)
}


func (handler *AuthHandler) SetupRoutes(server *http.ServeMux, baseUrl string, protected *deps.AuthDependency) {
	server.HandleFunc(baseUrl+"/auth/register", handler.RegisterUser)
	server.HandleFunc(baseUrl+"/auth/login", handler.LoginUser)
	server.HandleFunc(baseUrl+"/auth/refresh", handler.RefreshToken)
	server.HandleFunc(baseUrl+"/auth/current", handler.GetCurrentUser)
}