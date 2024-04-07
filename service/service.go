package service

import (
	"context"
	"encoding/json"
	"keeper/database"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	HealthService(http.ResponseWriter, *http.Request)
	GetTodo(http.ResponseWriter, *http.Request)
	GetTodos(http.ResponseWriter, *http.Request)
	CreateTodo(http.ResponseWriter, *http.Request)
	UpdateTodo(http.ResponseWriter, *http.Request)
	DeleteTodo(http.ResponseWriter, *http.Request)
	Close()
}

type service struct {
	db     database.Database
	logger *slog.Logger
}

func New(db database.Database) Service {
	logger := slog.Default()
	return &service{db: db, logger: logger}
}

func (s *service) HealthService(res http.ResponseWriter, req *http.Request) {
	s.logger.Info("control inside HealthService")
	content, err := json.Marshal(s.db.Health())
	if err != nil {
		s.logger.Warn("could not create json.")
		res.Write([]byte(err.Error()))
	}
	res.Write(content)
}

func (s *service) GetTodo(res http.ResponseWriter, req *http.Request) {
	s.logger.Info("control inside GetTodo")
	idStr := req.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		s.logger.Warn("invalid request")
		res.Write([]byte(err.Error()))
	}

	// timeout in 3 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	// to stop after 3 seconds.
	defer cancel()

	todo, err := s.db.(database.TodoModel).GetTodo(ctx, id)
	if err != nil {
		// handle failed to get from db.
		s.logger.Warn(err.Error())
		res.Write([]byte("error getting data from db"))
	}

	content, err := json.Marshal(todo)
	if err != nil {
		s.logger.Warn("could not create json.")
		res.Write([]byte(err.Error()))
	}
	res.Write(content)
}

func (s *service) GetTodos(res http.ResponseWriter, req *http.Request) {
	s.logger.Info("control inside GetTodos")
	// timeout in 3 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	// to stop after 3 seconds.
	defer cancel()

	todos, err := s.db.(database.TodoModel).ListTodos(ctx)

	if err != nil {
		// handle failed to get from db.
		s.logger.Warn(err.Error())
		res.Write([]byte("error getting data from db"))
	}

	content, err := json.Marshal(todos)
	if err != nil {
		s.logger.Warn("could not create json.")
		res.Write([]byte(err.Error()))
	}
	res.Write(content)
}

func (s *service) CreateTodo(res http.ResponseWriter, req *http.Request) {
	s.logger.Info("control inside CreateTodo")
	decoder := json.NewDecoder(req.Body)
	var todo database.Todo
	err := decoder.Decode(&todo)

	if err != nil {
		s.logger.Warn("invalid request")
		res.Write([]byte(err.Error()))
	}

	// timeout in 3 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	// to stop after 3 seconds.
	defer cancel()

	err = s.db.(database.TodoModel).Save(ctx, &todo)

	if err != nil {
		s.logger.Warn("internal server error")
		res.Write([]byte(err.Error()))
	}

	content, err := json.Marshal(&todo)
	if err != nil {
		s.logger.Warn("could not create json.")
		res.Write([]byte(err.Error()))
	}
	res.Write(content)
}

func (s *service) UpdateTodo(res http.ResponseWriter, req *http.Request) {
	s.logger.Info("control inside UpdateTodo")

	idStr := req.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		s.logger.Warn("invalid request")
		res.Write([]byte(err.Error()))
	}

	decoder := json.NewDecoder(req.Body)
	var todo database.Todo
	err = decoder.Decode(&todo)

	if err != nil {
		s.logger.Warn("invalid request")
		res.Write([]byte(err.Error()))
	}

	// timeout in 3 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	// to stop after 3 seconds.
	defer cancel()

	err = s.db.(database.TodoModel).UpdateTodo(ctx, id, &todo)

	if err != nil {
		s.logger.Warn("internal server error")
		res.Write([]byte(err.Error()))
	}

	content, err := json.Marshal(&todo)
	if err != nil {
		s.logger.Warn("could not create json.")
		res.Write([]byte(err.Error()))
	}
	res.Write(content)
}

func (s *service) DeleteTodo(res http.ResponseWriter, req *http.Request) {
	s.logger.Info("control inside DeleteTodo")

	idStr := req.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		s.logger.Warn("invalid request")
		res.Write([]byte(err.Error()))
	}

	// timeout in 3 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	// to stop after 3 seconds.
	defer cancel()

	err = s.db.(database.TodoModel).DeleteTodo(ctx, id)

	if err != nil {
		s.logger.Warn("internal server error")
		res.Write([]byte(err.Error()))
	}

	if err != nil {
		s.logger.Warn("could not create json.")
		res.Write([]byte(err.Error()))
	}
	res.Write([]byte("deleted."))
}

func (s *service) Close() {
	s.db.Close()
}
