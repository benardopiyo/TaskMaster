// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"todo-app/db"
// 	"todo-app/handlers"
// 	"todo-app/utils"
// )

// func main() {
// 	db.InitDB()
// 	defer db.CloseDB()

// 	go func() { // Run reminder checking in a separate goroutine
// 		for {
// 			utils.CheckReminders()
// 			time.Sleep(5 * time.Minute) // Check every 30 minutes
// 		}
// 	}()

// 	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

// 	http.HandleFunc("/", handlers.IndexHandler)
// 	http.HandleFunc("/create", handlers.CreateHandler)
// 	http.HandleFunc("/delete", handlers.DeleteHandler)
// 	http.HandleFunc("/update", handlers.UpdateHandler)
// 	http.HandleFunc("/edit", handlers.EditHandler)
// 	http.HandleFunc("/complete", handlers.CompleteHandler) // Add complete handler
// 	http.HandleFunc("/completed_tasks", handlers.CompletedTasksHandler) // Add completed tasks handler

// 	fmt.Println("🚀 Server is running at http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"todo-app/db"
	"todo-app/handlers"
	"todo-app/utils"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	go func() {
		for {
			utils.CheckReminders()
			time.Sleep(5 * time.Minute)
		}
	}()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/", handlers.IndexHandler)
	protectedMux.HandleFunc("/create", handlers.CreateHandler)
	protectedMux.HandleFunc("/delete", handlers.DeleteHandler)
	protectedMux.HandleFunc("/update", handlers.UpdateHandler)
	protectedMux.HandleFunc("/edit", handlers.EditHandler)
	protectedMux.HandleFunc("/complete", handlers.CompleteHandler)
	protectedMux.HandleFunc("/completed_tasks", handlers.CompletedTasksHandler)
	protectedMux.HandleFunc("/profile", handlers.ProfileHandler) // Add profile handler

	http.Handle("/", handlers.AuthMiddleware(protectedMux)) // Apply AuthMiddleware to protected routes

	fmt.Println("🚀 Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
