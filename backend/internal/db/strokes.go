package db

import (
	"encoding/json"
	"log"
	"time"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Stroke struct {
	ID           int       `json:"id"`
	WhiteboardID int       `json:"whiteboardID"` // Use camel case to match the frontend
	OwnerID      int       `json:"ownerID"`
	Path         []Point   `json:"path"`
	Color        string    `json:"color"`
	Width        int       `json:"width"`
	CreatedAt    time.Time `json:"created_at"`
}

// InsertStroke inserts a stroke into the strokes table and logs the process
func InsertStroke(stroke *Stroke) error {
	log.Println("Inserting new stroke:", stroke)
	log.Printf("Inserting stroke with WhiteboardID: %v", stroke.WhiteboardID)

	// Convert the Path to JSON
	pathJSON, err := json.Marshal(stroke.Path)
	if err != nil {
		log.Println("Error marshaling stroke path:", err)
		return err
	}

	query := `INSERT INTO strokes (whiteboard_id, owner_id, path, color, width, created_at)
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, stroke.WhiteboardID, stroke.OwnerID, pathJSON, stroke.Color, stroke.Width, stroke.CreatedAt)
	if err != nil {
		log.Println("Error inserting stroke into database:", err)
		return err
	}

	log.Println("Stroke inserted successfully, result:", result)
	return nil
}

func GetStrokesByWhiteboardID(whiteboardID int) ([]Stroke, error) {
	log.Printf("Fetching strokes for WhiteboardID: %v", whiteboardID)

	var strokes []Stroke
	query := `SELECT id, whiteboard_id, owner_id, path, color, width, created_at
              FROM strokes WHERE whiteboard_id = ? ORDER BY created_at ASC`

	rows, err := db.Query(query, whiteboardID)
	if err != nil {
		log.Println("Error fetching strokes from database:", err)
		return nil, err
	}
	defer rows.Close()

	log.Printf("Processing rows for WhiteboardID: %v", whiteboardID)
	for rows.Next() {
		var stroke Stroke
		var pathStr string      // Temporarily store path as string
		var createdAtStr string // Temporarily store created_at as string

		// Scan into appropriate types
		if err := rows.Scan(&stroke.ID, &stroke.WhiteboardID, &stroke.OwnerID, &pathStr, &stroke.Color, &stroke.Width, &createdAtStr); err != nil {
			log.Println("Error scanning stroke data:", err)
			return nil, err
		}

		// Parse path JSON
		if err := json.Unmarshal([]byte(pathStr), &stroke.Path); err != nil {
			log.Println("Error unmarshaling stroke path:", err)
			return nil, err
		}

		// Parse created_at to time.Time
		stroke.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Println("Error parsing created_at:", err)
			return nil, err
		}

		strokes = append(strokes, stroke)
	}

	log.Printf("Successfully fetched %d strokes for WhiteboardID %v", len(strokes), whiteboardID)
	return strokes, nil
}
