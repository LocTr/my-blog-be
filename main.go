package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/LocTr/my-blog-be/models"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver for pgx
)

func main() {
	// connect to db
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// verify connectivity now so the container fails fast if DB is unreachable
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	defer db.Close()

	// create table if not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	//create router
	router := mux.NewRouter()

	// Add middleware
	router.Use(jsonContentTypeMiddleware)
	router.Use(corsMiddleware)

	// Register POST route BEFORE GET to ensure proper method matching
	router.HandleFunc("/users", createUser(db)).Methods("POST")
	router.HandleFunc("/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/users/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser(db)).Methods("DELETE")

	// Add a simple test endpoint to verify POST routing works
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Test endpoint called with method: %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"method":  r.Method,
			"message": "Test endpoint works",
		})
	}).Methods("POST", "GET")

	// Start the HTTP server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log ALL incoming requests to debug routing
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
		log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
} // get all users
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getUsers called with method: %s", r.Method)
		rows, err := db.Query("SELECT id, name, email FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

// get single user by ID
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var user models.User
		row := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id)
		if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		log.Printf("Fetched user: %+v", user)

		json.NewEncoder(w).Encode(user)
	}
}

// create user
func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("createUser called with method: %s", r.Method)

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("Failed to decode request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Creating user: %+v", user)

		// insert user into db
		err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
			user.Name, user.Email).Scan(&user.ID)
		if err != nil {
			log.Printf("Failed to insert user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Successfully created user: %+v", user)
		json.NewEncoder(w).Encode(user)
	}
}

// update user
func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// update user in db
		_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3",
			user.Name, user.Email, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

// delete user
func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		// delete user from db
		_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
