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
	Deleted      bool      `json:"deleted"`
	MinX         float64   `json:"minX"`
	MaxX         float64   `json:"maxX"`
	MinY         float64   `json:"minY"`
	MaxY         float64   `json:"maxY"`
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

	query := `INSERT INTO strokes (whiteboard_id, owner_id, path, color, width, created_at, deleted, minX, maxX, minY, maxY)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, stroke.WhiteboardID, stroke.OwnerID, pathJSON, stroke.Color, stroke.Width, stroke.CreatedAt,
		stroke.Deleted, stroke.MinX, stroke.MaxX, stroke.MinY, stroke.MaxY)

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
	query := `SELECT id, whiteboard_id, owner_id, path, color, width, created_at, minX, maxX, minY, maxY, deleted
			FROM strokes
			WHERE whiteboard_id = ? AND deleted = false
			ORDER BY created_at ASC`

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
		if err := rows.Scan(&stroke.ID, &stroke.WhiteboardID, &stroke.OwnerID, &pathStr, &stroke.Color, &stroke.Width, &createdAtStr, &stroke.MinX, &stroke.MaxX, &stroke.MinY, &stroke.MaxY, &stroke.Deleted); err != nil {
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

// MarkStrokesDeletedByBoundingBox marks strokes as deleted based on the eraser bounding box
func MarkStrokesDeletedByBoundingBox(whiteboardID int, minX, maxX, minY, maxY float64) error {
	log.Printf("Marking strokes as deleted for WhiteboardID: %v, BoundingBox: (%f, %f, %f, %f)", whiteboardID, minX, maxX, minY, maxY)

	query := `UPDATE strokes 
              SET deleted = true 
              WHERE whiteboard_id = ? AND deleted = false
              AND minX <= ? AND maxX >= ?
              AND minY <= ? AND maxY >= ?`

	// Ensure the bounding box parameters are in the correct order
	result, err := db.Exec(query, whiteboardID, maxX, minX, maxY, minY)
	if err != nil {
		log.Println("Error marking strokes as deleted in the database:", err)
		return err
	}

	log.Printf("Successfully marked strokes as deleted, result: %v", result)
	return nil
}
