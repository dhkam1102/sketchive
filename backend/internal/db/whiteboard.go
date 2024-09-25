package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// NOTE: Brian
// In Go (Golang), the backticks (`) are used to define struct tags,
// which provide metadata about struct fields. The json:"id" tag is telling Go’s
// encoding/json package how to marshal (convert a Go struct to JSON) and unmarshal
// (convert JSON to a Go struct) the field when encoding and decoding JSON data.

type Whiteboard struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	OwnerID      int       `json:"owner"`
	CreatedAt    time.Time `json:"created"`
	UpdatedAt    time.Time `json:"updated_at"`
	CurrentState string    `json:"data"` // Store JSON as a string
	//  for example: whiteboard.CurrentState = `{"strokes": [...], "shapes": [...]}`
}

// db should be set in main.go
var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetDB() *sql.DB {
	return db
}

func InsertWhiteboard(board *Whiteboard) error {
	query := `INSERT INTO whiteboards (name, owner_id, created_at, updated_at, current_state)
	          VALUES (?, ?, ?, ?, ?)`

	// Insert the whiteboard data into the database
	_, err := db.Exec(query, board.Name, board.OwnerID, board.CreatedAt, board.UpdatedAt, board.CurrentState)
	if err != nil {
		log.Println("Error inserting whiteboard:", err)
		return err
	}

	return nil
}

func GetWhiteboardById(id int) (*Whiteboard, error) {
	var whiteboard Whiteboard
	var createdAt, updatedAt []byte // Scan the timestamps as byte slices (strings) first

	database := GetDB()

	query := `SELECT id, name, owner_id, created_at, updated_at, current_state
			 From whiteboards WHERE id = ?`

	row := database.QueryRow(query, id)
	err := row.Scan(&whiteboard.ID, &whiteboard.Name, &whiteboard.OwnerID, &createdAt, &updatedAt, &whiteboard.CurrentState)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("whiteboard ID: %d was not found in the database", whiteboard.ID)
		}
		log.Println("Error on running SQL query: ", err)
		return nil, err
	}

	whiteboard.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		log.Println("Error parsing created_at:", err)
		return nil, err
	}

	whiteboard.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		log.Println("Error parsing updated_at:", err)
		return nil, err
	}

	return &whiteboard, nil
}

func UpdateWhiteboard(id int, whiteboard *Whiteboard) error {
	database := GetDB()
	query := `UPDATE whiteboards 
              SET name = ?, current_state = ?, updated_at = ?
              WHERE id = ?`

	_, err := database.Exec(query, whiteboard.Name, whiteboard.CurrentState, time.Now(), id)
	if err != nil {
		log.Println("Error updating whiteboard:", err)
		return err
	}
	return nil
}

func DeleteWhiteboard(id int) error {
	database := GetDB()
	query := "DELETE FROM whiteboards WHERE id = ?"

	_, err := database.Exec(query, id)
	if err != nil {
		log.Println("Error deleting whiteboard:", err)
		return err
	}

	return nil
}
