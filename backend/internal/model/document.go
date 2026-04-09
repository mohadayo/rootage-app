package model

import "time"

type Document struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Filename   string    `json:"filename"`
	Content    string    `json:"-"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type DocumentChunk struct {
	ID          string    `json:"id"`
	DocumentID  string    `json:"document_id"`
	ChunkIndex  int       `json:"chunk_index"`
	Content     string    `json:"content"`
	Embedding   []float64 `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
}
