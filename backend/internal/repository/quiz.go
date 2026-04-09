package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/rootage-ses-quiz/backend/internal/model"
)

type QuizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) CreateSession(ctx context.Context, session *model.QuizSession) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO quiz_sessions (user_id, category_id, difficulty, total) VALUES ($1, $2, $3, $4)
		 RETURNING id, started_at`,
		session.UserID, session.CategoryID, session.Difficulty, session.Total,
	).Scan(&session.ID, &session.StartedAt)
}

func (r *QuizRepository) AddAnswer(ctx context.Context, answer *model.QuizAnswer) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO quiz_answers (session_id, question_id, selected_index, is_correct)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		answer.SessionID, answer.QuestionID, answer.SelectedIndex, answer.IsCorrect,
	).Scan(&answer.ID)
}

func (r *QuizRepository) FinishSession(ctx context.Context, sessionID string, score int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE quiz_sessions SET score = $1, finished_at = NOW() WHERE id = $2`,
		score, sessionID,
	)
	return err
}

func (r *QuizRepository) GetSession(ctx context.Context, sessionID string) (*model.QuizSession, error) {
	var s model.QuizSession
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, category_id, difficulty, score, total, started_at, finished_at
		 FROM quiz_sessions WHERE id = $1`, sessionID,
	).Scan(&s.ID, &s.UserID, &s.CategoryID, &s.Difficulty, &s.Score, &s.Total, &s.StartedAt, &s.FinishedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *QuizRepository) GetSessionAnswers(ctx context.Context, sessionID string) ([]model.QuizAnswer, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, session_id, question_id, selected_index, is_correct
		 FROM quiz_answers WHERE session_id = $1`, sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []model.QuizAnswer
	for rows.Next() {
		var a model.QuizAnswer
		if err := rows.Scan(&a.ID, &a.SessionID, &a.QuestionID, &a.SelectedIndex, &a.IsCorrect); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, rows.Err()
}

func (r *QuizRepository) GetSessionsByUser(ctx context.Context, userID string) ([]model.QuizSession, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, category_id, difficulty, score, total, started_at, finished_at
		 FROM quiz_sessions WHERE user_id = $1 ORDER BY started_at DESC LIMIT 100`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.QuizSession
	for rows.Next() {
		var s model.QuizSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.CategoryID, &s.Difficulty, &s.Score, &s.Total, &s.StartedAt, &s.FinishedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

func (r *QuizRepository) GetWrongAnswers(ctx context.Context, userID string) ([]model.QuizAnswerDetail, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT qa.id, qa.session_id, qa.question_id, qa.selected_index, qa.is_correct,
		        q.text, q.choices, q.correct_index, q.explanation, c.name
		 FROM quiz_answers qa
		 JOIN quiz_sessions qs ON qa.session_id = qs.id
		 JOIN questions q ON qa.question_id = q.id
		 JOIN categories c ON q.category_id = c.id
		 WHERE qs.user_id = $1 AND qa.is_correct = false
		 ORDER BY qs.started_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []model.QuizAnswerDetail
	for rows.Next() {
		var d model.QuizAnswerDetail
		if err := rows.Scan(&d.ID, &d.SessionID, &d.QuestionID, &d.SelectedIndex, &d.IsCorrect,
			&d.QuestionText, &d.Choices, &d.CorrectIndex, &d.Explanation, &d.CategoryName); err != nil {
			return nil, err
		}
		details = append(details, d)
	}
	return details, rows.Err()
}

type CategoryStats struct {
	CategoryID   string
	CategoryName string
	TotalAnswers int
	Correct      int
}

func (r *QuizRepository) GetCategoryStats(ctx context.Context, userID string) ([]CategoryStats, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.id, c.name,
		        COUNT(qa.id) as total_answers,
		        COUNT(CASE WHEN qa.is_correct THEN 1 END) as correct
		 FROM quiz_answers qa
		 JOIN quiz_sessions qs ON qa.session_id = qs.id
		 JOIN questions q ON qa.question_id = q.id
		 JOIN categories c ON q.category_id = c.id
		 WHERE qs.user_id = $1 AND qs.finished_at IS NOT NULL
		 GROUP BY c.id, c.name
		 ORDER BY c.name`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []CategoryStats
	for rows.Next() {
		var s CategoryStats
		if err := rows.Scan(&s.CategoryID, &s.CategoryName, &s.TotalAnswers, &s.Correct); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

func (r *QuizRepository) GetTotalSessions(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM quiz_sessions WHERE user_id = $1 AND finished_at IS NOT NULL`, userID,
	).Scan(&count)
	return count, err
}

func (r *QuizRepository) GetTotalQuestionsAndCorrect(ctx context.Context, userID string) (int, int, error) {
	var total, correct int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(qa.id), COUNT(CASE WHEN qa.is_correct THEN 1 END)
		 FROM quiz_answers qa
		 JOIN quiz_sessions qs ON qa.session_id = qs.id
		 WHERE qs.user_id = $1 AND qs.finished_at IS NOT NULL`, userID,
	).Scan(&total, &correct)
	return total, correct, err
}

func (r *QuizRepository) GetDailyStats(ctx context.Context, userID string) ([]DailyStats, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DATE(qs.finished_at) as date,
		        COUNT(qa.id) as total,
		        COUNT(CASE WHEN qa.is_correct THEN 1 END) as correct
		 FROM quiz_answers qa
		 JOIN quiz_sessions qs ON qa.session_id = qs.id
		 WHERE qs.user_id = $1 AND qs.finished_at IS NOT NULL
		 GROUP BY DATE(qs.finished_at)
		 ORDER BY date`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []DailyStats
	for rows.Next() {
		var s DailyStats
		if err := rows.Scan(&s.Date, &s.Total, &s.Correct); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

type BestScore struct {
	CategoryID   string
	CategoryName string
	Difficulty   string
	BestScore    int
	Total        int
	Percentage   float64
}

func (r *QuizRepository) GetBestScores(ctx context.Context, userID string) ([]BestScore, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT qs.category_id, c.name, qs.difficulty,
		        MAX(qs.score) as best_score,
		        MAX(qs.total) as total,
		        MAX(CASE WHEN qs.total > 0 THEN qs.score * 100.0 / qs.total ELSE 0 END) as best_pct
		 FROM quiz_sessions qs
		 JOIN categories c ON qs.category_id = c.id
		 WHERE qs.user_id = $1 AND qs.finished_at IS NOT NULL
		 GROUP BY qs.category_id, c.name, qs.difficulty
		 ORDER BY c.name, qs.difficulty`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []BestScore
	for rows.Next() {
		var b BestScore
		if err := rows.Scan(&b.CategoryID, &b.CategoryName, &b.Difficulty, &b.BestScore, &b.Total, &b.Percentage); err != nil {
			return nil, err
		}
		results = append(results, b)
	}
	return results, rows.Err()
}

type DailyStats struct {
	Date    string
	Total   int
	Correct int
}

// Suppress unused import
var _ = json.Marshal
