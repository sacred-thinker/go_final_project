package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/internal/config"
	"go_final_project/internal/handler"
	"go_final_project/internal/repository"
	"go_final_project/internal/service"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	log.Println("Database открыта.")

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	db, err := repository.InitDB("scheduler.db")
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)
	authHandler := handler.NewAuthHandler(cfg)
	nextDateHandler := handler.NewNextDateHandler()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/signin", authHandler.SignIn)
	http.HandleFunc("/api/task", authHandler.AuthMiddleware(taskHandler.HandleTask))
	http.HandleFunc("/api/tasks", authHandler.AuthMiddleware(taskHandler.GetTasks))
	http.HandleFunc("/api/task/done", authHandler.AuthMiddleware(taskHandler.TaskDone))
	http.HandleFunc("/api/nextdate", nextDateHandler.HandleNextDate)

	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
