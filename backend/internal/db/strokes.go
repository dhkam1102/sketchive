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

// func getStrocksByWhiteboardID()
