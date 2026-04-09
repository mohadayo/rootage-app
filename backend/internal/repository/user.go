package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/rootage-ses-quiz/backend/internal/model"
)

type UserSummary struct {
	ID             string
	Name           string
	Email          string
	CreatedAt      time.Time
	LastQuizAt     *time.Time
	TotalCorrect   int
	TotalQuestions  int
}

type UserCategoryStats struct {
	UserID       string
	CategoryID   string
	CategoryName string
	TotalAnswers int
	Correct      int
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, name, role) VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at, updated_at`,
		user.Email, user.PasswordHash, user.Name, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`,
		passwordHash, id,
	)
	return err
}

func (r *UserRepository) ListUserSummaries(ctx context.Context) ([]UserSummary, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT u.id, u.name, u.email, u.created_at,
		        MAX(qs.finished_at) AS last_quiz_at,
		        COALESCE(SUM(qs.score), 0) AS total_correct,
		        COALESCE(SUM(qs.total), 0) AS total_questions
		 FROM users u
		 LEFT JOIN quiz_sessions qs ON qs.user_id = u.id AND qs.finished_at IS NOT NULL
		 WHERE u.role = 'user'
		 GROUP BY u.id, u.name, u.email, u.created_at
		 ORDER BY u.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []UserSummary
	for rows.Next() {
		var s UserSummary
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.CreatedAt, &s.LastQuizAt, &s.TotalCorrect, &s.TotalQuestions); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}

func (r *UserRepository) ListAllCategoryStats(ctx context.Context) ([]UserCategoryStats, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT qs.user_id, c.id AS category_id, c.name AS category_name,
		        COUNT(qa.id) AS total_answers,
		        COUNT(CASE WHEN qa.is_correct THEN 1 END) AS correct
		 FROM quiz_sessions qs
		 JOIN quiz_answers qa ON qa.session_id = qs.id
		 JOIN questions q ON qa.question_id = q.id
		 JOIN categories c ON q.category_id = c.id
		 WHERE qs.finished_at IS NOT NULL
		 GROUP BY qs.user_id, c.id, c.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []UserCategoryStats
	for rows.Next() {
		var s UserCategoryStats
		if err := rows.Scan(&s.UserID, &s.CategoryID, &s.CategoryName, &s.TotalAnswers, &s.Correct); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}
