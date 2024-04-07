package service_test

import (
	"bytes"
	"encoding/json"
	"io"
	"keeper/database"
	"keeper/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() *httptest.Server {
	mockDb := database.New(":memory:")
	svc := service.New(mockDb)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /todo", svc.CreateTodo)
	mux.HandleFunc("GET /todo", svc.GetTodos)
	mux.HandleFunc("GET /todo/{id}", svc.GetTodo)
	mux.HandleFunc("PUT /todo/{id}", svc.UpdateTodo)
	mux.HandleFunc("DELETE /todo/{id}", svc.DeleteTodo)
	test_server := httptest.NewServer(mux)
	return test_server
}

func TestEnd2End(t *testing.T) {
	// Arrange
	ts := setup()

	defer ts.Close()

	test_todo := &database.Todo{
		Title:       "Test todo",
		Description: "Test description",
		Completed:   false,
	}

	body, _ := json.Marshal(test_todo)

	// Act - create todo
	resp, err := http.Post(ts.URL+"/todo", "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	// Assert - create todo
	resBody, err := io.ReadAll(resp.Body)

	singleTodo := "{\"Id\":1,\"Title\":\"Test todo\",\"Description\":\"Test description\",\"Completed\":false}"

	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	if singleTodo != string(resBody) {
		t.Errorf("expected response body to be %v; got %v", singleTodo, string(resBody))
	}

	// Act - get todos
	resp, err = http.Get(ts.URL + "/todo")

	if err != nil {
		t.Fatal(err)
	}

	// Assert - create todo
	resBody, err = io.ReadAll(resp.Body)

	arrayOfSingleTodo := "[{\"Id\":1,\"Title\":\"Test todo\",\"Description\":\"Test description\",\"Completed\":false}]"

	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	if arrayOfSingleTodo != string(resBody) {
		t.Errorf("expected response body to be %v; got %v", arrayOfSingleTodo, string(resBody))
	}
}
