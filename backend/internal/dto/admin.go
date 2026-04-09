package dto

import (
	"encoding/json"
	"time"
)

type CreateQuestionRequest struct {
	CategoryID   string          `json:"category_id"`
	Text         string          `json:"text"`
	Choices      json.RawMessage `json:"choices"`
	CorrectIndex int             `json:"correct_index"`
	Explanation  string          `json:"explanation"`
	Difficulty   string          `json:"difficulty"`
}

type UpdateQuestionRequest struct {
	Text         string          `json:"text"`
	Choices      json.RawMessage `json:"choices"`
	CorrectIndex int             `json:"correct_index"`
	Explanation  string          `json:"explanation"`
	Difficulty   string          `json:"difficulty"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateDocumentTextRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ImportError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

type ImportQuestionsResponse struct {
	Imported int           `json:"imported"`
	Errors   []ImportError `json:"errors"`
}

// Guide

type CreateGuideCategoryRequest struct {
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

type UpdateGuideCategoryRequest struct {
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

type CreateGuideRequest struct {
	CategoryID  string `json:"category_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	IsPublished bool   `json:"is_published"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateGuideRequest struct {
	CategoryID  string `json:"category_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	IsPublished bool   `json:"is_published"`
	SortOrder   int    `json:"sort_order"`
}

// Admin Summary

type AdminSummary struct {
	TotalUsers     int `json:"total_users"`
	TotalQuestions int `json:"total_questions"`
	TotalDocuments int `json:"total_documents"`
	TotalGuides    int `json:"total_guides"`
}

// User Progress

type UserProgressItem struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	Email          string                  `json:"email"`
	CreatedAt      time.Time               `json:"created_at"`
	LastQuizAt     *time.Time              `json:"last_quiz_at"`
	TotalCorrect   int                     `json:"total_correct"`
	TotalQuestions int                     `json:"total_questions"`
	Percentage     float64                 `json:"percentage"`
	CategoryStats  []UserCategoryStatsItem `json:"category_stats"`
}

type UserCategoryStatsItem struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAnswers int     `json:"total_answers"`
	Correct      int     `json:"correct"`
	Percentage   float64 `json:"percentage"`
}
