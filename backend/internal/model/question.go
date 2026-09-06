package model

import (
	"encoding/json"
	"time"
)

type Question struct {
	ID           string          `json:"id"`
	CategoryID   string          `json:"category_id"`
	Text         string          `json:"text"`
	Choices      json.RawMessage `json:"choices"`
	CorrectIndex int             `json:"correct_index"`
	Explanation  string          `json:"explanation"`
	Difficulty   string          `json:"difficulty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
