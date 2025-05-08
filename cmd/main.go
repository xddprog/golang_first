package main

import (
	"golang/internal/handlers/dependencies"
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

	// redis := connections.NewRedisConnection()

	// rabbit, err := connections.NewRabbitConnection()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	server := http.NewServeMux()

	userHandler, _ := setup.InitNewHandler(&handlers.UserHandler{}, db)
	authHandler, _ := setup.InitNewHandler(&handlers.AuthHandler{}, db)
	documentHandler, _ := setup.InitNewHandler(&handlers.DocumentHandler{}, db)
	
	authDependency := deps.NewAuthDependency(authHandler.Service)

	userHandler.SetupRoutes(server, "/api/v1", authDependency)
	authHandler.SetupRoutes(server, "/api/v1", authDependency)
	documentHandler.SetupRoutes(server, "/api/v1", authDependency)

	http.ListenAndServe("localhost:8000", server)
}
