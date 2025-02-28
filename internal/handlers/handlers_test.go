package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	db "go_todo/internal/database"
	"go_todo/internal/models"
)

func SetupTestDB() {
	db.InitDB() // Временная бд
	db.Migrate()
	db.DB.Exec("DELETE FROM tasks")
}

func TestGetTask(t *testing.T) {
	SetupTestDB()

	db.DB.Exec("INSERT INTO tasks (title, content) VALUES ('Test Task 1', 'Test Content 1')")

	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var tasks []models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Title != "Test Task 1" || tasks[0].Content != "Test Content 1" {
		t.Errorf("Expected one task with title 'Test Task 1' and content 'Test Content 1', got %v", tasks)
	}

}

func TestPostTask(t *testing.T) {
	SetupTestDB()

	task := models.Task{Title: "Test Task2", Content: "Test Content2"}
	body, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var createdTask models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &createdTask)
	if err != nil {
		t.Fatalf("Failed to parse response body: %s", err)
	}

	// Проверяем значения
	if createdTask.Title != "Test Task2" || createdTask.Content != "Test Content2" {
		t.Errorf("Unexpected response body: %+v", createdTask)
	}

	// Проверяем временные метки
	if createdTask.CreatedAT.IsZero() || createdTask.UpdatedAT.IsZero() {
		t.Errorf("CreatedAT or UpdatedAT is zero: %+v", createdTask)
	}
}
