package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"todo-app/db"
	"todo-app/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func UserRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		userID := uuid.New().String() // Generate UUID
		_, err = db.DB.Exec("INSERT INTO users(id, username, password) VALUES(?, ?, ?)", userID, username, hashedPassword)
		if err != nil {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    userID, // Set UUID cookie
			HttpOnly: true,
			Path:     "/",
		})

		http.Redirect(w, r, "/profile", http.StatusSeeOther) // Redirect to profile creation
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

// UserLogin handles user login
func UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user models.User
		err := db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    user.ID, // Set UUID cookie
			HttpOnly: true,
			Path:     "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

// ProfileHandler handles profile creation/update
func UserProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")

		// Handle image upload
		file, header, err := r.FormFile("image")
		var imagePath string

		if err == nil {
			defer file.Close()
			imageName := uuid.New().String() + filepath.Ext(header.Filename)
			imagePath = filepath.Join("static/images", imageName)
			outFile, err := os.Create(imagePath)
			if err != nil {
				http.Error(w, "Error saving image", http.StatusInternalServerError)
				return
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, file); err != nil {
				http.Error(w, "Error saving image", http.StatusInternalServerError)
				return
			}
			imagePath = "/" + imagePath
		}

		_, err = db.DB.Exec("INSERT OR REPLACE INTO profiles(user_id, name, email, image_path) VALUES(?, ?, ?, ?)", userID, name, email, imagePath)
		if err != nil {
			http.Error(w, "Error updating profile", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/profile.html"))
	tmpl.Execute(w, nil)
}

// Authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := getUserIDFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getUserIDFromCookie(r *http.Request) (int, error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return 0, err
	}

	var userID int
	_, err = fmt.Sscanf(cookie.Value, "%d", &userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// IndexHandler with sorting & filtering
func Dasboard(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")   // e.g., "due_date"
	filter := r.URL.Query().Get("filter") // e.g., "pending"

	query := "SELECT id, title, description, notes, due_date, status FROM todos ORDER BY due_date ASC"
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
func CreateTask(w http.ResponseWriter, r *http.Request) {
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
		// taskReminder(title, dueDate)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// DeleteHandler handles deleting a todo
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// scheduleReminder triggers a reminder notification
// func taskReminder(taskTitle, dueDateStr string) {
// 	dueTime, err := time.Parse("2006-01-02T15:04", dueDateStr)
// 	if err != nil {
// 		return
// 	}

// 	delay := time.Until(dueTime)
// 	if delay > 0 {
// 		time.AfterFunc(delay, func() {
// 			log.Printf("🔔 Reminder: Task '%s' is due now!", taskTitle)
// 		})
// 	}
// }

// UpdateHandler for editing tasks
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		title := r.FormValue("title")
		description := r.FormValue("description")
		notes := r.FormValue("notes")
		dueDate := r.FormValue("due_date")
		status := r.FormValue("status")

		_, err := db.DB.Exec("UPDATE todos SET title=?, description=?, notes=?, due_date=?, status=? WHERE id=?", title, description, notes, dueDate, status, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func EditTask(w http.ResponseWriter, r *http.Request) {
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

// CompleteHandler handles marking a todo as completed
func CompleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var todo models.Todo
	err := db.DB.QueryRow("SELECT id, title, description, notes, due_date FROM todos WHERE id=?", id).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Notes, &todo.DueDate)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	_, err = db.DB.Exec("INSERT INTO completed_todos(title, description, notes, due_date) VALUES(?, ?, ?, ?)",
		todo.Title, todo.Description, todo.Notes, todo.DueDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CompletedTasksHandler displays completed tasks
func CompletedTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, description, notes, due_date, completed_at FROM completed_todos ORDER BY completed_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	completedTodos := []models.CompletedTodo{}
	for rows.Next() {
		var todo models.CompletedTodo
		rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Notes, &todo.DueDate, &todo.CompletedAt)
		completedTodos = append(completedTodos, todo)
	}

	tmpl := template.Must(template.ParseFiles("templates/complete.html"))
	tmpl.Execute(w, completedTodos)
}
