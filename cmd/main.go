package main

import (
	"golang/internal/handlers/setup"
	"golang/internal/handlers/v1"
	"golang/internal/infrastructure/database/connections"
	"log"
	"net/http"
)


func main() {
	db, err := db_connections.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServeMux()

	userHandler, err := setup.InitNewHandler(&handlers.UserHandler{}, db)
	userHandler.SetupRoutes(server, "/api/v1")

	http.ListenAndServe("localhost:8000", server)
}
