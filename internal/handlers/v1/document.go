package handlers

import (
	"encoding/json"
	"golang/internal/core/services"
	deps "golang/internal/handlers/dependencies"
	"golang/internal/infrastructure/database/models"
	apierrors "golang/internal/infrastructure/errors"
	"net/http"
	"strconv"
)


type DocumentHandler struct {
	Service *services.DocumentService
}


func (handler *DocumentHandler) CreateDocument(response http.ResponseWriter, request *http.Request, user *models.UserModel) {
	response.Header().Set("Content-Type", "application/json")

	document, err := handler.Service.CreateDocument(request.Context(), user.Id, request.Body)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(response).Encode(document); err != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrEncodingError)
	}
}

func (handler *DocumentHandler) GetDocumentById(response http.ResponseWriter, request *http.Request, user *models.UserModel) {
	response.Header().Set("Content-Type", "application/json")

	documentId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	document, serviceErr := handler.Service.GetDocumentById(request.Context(), documentId)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(document); err != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrEncodingError)
	}
}


func (handler *DocumentHandler) UpdateDocument(response http.ResponseWriter, request *http.Request, user *models.UserModel) {
	response.Header().Set("Content-Type", "application/json")

	documentId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	document, serviceErr := handler.Service.UpdateDocument(request.Context(), documentId, request.Body)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(document); err != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrEncodingError)
	}
}


func (handler *DocumentHandler) DeleteDocument(response http.ResponseWriter, request *http.Request, user *models.UserModel) {
	response.Header().Set("Content-Type", "application/json")

	documentId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	handler.Service.DeleteDocument(request.Context(), documentId)

	response.WriteHeader(http.StatusNoContent)
}


func (handler *DocumentHandler) DocumentEditWebsocket(response http.ResponseWriter, request *http.Request, user *models.UserModel) {

}


func (handler *DocumentHandler) SetupRoutes(server *http.ServeMux, baseUrl string, d *deps.AuthDependency) {
	server.HandleFunc("POST " + baseUrl+ "/documents", d.Protected(handler.CreateDocument))
	server.HandleFunc("GET " + baseUrl+ "/documents", d.Protected(handler.GetDocumentById))
	server.HandleFunc("PUT " + baseUrl+ "/documents", d.Protected(handler.UpdateDocument))
	server.HandleFunc("DELETE " + baseUrl+ "/documents", d.Protected(handler.DeleteDocument))
}