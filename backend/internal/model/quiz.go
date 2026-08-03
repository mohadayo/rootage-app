package model

import (
	"encoding/json"
	"time"
)

type QuizSession struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	CategoryID string     `json:"category_id"`
	Difficulty string     `json:"difficulty"`
	Score      int        `json:"score"`
	Total      int        `json:"total"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	// このセッションで出題した問題のID。回答時の検証に使うためクライアントには返さない
	QuestionIDs json.RawMessage `json:"-"`
}

type QuizAnswer struct {
	ID            string `json:"id"`
	SessionID     string `json:"session_id"`
	QuestionID    string `json:"question_id"`
	SelectedIndex int    `json:"selected_index"`
	IsCorrect     bool   `json:"is_correct"`
}

type QuizAnswerDetail struct {
	QuizAnswer
	QuestionText string          `json:"question_text"`
	Choices      json.RawMessage `json:"choices"`
	CorrectIndex int             `json:"correct_index"`
	Explanation  string          `json:"explanation"`
	CategoryName string          `json:"category_name"`
}
