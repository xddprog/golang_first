package handlers

import (
	"context"
	"encoding/json"
	deps "golang/internal/handlers/dependencies"
	"golang/internal/infrastructure/database/models"
	apierrors "golang/internal/infrastructure/errors"
	"log"
	"strconv"

	socketio "github.com/googollee/go-socket.io"
)


func (handler *DocumentHandler) notifyUsers(documentId string) {
	handler.Mutex.RLock()
	defer handler.Mutex.RUnlock()

	users := make([]*models.BaseUserModel, 0, len(handler.Connections[documentId]))

	for _, user := range handler.Connections[documentId] {
		users = append(users, &user)
	}

	handler.Socket.BroadcastToRoom("/", "doc_" + documentId, "users_update", users)
}


func (handler *DocumentHandler) HandleConnect(s socketio.Conn) error {
	s.SetContext("")
	return nil
}


func (handler *DocumentHandler) HandleJoinDocument(s socketio.Conn, documentId string, user *models.BaseUserModel) {
	convDocumentId, err := strconv.Atoi(documentId)
	if err != nil {
		s.Emit("error", "invalid documentId")
		return
	}

	_, dbErr := handler.DocumentService.GetDocumentById(context.Background(), convDocumentId, user.Id)
	if dbErr != nil {
		s.Emit("error", dbErr.Error())
		return
	}

	s.Join("doc_" + documentId)

	handler.Mutex.Lock()
	if handler.Connections[documentId] == nil {
		handler.Connections[documentId] = make(map[string]models.BaseUserModel)
	}
	handler.Connections[documentId][s.ID()] = *user
	handler.Mutex.Unlock()

	document, _ := handler.DocumentService.GetDocumentById(context.Background(), convDocumentId, user.Id)
	s.Emit("document_state", document)

	handler.notifyUsers(documentId)
}


func (handler *DocumentHandler) HandleDocumentUpdate(s socketio.Conn, data string, user *models.BaseUserModel) {
	var update models.UpdateDocumentContent

	if err := json.Unmarshal([]byte(data), &update); err != nil {
		s.Emit("error", apierrors.ErrEncodingError.Error())
		return 
	}

	_, err := strconv.Atoi(update.DocumentId)
	if err != nil {
		s.Emit("error", "invalid documentId")
		return
	}

	document, updateErr := handler.DocumentService.UpdateDocumentContent(
		context.Background(), user.Id, update.DocumentId, update.Content,
	)
	if updateErr != nil {
		s.Emit("error", updateErr.Error())
		return
	}

	s.Emit("document_updated", document)

}


func (handler *DocumentHandler) HandlerCursorMove(s socketio.Conn, data string, user *models.BaseUserModel) {
	var move models.CursorMove

	if err := json.Unmarshal([]byte(data), &move); err != nil {
		s.Emit("error", apierrors.ErrEncodingError.Error())
		return
	}

	handler.Socket.BroadcastToRoom("/", "doc_" + move.DocumentId, "cursor_move", move.Position)
}


func (handler *DocumentHandler) HandleDisconnect(s socketio.Conn, reason string) {
	handler.Mutex.Lock()
	defer handler.Mutex.Unlock()

	for documentId, users := range handler.Connections {
		if _, exists := users[s.ID()]; exists {
			delete(handler.Connections[documentId], s.ID())
			go handler.notifyUsers(documentId)
		}
	}
}


func (handler *DocumentHandler) RunWebsocket() {
	go func() {
		if err := handler.Socket.Serve(); err != nil {
			log.Printf("Socket.IO error: %v", err)
		}
	}()
}


func (handler *DocumentHandler) SetupSocket(server *socketio.Server, d *deps.AuthDependency) {
	handler.Socket.OnConnect("/", d.ProtectConnect(handler.HandleConnect))
	handler.Socket.OnEvent("/", "join", d.ProtectEvent(handler.HandleJoinDocument))
	handler.Socket.OnEvent("/", "update", d.ProtectEvent(handler.HandleDocumentUpdate))
	handler.Socket.OnEvent("/", "cursor_move", d.ProtectEvent(handler.HandlerCursorMove))
	handler.Socket.OnDisconnect("/", handler.HandleDisconnect)
}