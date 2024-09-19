package api

import (
	"encoding/json"
	"net/http"
	"sketchive/internal/db"
	"time"
)

func CreateWhiteboard(w http.ResponseWriter, r *http.Request) {
	var newBoard db.Whiteboard
	newBoard.Name = "Untitled"
	newBoard.OwnerID = 0 // default for right now, use actual user data next time
	newBoard.CreatedAt = time.Now()
	newBoard.UpdatedAt = newBoard.CreatedAt
	newBoard.CurrentState = ""

	err := db.InsertWhiteboard(&newBoard)
	if err != nil {
		http.Error(w, "Failed to insert whiteboard", http.StatusInternalServerError)
		return
	}

	// Returning the newly created whiteboard as a JSON response
	json.NewEncoder(w).Encode(newBoard)

}
