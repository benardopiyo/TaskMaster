package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"todo-app/db"
	"todo-app/models"
)

// IndexHandler with sorting & filtering
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")   // e.g., "due_date"
	filter := r.URL.Query().Get("filter") // e.g., "pending"

	query := "SELECT id, title, description, notes, due_date, status FROM todos"
	if filter != "" {
		query += " AND status = '" + filter + "'"
	}
	if sortBy != "" {
		query += " ORDER BY " + sortBy
	}

	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var todo models.Todo
		rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Notes, &todo.DueDate, &todo.Status)
		todos = append(todos, todo)
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, todos)
}

// CreateHandler handles creating a new todo
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		description := r.FormValue("description")
		notes := r.FormValue("notes")
		dueDate := r.FormValue("due_date")

		_, err := db.DB.Exec("INSERT INTO todos(title, description, notes, due_date, status) VALUES(?, ?, ?, ?, 'pending')",
			title, description, notes, dueDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Schedule a reminder
		scheduleReminder(title, dueDate)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// DeleteHandler handles deleting a todo
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// scheduleReminder triggers a reminder notification
func scheduleReminder(taskTitle, dueDateStr string) {
	dueTime, err := time.Parse("2006-01-02T15:04", dueDateStr)
	if err != nil {
		return
	}

	delay := time.Until(dueTime)
	if delay > 0 {
		time.AfterFunc(delay, func() {
			log.Printf("🔔 Reminder: Task '%s' is due now!", taskTitle)
		})
	}
}

// UpdateHandler for editing tasks
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		title := r.FormValue("title")
		description := r.FormValue("description")
		notes := r.FormValue("notes")
		dueDate := r.FormValue("due_date")
		status := r.FormValue("status")

		_, err := db.DB.Exec("UPDATE todos SET title=?, description=?, due_date=?, status=? WHERE id=?", title, description, notes, dueDate, status, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    var todo models.Todo
    err := db.DB.QueryRow("SELECT id, title, description, notes, due_date, status FROM todos WHERE id=?", id).
        Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Notes, &todo.DueDate, &todo.Status)
    if err != nil {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    tmpl := template.Must(template.ParseFiles("templates/edit.html"))
    tmpl.Execute(w, todo)
}


