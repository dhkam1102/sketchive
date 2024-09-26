package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sketchive/internal/db"
	"strconv"
	"time"
)

// AddStroke adds a new stroke to the database and logs relevant details
func AddStroke(w http.ResponseWriter, r *http.Request) {
	log.Println("AddStroke API called")

	var newStroke db.Stroke
	err := json.NewDecoder(r.Body).Decode(&newStroke)
	if err != nil {
		log.Println("Error decoding stroke data:", err)
		http.Error(w, "Error decoding stroke", http.StatusBadRequest)
		return
	}

	newStroke.CreatedAt = time.Now()
	log.Printf("Decoded stroke data: %+v\n", newStroke)

	err = db.InsertStroke(&newStroke)
	if err != nil {
		log.Println("Error inserting stroke into database:", err)
		http.Error(w, "Error inserting stroke", http.StatusInternalServerError)
		return
	}

	log.Println("Stroke inserted successfully")
	json.NewEncoder(w).Encode(newStroke)
}

// GetStrokesHistoryByWhiteboard retrieves stroke history for a specific whiteboard
func GetStrokesHistoryByWhiteboard(w http.ResponseWriter, r *http.Request) {
	log.Println("GetStrokesHistoryByWhiteboard API called")

	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		log.Println("Error: missing whiteboard ID in request")
		http.Error(w, "Failed to get whiteboard's ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(whiteboardID)
	if err != nil {
		log.Println("Error converting whiteboard ID to int:", err)
		http.Error(w, "Can't convert whiteboardID to int", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching stroke history for whiteboard ID: %d\n", id)
	strokes, err := db.GetStrokesByWhiteboardID(id)
	if err != nil {
		log.Println("Error retrieving stroke history from database:", err)
		http.Error(w, "Failed to retrieve strokes history", http.StatusInternalServerError)
		return
	}

	// Ensure strokes is not null, return an empty array if no strokes found
	if strokes == nil {
		log.Println("No strokes found for whiteboard ID", id)
		strokes = []db.Stroke{} // Ensure an empty array is returned
	}

	log.Printf("Successfully retrieved %d strokes for whiteboard ID %d\n", len(strokes), id)
	json.NewEncoder(w).Encode(strokes)
}
