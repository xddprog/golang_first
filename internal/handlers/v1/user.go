package handlers

import (
	"encoding/json"
	"golang/internal/core/services"
	"net/http"
	"strconv"
)


type UserHandler struct {
	Service *services.UserService
}


func (handler *UserHandler) GetUserById(response http.ResponseWriter, request *http.Request) {
	userId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		http.Error(response, "Invalid user id", http.StatusUnprocessableEntity)
	}
	user, err := handler.Service.GetUserById(request.Context(), userId)

    response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(response).Encode(user); err != nil {
		http.Error(response, "Encoding error", http.StatusInternalServerError)
	}
}


func (handler *UserHandler) SetupRoutes(server *http.ServeMux, baseUrl string) {
	server.HandleFunc("GET " + baseUrl+ "/user/{id}", handler.GetUserById)
}