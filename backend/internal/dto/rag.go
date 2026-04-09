package dto

import "time"

type RAGAskRequest struct {
	Question string              `json:"question"`
	History  []RAGHistoryMessage `json:"history,omitempty"`
}

type RAGHistoryMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RAGResponse struct {
	Answer  string      `json:"answer"`
	Sources []RAGSource `json:"sources"`
}

type RAGSource struct {
	DocumentID    string `json:"document_id"`
	DocumentTitle string `json:"document_title"`
	ChunkIndex    int    `json:"chunk_index"`
	Content       string `json:"content"`
}

type ChatHistoryItem struct {
	SessionID string          `json:"session_id"`
	Messages  []ChatMessageDTO `json:"messages"`
	CreatedAt time.Time       `json:"created_at"`
}

type ChatMessageDTO struct {
	Role      string      `json:"role"`
	Content   string      `json:"content"`
	Sources   []RAGSource `json:"sources,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}
