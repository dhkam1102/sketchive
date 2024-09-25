package api

import (
	"encoding/json"
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
	newBoard.CurrentState = "{}"

	err := db.InsertWhiteboard(&newBoard)
	if err != nil {
		http.Error(w, "Failed to insert whiteboard", http.StatusInternalServerError)
		return
	}

	// Returning the newly created whiteboard as a JSON response
	json.NewEncoder(w).Encode(newBoard)
}

func GetWhiteboard(w http.ResponseWriter, r *http.Request) {
	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		http.Error(w, "Failed to get whiteboard's ID", http.StatusBadRequest)
		return
	}

	// Convert string to int
	id, err := strconv.Atoi(whiteboardID)
	if err != nil {
		http.Error(w, "Failed to convert whiteboardID to int", http.StatusBadRequest)
		return
	}

	// Correctly pass the integer ID to the db function
	whiteboard, err := db.GetWhiteboardById(id)
	if err != nil {
		http.Error(w, "Failed to get whiteboard by its ID", http.StatusInternalServerError)
		return
	}

	// Sending Whiteboard data as JSON
	json.NewEncoder(w).Encode(whiteboard)
}

func UpdateWhiteboard(w http.ResponseWriter, r *http.Request) {
	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		http.Error(w, "Missing whiteboard ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(whiteboardID)
	if err != nil {
		http.Error(w, "Failed to convert whiteboardID to int", http.StatusBadRequest)
		return
	}

	var updatedBoard db.Whiteboard
	// Decoding the whiteboard's body from the request
	err = json.NewDecoder(r.Body).Decode(&updatedBoard)
	if err != nil {
		http.Error(w, "Invalid body request", http.StatusBadRequest)
		return
	}

	updatedBoard.UpdatedAt = time.Now()

	// Correct function call with integer ID
	err = db.UpdateWhiteboard(id, &updatedBoard)
	if err != nil {
		http.Error(w, "Failed to update the whiteboard", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedBoard)
}

func DeleteWhiteboard(w http.ResponseWriter, r *http.Request) {
	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		http.Error(w, "Missing whiteboard ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(whiteboardID)
	if err != nil {
		http.Error(w, "Failed to convert whiteboardID to int", http.StatusBadRequest)
		return
	}

	// Correctly pass integer ID to the db function
	err = db.DeleteWhiteboard(id)
	if err != nil {
		http.Error(w, "Failed to delete whiteboard", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Whiteboard deleted successfully"})
}
