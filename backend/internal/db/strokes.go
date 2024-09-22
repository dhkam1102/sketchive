package db

import (
	"log"
	"time"
)

type Stroke struct {
	ID           int       `json:"id"`
	WhiteboardID int       `json:"whiteboard_id"`
	UserID       int       `json:"user_id"`
	StrokeData   string    `json:"stroke_data"`
	CreatedAt    time.Time `json:"created_at"`
}

func InsertStroke(stoke *Stroke) error {
	query := `INSERT INTO strokes (whiteboard_id, user_id, stroke_data, created_at)
			VALUES (?, ? ?, ?)`

	database := GetDB()
	_, err := database.Exec(query, stroke.WhiteboardID, stroke.UserID, stroke.StrokeData, stroke.CreatedAt)
	if err != nil {
		log.Println("Error inserting stroke:", err)
		return err
	}

	return nil
}

func getStrocksByWhiteboardID(WhiteboardID string) {
	var strokes []Stroke
	database := GetDB()

	query := `SELECT id, whiteboard_id, user_id, stroke_data, created_at
	FROM strokes WHERE whiteboard_id = ? ORDER BY created_at`

	rows, err := database.Query(qeury, whitebaordID)
	if err != nil {
		log.Println("Error fetching strokes:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stroke Stroke
		if err := rows.Scan(&stroke.ID, &stroke.WhiteboardID, &stroke.UserID, &stroke.StrokeData, &stroke.CreatedAt); err != nil {
			log.Println("Error scanning stroke data:", err)
			return nil, err
		}
		strokes = append(strokes, stroke)
	}

	return strokes, nil

}
