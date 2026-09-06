package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rootage-ses-quiz/backend/internal/model"
)

type QuestionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) GetRandomByCategory(ctx context.Context, categoryID string, limit int) ([]model.Question, error) {
	return r.GetRandomByCategoryAndDifficulty(ctx, categoryID, "", limit)
}

func (r *QuestionRepository) GetRandomByCategoryAndDifficulty(ctx context.Context, categoryID, difficulty string, limit int) ([]model.Question, error) {
	var rows *sql.Rows
	var err error

	if difficulty != "" {
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, category_id, text, choices, correct_index, explanation, difficulty, created_at, updated_at
			 FROM questions WHERE category_id = $1 AND difficulty = $2 AND is_active ORDER BY RANDOM() LIMIT $3`,
			categoryID, difficulty, limit,
		)
	} else {
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, category_id, text, choices, correct_index, explanation, difficulty, created_at, updated_at
			 FROM questions WHERE category_id = $1 AND is_active ORDER BY RANDOM() LIMIT $2`,
			categoryID, limit,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []model.Question
	for rows.Next() {
		var q model.Question
		if err := rows.Scan(&q.ID, &q.CategoryID, &q.Text, &q.Choices, &q.CorrectIndex, &q.Explanation, &q.Difficulty, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

func (r *QuestionRepository) GetByID(ctx context.Context, id string) (*model.Question, error) {
	var q model.Question
	err := r.db.QueryRowContext(ctx,
		`SELECT id, category_id, text, choices, correct_index, explanation, difficulty, created_at, updated_at
		 FROM questions WHERE id = $1`, id,
	).Scan(&q.ID, &q.CategoryID, &q.Text, &q.Choices, &q.CorrectIndex, &q.Explanation, &q.Difficulty, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *QuestionRepository) GetByIDs(ctx context.Context, ids []string) (map[string]*model.Question, error) {
	if len(ids) == 0 {
		return map[string]*model.Question{}, nil
	}
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(
		`SELECT id, category_id, text, choices, correct_index, explanation, difficulty, created_at, updated_at
		 FROM questions WHERE id IN (%s)`,
		strings.Join(placeholders, ","),
	)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*model.Question)
	for rows.Next() {
		var q model.Question
		if err := rows.Scan(&q.ID, &q.CategoryID, &q.Text, &q.Choices, &q.CorrectIndex, &q.Explanation, &q.Difficulty, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, err
		}
		result[q.ID] = &q
	}
	return result, rows.Err()
}

func (r *QuestionRepository) List(ctx context.Context, categoryID string) ([]model.Question, error) {
	query := `SELECT id, category_id, text, choices, correct_index, explanation, difficulty, created_at, updated_at FROM questions`
	var rows *sql.Rows
	var err error

	if categoryID != "" {
		query += ` WHERE category_id = $1 AND is_active ORDER BY created_at DESC`
		rows, err = r.db.QueryContext(ctx, query, categoryID)
	} else {
		query += ` WHERE is_active ORDER BY created_at DESC`
		rows, err = r.db.QueryContext(ctx, query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []model.Question
	for rows.Next() {
		var q model.Question
		if err := rows.Scan(&q.ID, &q.CategoryID, &q.Text, &q.Choices, &q.CorrectIndex, &q.Explanation, &q.Difficulty, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

func (r *QuestionRepository) Create(ctx context.Context, q *model.Question) error {
	if q.Difficulty == "" {
		q.Difficulty = "beginner"
	}
	return r.db.QueryRowContext(ctx,
		`INSERT INTO questions (category_id, text, choices, correct_index, explanation, difficulty)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		q.CategoryID, q.Text, q.Choices, q.CorrectIndex, q.Explanation, q.Difficulty,
	).Scan(&q.ID, &q.CreatedAt, &q.UpdatedAt)
}

func (r *QuestionRepository) Update(ctx context.Context, q *model.Question) error {
	if q.Difficulty == "" {
		q.Difficulty = "beginner"
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE questions SET text = $1, choices = $2, correct_index = $3, explanation = $4, difficulty = $5, updated_at = NOW()
		 WHERE id = $6`,
		q.Text, q.Choices, q.CorrectIndex, q.Explanation, q.Difficulty, q.ID,
	)
	return err
}

// Deactivate は問題を論理削除する（出題・一覧から外すが回答履歴は残す）。
func (r *QuestionRepository) Deactivate(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE questions SET is_active = FALSE, updated_at = NOW() WHERE id = $1`, id)
	return err
}
