package handlers

import (
	"golang/internal/core/services"
	"golang/internal/handlers/dependencies"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"golang/internal/utils"
	"net/http"
	"strconv"
)


type UserHandler struct {
	UserService *services.UserService
	DocumentService *services.DocumentService
}


func (handler *UserHandler) GetUserById(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	userId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	userGet, serviceErr := handler.UserService.GetUserById(request.Context(), userId)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	utils.WriteJSONResponse(response, http.StatusOK, userGet)
}


func (handler *UserHandler) UpdateUser(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	user, err := handler.UserService.UpdateUser(request.Context(), user.Id, request.Body)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	utils.WriteJSONResponse(response, http.StatusOK, user)
}


func (handler *UserHandler) DeleteUser(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	err := handler.UserService.DeleteUser(request.Context(), user.Id)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
	utils.WriteJSONResponse(response, http.StatusOK, `{"message": "User deleted successfully"}`)
}


func (handler *UserHandler) GetUserDocuments(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	limit, offset := utils.GetLimitAndOffset(request)

	userId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	documents, serviceErr := handler.DocumentService.GetUserDocuments(request.Context(), userId, limit, offset)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	utils.WriteJSONResponse(response, http.StatusOK, documents)
}


func (handler *UserHandler) SetupRoutes(server *http.ServeMux, baseUrl string, d *deps.AuthDependency) {
	server.HandleFunc("GET " + baseUrl+ "/user/{id}", d.Protected(handler.GetUserById))
	server.HandleFunc("PUT " + baseUrl + "/user", d.Protected(handler.UpdateUser))
	server.HandleFunc("DELETE " + baseUrl + "/user", d.Protected(handler.DeleteUser))
	server.HandleFunc("GET " + baseUrl + "/user", d.Protected(handler.GetUserDocuments))
}