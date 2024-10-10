package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	// "sketchive/internal/api"
	"sketchive/internal/db"

	_ "github.com/go-sql-driver/mysql"
)

// OLD CODE using HTTP requests

// func main() {
// 	//dsn: Data Source Name
// 	dsn := "root:@tcp(127.0.0.1:3306)/sketchive"
// 	database, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatal("Could not grad the connvection at first place", err)
// 	}

// 	// Ping() checks if the connection is alive
// 	err = database.Ping()
// 	if err != nil {
// 		log.Fatal("Lost Database connection failed:", err)
// 	} else {
// 		fmt.Println("Successfully connected to database!")
// 	}

// 	db.SetDB(database)

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/whiteboards", func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case "GET":
// 			api.GetWhiteboard(w, r) // Handle GET request
// 		case "POST":
// 			api.CreateWhiteboard(w, r) // Handle POST request (if needed)
// 		case "PUT":
// 			api.UpdateWhiteboard(w, r) // Handle PUT request
// 		case "DELETE":
// 			api.DeleteWhiteboard(w, r) // Handle DELETE request
// 		default:
// 			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		}
// 	})

// 	mux.HandleFunc("/whiteboards/clear", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == "DELETE" {
// 			api.ClearWhiteboardHandler(w, r)
// 		} else {
// 			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		}
// 	})

// 	// Stroke related endpoints
// 	mux.HandleFunc("/strokes", func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case "POST":
// 			api.AddStroke(w, r) // Handle adding a stroke
// 		case "GET":
// 			api.GetStrokesHistoryByWhiteboard(w, r) // Handle fetching stroke history
// 		default:
// 			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		}
// 	})

// 	// Endpoint for updating stroke status (marking strokes as deleted)
// 	mux.HandleFunc("/strokes/delete", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == "POST" { // POST to delete based on bounding box
// 			api.UpdateStrokeForDeletion(w, r)
// 		} else {
// 			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		}
// 	})

// 	// Start the server with CORS enabled
// 	fmt.Println("Starting server on :8080")
// 	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))

// }
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

// WebSocket upgrader to handle WebSocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// later will accecpt only from my domain
		// Allow all connections (adjust for production)
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var hub = Hub{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("ERROR reading message: %v", err)
			hub.unregister <- c
			break
		}
		hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("ERROR writing message: %v", err)
				return
			}
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.conn.Close()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer ws.Close()

	client := &Client{conn: ws, send: make(chan []byte, 256)}
	hub.register <- client

	go client.writePump()
	client.readPump()
}

func main() {
	//dsn: Data Source Name
	dsn := "root:@tcp(127.0.0.1:3306)/sketchive"
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Could not grab the connection:", err)
	}

	// Ping() checks if the connection is alive
	err = database.Ping()
	if err != nil {
		log.Fatal("Lost Database connection:", err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	db.SetDB(database)

	mux := http.NewServeMux()

	// WebSocket endpoint
	mux.HandleFunc("/ws", handleConnections)

	// Start the hub in a separate goroutine
	go hub.run()

	// Start the server with CORS enabled
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}
