package handlers

import (
	"encoding/json"
	db "go_todo/internal/database"
	"go_todo/internal/models"
	"log"
	"net/http"
)

// Получение всех задач
func GetTask(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/tasks request received")
	if r.Method != http.MethodGet {
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.DB.Query("SELECT id, title, content, created_at, updated_at FROM tasks")
	if err != nil {
		log.Printf("DB query failed: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.CreatedAT, &task.UpdatedAT)
		if err != nil {
			log.Printf("DB scan failed: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
	log.Println("GET /api/tasks request completed")
}

// Создание задачи
func CreateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /api/tasks request received")
	if r.Method != http.MethodPost {
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("JSON decode failed: %s", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if task.Title == "" || task.Content == "" {
		log.Println("Title and content are required")
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("INSERT INTO tasks (title, content) VALUES (?, ?)", task.Title, task.Content)
	if err != nil {
		log.Printf("DB query failed: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("DB query failed: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем созданную задачу из базы данных
	row := db.DB.QueryRow("SELECT id, title, content, created_at, updated_at FROM tasks WHERE id = ?", id)
	err = row.Scan(&task.ID, &task.Title, &task.Content, &task.CreatedAT, &task.UpdatedAT)
	if err != nil {
		log.Printf("DB scan failed: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
	log.Println("POST /api/tasks request completed")
}
