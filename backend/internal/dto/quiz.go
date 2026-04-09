package dto

import (
	"encoding/json"
	"time"
)

type QuizStartRequest struct {
	CategoryID string `json:"category_id"`
	Difficulty string `json:"difficulty"`
}

type QuizStartResponse struct {
	SessionID string         `json:"session_id"`
	Questions []QuizQuestion `json:"questions"`
}

type QuizQuestion struct {
	ID      string          `json:"id"`
	Text    string          `json:"text"`
	Choices json.RawMessage `json:"choices"`
}

type AnswerRequest struct {
	SessionID     string `json:"session_id"`
	QuestionID    string `json:"question_id"`
	SelectedIndex int    `json:"selected_index"`
}

type AnswerResponse struct {
	IsCorrect    bool   `json:"is_correct"`
	CorrectIndex int    `json:"correct_index"`
	Explanation  string `json:"explanation"`
}

type QuizFinishRequest struct {
	SessionID string `json:"session_id"`
}

type QuizResultResponse struct {
	SessionID  string         `json:"session_id"`
	Score      int            `json:"score"`
	Total      int            `json:"total"`
	Percentage float64        `json:"percentage"`
	Answers    []AnswerDetail `json:"answers"`
}

type AnswerDetail struct {
	QuestionID    string          `json:"question_id"`
	QuestionText  string          `json:"question_text"`
	Choices       json.RawMessage `json:"choices"`
	SelectedIndex int             `json:"selected_index"`
	CorrectIndex  int             `json:"correct_index"`
	IsCorrect     bool            `json:"is_correct"`
	Explanation   string          `json:"explanation"`
}

type QuizHistoryItem struct {
	SessionID    string     `json:"session_id"`
	CategoryID   string     `json:"category_id"`
	CategoryName string     `json:"category_name"`
	Score        int        `json:"score"`
	Total        int        `json:"total"`
	Percentage   float64    `json:"percentage"`
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
}

type QuizStatsResponse struct {
	TotalSessions  int                 `json:"total_sessions"`
	CategoryStats  []CategoryStatsItem `json:"category_stats"`
	WeakCategory   string              `json:"weak_category"`
}

type CategoryStatsItem struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAnswers int     `json:"total_answers"`
	Correct      int     `json:"correct"`
	Percentage   float64 `json:"percentage"`
}

type ReviewItem struct {
	QuestionID   string          `json:"question_id"`
	QuestionText string          `json:"question_text"`
	Choices      json.RawMessage `json:"choices"`
	CorrectIndex int             `json:"correct_index"`
	Explanation  string          `json:"explanation"`
	CategoryName string          `json:"category_name"`
}

type UserDashboardResponse struct {
	TotalSessions     int                 `json:"total_sessions"`
	TotalQuestions    int                 `json:"total_questions"`
	TotalCorrect     int                 `json:"total_correct"`
	OverallPercentage float64            `json:"overall_percentage"`
	CategoryStats    []CategoryStatsItem `json:"category_stats"`
	WeakCategory     string              `json:"weak_category"`
	DailyStats       []DailyStatsItem    `json:"daily_stats"`
	CurrentStreak    int                 `json:"current_streak"`
	Achievements     []AchievementItem   `json:"achievements"`
}

type AchievementItem struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Difficulty   string  `json:"difficulty"`
	BestScore    int     `json:"best_score"`
	Total        int     `json:"total"`
	Percentage   float64 `json:"percentage"`
	Cleared      bool    `json:"cleared"`
}

type DailyStatsItem struct {
	Date       string  `json:"date"`
	Total      int     `json:"total"`
	Correct    int     `json:"correct"`
	Percentage float64 `json:"percentage"`
}
