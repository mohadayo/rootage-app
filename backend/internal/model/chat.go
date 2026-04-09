package model

import (
	"encoding/json"
	"time"
)

type ChatSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMessage struct {
	ID        string          `json:"id"`
	SessionID string          `json:"session_id"`
	Role      string          `json:"role"`
	Content   string          `json:"content"`
	Sources   json.RawMessage `json:"sources,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type Source struct {
	DocumentID    string `json:"document_id"`
	DocumentTitle string `json:"document_title"`
	ChunkIndex    int    `json:"chunk_index"`
	Content       string `json:"content"`
}
