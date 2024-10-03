package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sketchive/internal/db"
	"strconv"
	"time"
)

func CreateWhiteboard(w http.ResponseWriter, r *http.Request) {
	var newBoard db.Whiteboard
	newBoard.Name = "Untitled"
	newBoard.OwnerID = 1 // default for right now, use actual user data next time
	newBoard.CreatedAt = time.Now()
	newBoard.UpdatedAt = newBoard.CreatedAt

	err := db.InsertWhiteboard(&newBoard)
	if err != nil {
		log.Println("Error inserting whiteboard: ", err)
		http.Error(w, "Failed to insert whiteboard", http.StatusInternalServerError)
		return
	}

	// Returning the newly created whiteboard as a JSON response
	json.NewEncoder(w).Encode(newBoard)
}

func GetWhiteboard(w http.ResponseWriter, r *http.Request) {
	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		log.Println("Error: missing whiteboard ID in the request (GetWhiteboard()")
		http.Error(w, "Failed to get whiteboard's ID", http.StatusBadRequest)
		return
	}

	// Convert string to int
	id, err := strconv.Atoi(whiteboardID)
	if err != nil {
		log.Println("Error converting whiteboard ID to int (GetWhiteboard())")
		http.Error(w, "Failed to convert whiteboardID to int", http.StatusBadRequest)
		return
	}

	// Correctly pass the integer ID to the db function
	whiteboard, err := db.GetWhiteboardById(id)
	if err != nil {
		log.Println("Error fetching whiteboard by ID (GetWhiteboard()")
		http.Error(w, "Failed to get whiteboard by its ID", http.StatusInternalServerError)
		return
	}

	// Sending Whiteboard data as JSON
	json.NewEncoder(w).Encode(whiteboard)
}

func UpdateWhiteboard(w http.ResponseWriter, r *http.Request) {
	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		log.Println("Error: missing whiteboard ID in request (UpdateWhiteboard())")
		http.Error(w, "Missing whiteboard ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(whiteboardID)
	if err != nil {
		log.Println("Error converting whiteboard ID to int (UpdateWhiteboard()):", err)
		http.Error(w, "Failed to convert whiteboardID to int", http.StatusBadRequest)
		return
	}

	var updatedBoard db.Whiteboard
	// Decoding the whiteboard's body from the request
	err = json.NewDecoder(r.Body).Decode(&updatedBoard)
	if err != nil {
		log.Println("Error decoding request body (UpdateWhiteboard()):", err)
		http.Error(w, "Invalid body request", http.StatusBadRequest)
		return
	}

	updatedBoard.UpdatedAt = time.Now()
	// Correct function call with integer ID
	err = db.UpdateWhiteboard(id, &updatedBoard)
	if err != nil {
		log.Println("Error updating whiteboard (UpdateWhiteboard()):", err)
		http.Error(w, "Failed to update the whiteboard", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedBoard)
}

func DeleteWhiteboard(w http.ResponseWriter, r *http.Request) {
	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		log.Println("Error: missing whiteboard ID in request (DeleteWhiteboard())")
		http.Error(w, "Missing whiteboard ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(whiteboardID)
	if err != nil {
		log.Println("Error converting whiteboard ID to int (DeleteWhiteboard()):", err)
		http.Error(w, "Failed to convert whiteboardID to int", http.StatusBadRequest)
		return
	}

	// Correctly pass integer ID to the db function
	err = db.DeleteWhiteboard(id)
	if err != nil {
		log.Println("Error deleting whiteboard (DeleteWhiteboard()):", err)
		http.Error(w, "Failed to delete whiteboard", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Whiteboard deleted successfully"})
}

func ClearWhiteboardHandler(w http.ResponseWriter, r *http.Request) {
	whiteboardIDStr := r.URL.Query().Get("id")
	if whiteboardIDStr == "" {
		log.Println("Error: missing whiteboard ID in request (ClearWhiteboardHandler())")
		http.Error(w, "Missing whiteboard ID", http.StatusBadRequest)
		return
	}

	whiteboardID, err := strconv.Atoi(whiteboardIDStr)
	if err != nil {
		log.Println("Error converting whiteboard ID to int (ClearWhiteboardHandler()) :", err)
		http.Error(w, "Invalid whiteboard ID", http.StatusBadRequest)
		return
	}

	// Call the DB function to clear strokes for the whiteboard
	err = db.ClearStrokesByWhiteboardID(whiteboardID)
	if err != nil {
		log.Println("Error deleting whiteboard (ClearWhiteboardHandler()):", err)
		http.Error(w, "Failed to clear strokes", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully cleared strokes for whiteboard ID %d\n", whiteboardID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Whiteboard cleared successfully"})
}
