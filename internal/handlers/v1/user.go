package handlers

import (
	"encoding/json"
	"golang/internal/core/services"
	"golang/internal/infrastructure/errors"
	"net/http"
	"strconv"
)


type UserHandler struct {
	Service *services.UserService
}


func (handler *UserHandler) GetUserById(response http.ResponseWriter, request *http.Request) {
	userId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	user, serviceErr := handler.Service.GetUserById(request.Context(), userId)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(response).Encode(user); err != nil {
		http.Error(response, "Encoding error", http.StatusInternalServerError)
	}
}


func (handler *UserHandler) SetupRoutes(server *http.ServeMux, baseUrl string) {
	server.HandleFunc("GET " + baseUrl+ "/user/{id}", handler.GetUserById)
	// server.HandleFunc("PUT " + baseUrl + "/user/{id}", handler.UpdateUser)
	// server.HandleFunc("DELETE " + baseUrl + "/user/{id}", handler.DeleteUser)
}