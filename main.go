package main

import (
	"keeper/database"
	"keeper/service"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	logger := slog.Default()
	logger.Info("starting... http://localhost:9090")

	db := database.New("todos.db")
	svc := service.New(db)
	defer svc.Close()

	// this is go 1.22 new feature, in real world people prefer `echo, fiber or chi` routers.
	http.HandleFunc("GET /", svc.HealthService)
	http.HandleFunc("GET /todo", svc.GetTodos)
	http.HandleFunc("POST /todo", svc.CreateTodo)
	http.HandleFunc("GET /todo/{id}", svc.GetTodo)
	http.HandleFunc("PUT /todo/{id}", svc.UpdateTodo)
	http.HandleFunc("DELETE /todo/{id}", svc.DeleteTodo)

	log.Fatal(http.ListenAndServe(":9090", nil))
}
