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

	go func() { // Run reminder checking in a separate goroutine
		for {
			utils.CheckReminders()
			time.Sleep(30 * time.Minute) // Check every 30 minutes
		}
	}()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/create", handlers.CreateHandler)
	http.HandleFunc("/delete", handlers.DeleteHandler)
	http.HandleFunc("/update", handlers.UpdateHandler)
	http.HandleFunc("/edit", handlers.EditHandler)

	fmt.Println("🚀 Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
