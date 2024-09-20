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

func GetWhiteboard(w http.ResponseWriter, r *http.request) {
	// later url will contain the board id
	whiteboardID := r.URL.Query().Get(id)
	if whiteboardID == "" {
		http.Error(w, "Failed to get whiteboard's ID", http.StatusInternalServerError)
		return
	}

	// need to build GetWhiteboard()
	whiteboard, err := db.GetWhiteboardById(whiteboardID)
	if err != nil {
		http.Error(w, "Failed to get whiteboard by its ID", http.StatusInternalServerError)
		return
	}
	// sending Whiteboard data as json
	json.NewEncoder(w).Encode(whiteboard)
}
