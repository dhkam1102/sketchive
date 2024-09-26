package db

import (
	"log"
	"time"
)

// Stroke struct to hold stroke data
type Stroke struct {
	ID           int       `json:"id"`
	WhiteboardID int       `json:"whiteboard_id"` // Foreign key to whiteboard
	OwnerID      int       `json:"owner_id"`      // Who created this stroke
	XStart       float64   `json:"x_start"`       // Starting x coordinate of the stroke
	YStart       float64   `json:"y_start"`
	XEnd         float64   `json:"x_end"`
	YEnd         float64   `json:"y_end"`
	Color        string    `json:"color"` // Stroke color
	Width        int       `json:"width"` // Stroke width
	CreatedAt    time.Time `json:"created_at"`
}

// InsertStroke inserts a stroke into the strokes table and logs the process
func InsertStroke(stroke *Stroke) error {
	log.Println("Inserting new stroke:", stroke)

	query := `INSERT INTO strokes (whiteboard_id, owner_id, x_start, y_start, x_end, y_end, color, width, created_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, stroke.WhiteboardID, stroke.OwnerID, stroke.XStart, stroke.YStart, stroke.XEnd, stroke.YEnd, stroke.Color, stroke.Width, stroke.CreatedAt)
	if err != nil {
		log.Println("Error inserting stroke into database:", err)
		return err
	}

	log.Println("Stroke inserted successfully with ID:", stroke.ID)
	return nil
}

// GetStrokesByWhiteboardID fetches all strokes for a specific whiteboard and logs the process
func GetStrokesByWhiteboardID(whiteboardID int) ([]Stroke, error) {
	log.Printf("Fetching strokes for whiteboard ID: %d\n", whiteboardID)

	var strokes []Stroke
	query := `SELECT id, whiteboard_id, owner_id, x_start, y_start, x_end, y_end, color, width, created_at
              FROM strokes WHERE whiteboard_id = ? ORDER BY created_at ASC`

	rows, err := db.Query(query, whiteboardID)
	if err != nil {
		log.Println("Error fetching strokes from database:", err)
		return nil, err
	}
	defer rows.Close()

	log.Printf("Processing rows for whiteboard ID: %d\n", whiteboardID)
	for rows.Next() {
		var stroke Stroke
		var createdAtStr string // Temporary variable to hold created_at as string

		// Scan created_at as a string, then convert it to time.Time
		if err := rows.Scan(&stroke.ID, &stroke.WhiteboardID, &stroke.OwnerID, &stroke.XStart, &stroke.YStart, &stroke.XEnd, &stroke.YEnd, &stroke.Color, &stroke.Width, &createdAtStr); err != nil {
			log.Println("Error scanning stroke data:", err)
			return nil, err
		}

		// Parse created_at string into time.Time
		stroke.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Println("Error parsing created_at:", err)
			return nil, err
		}

		log.Printf("Fetched stroke: %+v\n", stroke)
		strokes = append(strokes, stroke)
	}

	log.Printf("Successfully fetched %d strokes for whiteboard ID %d\n", len(strokes), whiteboardID)
	return strokes, nil
}
