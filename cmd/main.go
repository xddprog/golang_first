package main

import (
	"golang/internal/handlers/setup"
	"golang/internal/handlers/v1"
	"golang/internal/infrastructure/database/connections"
	"log"
	"net/http"
)


func main() {
	db, err := connections.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServeMux()

	userHandler, _ := setup.InitNewHandler(&handlers.UserHandler{}, db)
	authHandler, _ := setup.InitNewHandler(&handlers.AuthHandler{}, db)
	
	userHandler.SetupRoutes(server, "/api/v1")
	authHandler.SetupRoutes(server, "/api/v1")
	
	http.ListenAndServe("localhost:8000", server)
}
