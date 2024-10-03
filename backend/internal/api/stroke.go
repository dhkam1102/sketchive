package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sketchive/internal/db"
	"strconv"
	"time"
)

// probably need modification
func calculateBoundingBox(points []db.Point) (float64, float64, float64, float64, error) {
	if len(points) == 0 {
		return 0, 0, 0, 0, fmt.Errorf("points slice is empty, cannot calculate bounding box")
	}
	minX, maxX := points[0].X, points[0].X
	minY, maxY := points[0].Y, points[0].Y
	for _, point := range points {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}
	return minX, maxX, minY, maxY, nil
}

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

	// Calculate bounding box for the stroke
	minX, maxX, minY, maxY, err := calculateBoundingBox(newStroke.Path)
	if err != nil {
		log.Println("Error calculating bounding box:", err)
		http.Error(w, "Failed to calculate bounding box", http.StatusBadRequest)
		return
	}
	newStroke.MinX = minX
	newStroke.MaxX = maxX
	newStroke.MinY = minY
	newStroke.MaxY = maxY

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

	// Filter out deleted strokes
	activeStrokes := []db.Stroke{}
	for _, stroke := range strokes {
		if !stroke.Deleted {
			activeStrokes = append(activeStrokes, stroke)
		}
	}

	log.Printf("Successfully retrieved %d strokes for whiteboard ID %d\n", len(strokes), id)
	json.NewEncoder(w).Encode(strokes)
}

// UpdateStrokeForDeletion marks strokes as deleted based on the eraser bounding box
func UpdateStrokeForDeletion(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateStrokeForDeletion API called")

	var eraserBox struct {
		WhiteboardID int     `json:"whiteboardID"`
		MinX         float64 `json:"minX"`
		MaxX         float64 `json:"maxX"`
		MinY         float64 `json:"minY"`
		MaxY         float64 `json:"maxY"`
	}

	err := json.NewDecoder(r.Body).Decode(&eraserBox)
	if err != nil {
		log.Println("Error decoding eraser bounding box data:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Marking strokes for deletion in whiteboard ID %d with bounding box (%f, %f, %f, %f)\n",
		eraserBox.WhiteboardID, eraserBox.MinX, eraserBox.MaxX, eraserBox.MinY, eraserBox.MaxY)

	// Update strokes in the database that fall within the eraser bounding box
	err = db.MarkStrokesDeletedByBoundingBox(eraserBox.WhiteboardID, eraserBox.MinX, eraserBox.MaxX, eraserBox.MinY, eraserBox.MaxY)
	if err != nil {
		log.Println("Error marking strokes as deleted:", err)
		http.Error(w, "Failed to mark strokes as deleted", http.StatusInternalServerError)
		return
	}

	log.Println("Strokes marked as deleted successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Strokes marked as deleted successfully"})
}
