package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"sketchive/internal/api"
	"sketchive/internal/db"

	_ "github.com/go-sql-driver/mysql"
)

// CORS: Cross-Origin Resource Sharing

// CORS middleware function
func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins (change for production)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	//dsn: Data Source Name
	dsn := "root:@tcp(127.0.0.1:3306)/sketchive"
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Could not grad the connvection at first place", err)
	}

	// Ping() checks if the connection is alive
	err = database.Ping()
	if err != nil {
		log.Fatal("Lost Database connection failed:", err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	db.SetDB(database)

	mux := http.NewServeMux()

	mux.HandleFunc("/whiteboards", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			api.GetWhiteboard(w, r) // Handle GET request
		case "POST":
			api.CreateWhiteboard(w, r) // Handle POST request (if needed)
		case "PUT":
			api.UpdateWhiteboard(w, r) // Handle PUT request
		case "DELETE":
			api.DeleteWhiteboard(w, r) // Handle DELETE request
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Add more routes for strokes or users if needed
	// mux.HandleFunc("/strokes", ...)

	// Start the server with CORS enabled
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))

}
