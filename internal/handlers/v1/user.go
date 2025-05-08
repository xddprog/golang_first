package handlers

import (
	"encoding/json"
	"golang/internal/core/services"
	"golang/internal/handlers/dependencies"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"net/http"
	"strconv"
)


type UserHandler struct {
	Service *services.UserService
}


func (handler *UserHandler) GetUserById(response http.ResponseWriter, request *http.Request, user *models.UserModel) {
	response.Header().Set("Content-Type", "application/json")

	userId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	userGet, serviceErr := handler.Service.GetUserById(request.Context(), userId)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	response.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(response).Encode(userGet); err != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrEncodingError)
	}
}


func (handler *UserHandler) UpdateUser(response http.ResponseWriter, request *http.Request, user *models.UserModel) {
	response.Header().Set("Content-Type", "application/json")

	user, err := handler.Service.UpdateUser(request.Context(), user.Id, request.Body)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(user); err != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrEncodingError)
	}
}


func (handler *UserHandler) DeleteUser(response http.ResponseWriter, request *http.Request, user *models.UserModel) {
	response.Header().Set("Content-Type", "application/json")

	err := handler.Service.DeleteUser(request.Context(), user.Id)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
	response.Write([]byte(`{"message": "User deleted successfully"}`))

}


func (handler *UserHandler) SetupRoutes(server *http.ServeMux, baseUrl string, d *deps.AuthDependency) {
	server.HandleFunc("GET " + baseUrl+ "/user/{id}", d.Protected(handler.GetUserById))
	server.HandleFunc("PUT " + baseUrl + "/user", d.Protected(handler.UpdateUser))
	server.HandleFunc("DELETE " + baseUrl + "/user", d.Protected(handler.DeleteUser))
}