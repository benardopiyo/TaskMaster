package routes

import (
	"net/http"

	"todo-app/handlers"
)

func Routes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/register", handlers.UserRegister)
	http.HandleFunc("/login", handlers.UserLogin)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/", handlers.IndexHandler)
	protectedMux.HandleFunc("/create", handlers.CreateTask)
	protectedMux.HandleFunc("/delete", handlers.DeleteTask)
	protectedMux.HandleFunc("/update", handlers.UpdateTask)
	protectedMux.HandleFunc("/edit", handlers.EditTask)
	protectedMux.HandleFunc("/complete", handlers.CompleteTask)
	protectedMux.HandleFunc("/completed_tasks", handlers.CompletedTasks)
	protectedMux.HandleFunc("/profile", handlers.UserProfile)

	http.Handle("/", handlers.AuthMiddleware(protectedMux)) // Apply AuthMiddleware to protected routes
}
