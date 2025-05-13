package handlers

import (
	"golang/internal/core/services"
	"golang/internal/handlers/dependencies"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"golang/internal/utils"
	"net/http"
	"strconv"
	"sync"

	"github.com/googollee/go-socket.io"
)


type DocumentHandler struct {
	DocumentService 	*services.DocumentService
	CommentService  	*services.CommentService
	Socket      		*socketio.Server
	Connections 		map[string]map[string]models.BaseUserModel
	Mutex 				sync.RWMutex
}


func (handler *DocumentHandler) CreateDocument(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	document, err := handler.DocumentService.CreateDocument(request.Context(), user.Id, request.Body)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}
	utils.WriteJSONResponse(response, http.StatusCreated, document)
}


func (handler *DocumentHandler) UpdateDocument(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	documentId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	document, serviceErr := handler.DocumentService.UpdateDocument(request.Context(), user.Id, documentId, request.Body)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	utils.WriteJSONResponse(response, http.StatusOK, document)
}


func (handler *DocumentHandler) DeleteDocument(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	documentId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	if err := handler.DocumentService.DeleteDocument(request.Context(), documentId, user.Id); err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func (handler *DocumentHandler) AddDocumentSnapshot(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	documentId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	snapshot, serviceErr := handler.DocumentService.AddDocumentSnapshot(request.Context(), user.Id, documentId)
	if serviceErr != nil {
		apierrors.WriteHTTPError(response, serviceErr)
		return
	}

	utils.WriteJSONResponse(response, http.StatusOK, snapshot)
}


func (handler *DocumentHandler) SendInvite(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	documentId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	smtpErr := handler.DocumentService.SendInvite(request.Context(), user.Id, request.PathValue("email"), documentId)
	if smtpErr != nil {
		apierrors.WriteHTTPError(response, smtpErr)
		return
	}

	utils.WriteJSONResponse(response, http.StatusOK, `{"detail": "Invite sent successfully"}`)
}


func (handler *DocumentHandler) GetComments(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	documentId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	if err := handler.DocumentService.CheckDocumentAccess(request.Context(), user.Id, documentId); err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	limit, offset := utils.GetLimitAndOffset(request)
	comments, dbErr := handler.CommentService.GetCommentsByDocument(request.Context(), documentId, limit, offset)
	if dbErr != nil {
	    apierrors.WriteHTTPError(response, dbErr)
		return
	}
	utils.WriteJSONResponse(response, http.StatusOK, comments)
}


func (handler *DocumentHandler) AddComment(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")
	documentId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	if err := handler.DocumentService.CheckDocumentAccess(request.Context(), user.Id, documentId); err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	comment, err := handler.CommentService.CreateComment(request.Context(), user, documentId, request.Body)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return	
	}
	
	utils.WriteJSONResponse(response, http.StatusOK, comment)
}


func (handler *DocumentHandler) UpdateComment(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")
	
	commentId, commentIdErr := strconv.Atoi(request.PathValue("commentId"))
	if commentIdErr != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrInvalidRequestBody)
		return
	}
	documentId, documentIdErr := strconv.Atoi(request.PathValue("documentId"))
	if documentIdErr != nil {
	    apierrors.WriteHTTPError(response, apierrors.ErrInvalidRequestBody)
		return
	}

	if err := handler.DocumentService.CheckDocumentAccess(request.Context(), user.Id, documentId); err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	comment, err := handler.CommentService.UpdateComment(request.Context(), commentId, user.Id, request.Body)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	utils.WriteJSONResponse(response, http.StatusOK, comment)

}


func (handler *DocumentHandler) DeleteComment(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	commentId, commentIdErr := strconv.Atoi(request.PathValue("commentId"))
	if commentIdErr != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrInvalidRequestBody)
		return
	}
	documentId, documentIdErr := strconv.Atoi(request.PathValue("documentId"))
	if documentIdErr != nil {
	    apierrors.WriteHTTPError(response, apierrors.ErrInvalidRequestBody)
		return
	}

	if err := handler.DocumentService.CheckDocumentAccess(request.Context(), user.Id, documentId); err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	err := handler.CommentService.DeleteComment(request.Context(), user.Id, commentId, documentId)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}


func (handler *DocumentHandler) GetCommentsReplies(response http.ResponseWriter, request *http.Request, user *models.BaseUserModel) {
	response.Header().Set("Content-Type", "application/json")

	commentId, commentIdErr := strconv.Atoi(request.PathValue("commentId"))
	if commentIdErr != nil {
		apierrors.WriteHTTPError(response, apierrors.ErrInvalidRequestBody)
		return
	}
	documentId, documentIdErr := strconv.Atoi(request.PathValue("documentId"))
	if documentIdErr != nil {
	    apierrors.WriteHTTPError(response, apierrors.ErrInvalidRequestBody)
		return
	}

	if err := handler.DocumentService.CheckDocumentAccess(request.Context(), user.Id, documentId); err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	comments, err := handler.CommentService.GetCommentsReplies(request.Context(), commentId)
	if err != nil {
		apierrors.WriteHTTPError(response, err)
		return
	}

	utils.WriteJSONResponse(response, http.StatusOK, comments)
}


func (handler *DocumentHandler) SetupRoutes(server *http.ServeMux, baseUrl string, d *deps.AuthDependency) {
	server.HandleFunc("POST " + baseUrl+ "/documents", d.Protected(handler.CreateDocument))
	server.HandleFunc("PUT " + baseUrl+ "/documents", d.Protected(handler.UpdateDocument))
	server.HandleFunc("DELETE " + baseUrl+ "/documents", d.Protected(handler.DeleteDocument))
	server.HandleFunc("POST " + baseUrl+ "/documents/{id}/invite", d.Protected(handler.SendInvite))
	server.HandleFunc("POST " + baseUrl+ "/documents/{id}/snapshot", d.Protected(handler.AddDocumentSnapshot))

	server.HandleFunc("POST " + baseUrl+ "/documents/{id}/comments", d.Protected(handler.AddComment))
	server.HandleFunc("GET " + baseUrl+ "/documents/{id}/comments", d.Protected(handler.GetComments))
	server.HandleFunc("GET " + baseUrl+ "/documents/{id}/comments/{commentId}", d.Protected(handler.GetCommentsReplies))
	server.HandleFunc("PUT " + baseUrl+ "/documents/{documentId}/comments/{commentId}", d.Protected(handler.UpdateComment))
	server.HandleFunc("DELETE " + baseUrl+ "/documents/{id}/comments/{commentId}", d.Protected(handler.DeleteComment))
	server.Handle(baseUrl + "documents/ws/{id}", handler.Socket)
}