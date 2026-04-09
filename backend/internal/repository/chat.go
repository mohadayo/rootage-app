package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/rootage-ses-quiz/backend/internal/model"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateSession(ctx context.Context, session *model.ChatSession) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO chat_sessions (user_id) VALUES ($1) RETURNING id, created_at`,
		session.UserID,
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *ChatRepository) AddMessage(ctx context.Context, msg *model.ChatMessage) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO chat_messages (session_id, role, content, sources) VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		msg.SessionID, msg.Role, msg.Content, msg.Sources,
	).Scan(&msg.ID, &msg.CreatedAt)
}

func (r *ChatRepository) GetSessionsByUser(ctx context.Context, userID string) ([]model.ChatSession, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, created_at FROM chat_sessions WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.ChatSession
	for rows.Next() {
		var s model.ChatSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

func (r *ChatRepository) GetMessages(ctx context.Context, sessionID string) ([]model.ChatMessage, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, session_id, role, content, sources, created_at
		 FROM chat_messages WHERE session_id = $1 ORDER BY created_at ASC`,
		sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []model.ChatMessage
	for rows.Next() {
		var m model.ChatMessage
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.Sources, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

// Suppress unused import
var _ = json.Marshal
