package setup

import (
	"fmt"
	"golang/internal/core/repositories"
	"golang/internal/core/services"
	"golang/internal/handlers/v1"
	"golang/internal/infrastructure/config"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/types"

	socketio "github.com/googollee/go-socket.io"
	"github.com/jackc/pgx/v5/pgxpool"
)


func InitNewHandler[T types.HandlerInterface](emptyHandler T, db *pgxpool.Pool) (T, error) {
	switch h := any(emptyHandler).(type) {
	case *handlers.UserHandler:
		userRepository := &repositories.UserRepository{DB: db}
		documentRepository := &repositories.DocumentRepository{DB: db}

		userService := &services.UserService{Repository: userRepository}
		documentService := &services.DocumentService{Repository: documentRepository}
		
		*h = handlers.UserHandler{UserService: userService, DocumentService: documentService}
		return any(h).(T), nil
		
	case *handlers.AuthHandler:		
		cfg, err := config.LoadJwtConfig()
		if err != nil {
			panic(err)
		}
		
		repository := &repositories.UserRepository{DB: db}
		service := &services.AuthService{Repository: repository, Config: cfg}
		*h = handlers.AuthHandler{Service: service}
		return any(h).(T), nil

	case *handlers.DocumentHandler:
		documentRepository := &repositories.DocumentRepository{DB: db}
		commentRepository := &repositories.CommentRepository{DB: db}

		documentService := &services.DocumentService{Repository: documentRepository}
		commentService := &services.CommentService{Repository: commentRepository}
		
		socket := socketio.NewServer(nil)
		*h = handlers.DocumentHandler{
			DocumentService: documentService, 
			CommentService: commentService,
			Socket: socket, 
			Connections: make(map[string]map[string]models.BaseUserModel),
		}
		return any(h).(T), nil
		
	default:
		return emptyHandler, fmt.Errorf("undefined handler type: %T", emptyHandler)
	}
}