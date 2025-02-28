package main

import (
	db "go_todo/internal/database"
	"go_todo/internal/handlers"
	"go_todo/internal/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file", err)
	}
	defer file.Close()
	//НАстройКа логгера для записи в файл
	log.SetOutput(file)
	log.Println("Application started.")
	//Логирование: Начало работы приложения
	log.Println("Starting the application...")
	// Иницилизация БД
	err = db.InitDB()
	if err != nil {
		log.Fatalf("DB connection is falled: %s", err)
	}
	defer db.DB.Close()
	log.Println("DB connection established.")

	//Миграции
	err = db.Migrate()
	if err != nil {
		log.Fatalf("DB migration is falled: %s", err)
	}
	log.Println("database migrated.")
	//Создание обработчика
	r := mux.NewRouter()
	//Регистрация роутов
	r.HandleFunc("/api/tasks", handlers.GetTask).Methods("GET")
	r.HandleFunc("/api/tasks", handlers.CreateTask).Methods("POST")

	loggedR := middleware.LoggingMiddleware(r)
	//запуск сервера
	log.Println("server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedR))
}
