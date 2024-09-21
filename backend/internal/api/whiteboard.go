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
	whiteboardID := r.URL.Query().Get("id")
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

func UpdateWhitebaord(w http.ResponseWriter, r *http.Request) {
	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		http.Error(w, "Missing whiteboard ID", http.StatusInternalServerError)
		return
	}

	var updatedBoard db.whiteboard
	// need more studing on how the decoder works (PUT, PATCH requests)
	// decoding the whiteboard's body from the request (which contains the new data)
	err := json.NewDecoder(r.Body).Decode(&updatedBoard)
	if err != nil {
		http.Error(w, "invalid body request")
		return
	}

	updatedBoard.UpdatedAt = time.Now()

	error = db.updateWhiteboard(whiteboardID, &updatedBoard)
	if error != nil {
		http.Error(w, "failed to upate the whitebaord", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedBoard)
}

fun DeleteWhiteboard(w http.ResponseWriter, r *http.Request) {
	whiteboardID := r.URL.Query().Get("id")
	if whiteboardID == "" {
		http.Error(w, "Missing whiteboard ID", http.StatusInternalServerError)
		return
	}

	err := db.DeleteWhiteboardByID(whiteboardID)
    if err != nil {
        http.Error(w, "Failed to delete whiteboard", http.StatusInternalServerError)
        return
    }

	json.NewEncoder(w).Encode(map[string]string{"message": "whiteboard deleted successfully"})
}

//func AddStroke()

// func GetStokesHistoryByWhiteboard()