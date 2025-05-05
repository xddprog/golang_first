package handlers

import (
	"encoding/json"
	"golang/internal/core/services"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"net/http"
)


type AuthHandler struct {
	Service *services.AuthService
}



func (handler *UserHandler) CreateUser(response http.ResponseWriter, request *http.Request) {
	var userForm models.CreateUserModel
	err := json.NewDecoder(request.Body).Decode(&userForm)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	user, serviceErr := handler.Service.CreateUser(request.Context(), userForm)
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
