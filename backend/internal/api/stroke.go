package api

import (
	"encoding/json"
	"net/http"
	"sketchive/internal/db"
	"strconv"
	"time"
)

func AddStroke(w http.ResponseWriter, r *http.Request) {
	var newStroke db.Stroke

	err := json.NewDecoder(r.Body).Decode(&newStroke)
	if err != nil {
		http.Error(w, "Error when decoding stroke", http.StatusBadRequest)
		return
	}

	newStroke.CreatedAt = time.Now()

	err = db.InsertStroke(&newStroke)
	if err != nil {
		http.Error(w, "Error when inserting stroke", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newStroke)
}

func GetStrokesHistoryByWhiteboard(w http.ResponseWriter, r *http.Request) {
	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		http.Error(w, "Failed to get whiteboard's ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(whiteboardID)
	if err != nil {
		http.Error(w, "Can't convert whiteboardID to int", http.StatusBadRequest)
		return
	}

	strokes, err := db.GetStrokesByWhiteboardID(id)
	if err != nil {
		http.Error(w, "Failed to retrieve strokes history", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(strokes)
}
