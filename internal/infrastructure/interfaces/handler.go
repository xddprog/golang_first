package interfaces

import (
	deps "golang/internal/handlers/dependencies"
	"net/http"
)


type HandlerInterface interface {
	SetupRoutes(server *http.ServeMux, baseUrl string, protected *deps.AuthDependency)
}
