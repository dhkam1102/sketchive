package db

import (
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
	OwnwerID     int       `json:"owner"`
	CreatedAt    time.Time `json:"created"`
	UpdatedAt    time.Time `json:"updated_at"`
	CurrentState string    `json:"data"` // Store JSON as a string
	//  for example: whiteboard.CurrentState = `{"strokes": [...], "shapes": [...]}`
}
