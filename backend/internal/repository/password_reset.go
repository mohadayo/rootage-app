package repository

import (
	"context"
	"database/sql"
	"time"
)

type PasswordResetToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(ctx context.Context, userID, token string, expiresAt time.Time) error {
	// 既存の未使用トークンを削除
	r.db.ExecContext(ctx, `DELETE FROM password_reset_tokens WHERE user_id = $1 AND used_at IS NULL`, userID)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	return err
}

func (r *PasswordResetRepository) FindValidToken(ctx context.Context, token string) (*PasswordResetToken, error) {
	var t PasswordResetToken
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token, expires_at, used_at, created_at
		 FROM password_reset_tokens
		 WHERE token = $1 AND used_at IS NULL AND expires_at > NOW()`,
		token,
	).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE password_reset_tokens SET used_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}
