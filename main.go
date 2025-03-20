package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"todo-app/db"
	"todo-app/routes"
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

	routes.Routes()

	fmt.Println("🚀 Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
