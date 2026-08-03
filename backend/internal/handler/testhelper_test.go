package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rootage-ses-quiz/backend/internal/middleware"
	"github.com/rootage-ses-quiz/backend/internal/repository"
	"github.com/rootage-ses-quiz/backend/internal/service"
)

const testJWTSecret = "test-jwt-secret-key"

type testEnv struct {
	db          *sql.DB
	mux         *http.ServeMux
	authSvc     *service.AuthService
	authHandler *AuthHandler
	quizHandler *QuizHandler
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()

	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://ses_quiz_user:ses_quiz_pass@localhost:5432/ses_quiz?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	quizRepo := repository.NewQuizRepository(db)

	resetRepo := repository.NewPasswordResetRepository(db)
	authSvc := service.NewAuthService(userRepo, resetRepo, testJWTSecret, "", "", "http://localhost:5173")
	quizSvc := service.NewQuizService(quizRepo, questionRepo, categoryRepo)

	authHandler := NewAuthHandler(authSvc)
	quizHandler := NewQuizHandler(quizSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/categories", middleware.Auth(testJWTSecret, quizHandler.ListCategories))
	mux.HandleFunc("POST /api/quiz/start", middleware.Auth(testJWTSecret, quizHandler.Start))
	mux.HandleFunc("POST /api/quiz/answer", middleware.Auth(testJWTSecret, quizHandler.Answer))
	mux.HandleFunc("POST /api/quiz/finish", middleware.Auth(testJWTSecret, quizHandler.Finish))
	mux.HandleFunc("GET /api/quiz/stats", middleware.Auth(testJWTSecret, quizHandler.Stats))
	mux.HandleFunc("GET /api/users/me/stats", middleware.Auth(testJWTSecret, quizHandler.Dashboard))
	mux.HandleFunc("GET /api/quiz/review", middleware.Auth(testJWTSecret, quizHandler.Review))

	return &testEnv{
		db:          db,
		mux:         mux,
		authSvc:     authSvc,
		authHandler: authHandler,
		quizHandler: quizHandler,
	}
}

func (e *testEnv) cleanup(t *testing.T, email string) {
	t.Helper()
	// テストユーザーと関連データを削除
	var userID string
	err := e.db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return
	}
	e.db.Exec("DELETE FROM quiz_answers WHERE session_id IN (SELECT id FROM quiz_sessions WHERE user_id = $1)", userID)
	e.db.Exec("DELETE FROM quiz_sessions WHERE user_id = $1", userID)
	e.db.Exec("DELETE FROM users WHERE id = $1", userID)
}

func (e *testEnv) request(method, path string, body any, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = &bytes.Buffer{}
	}

	r := httptest.NewRequest(method, path, reqBody)
	r.Header.Set("Content-Type", "application/json")
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	e.mux.ServeHTTP(w, r)
	return w
}

func (e *testEnv) registerUser(t *testing.T, name, email, password string) string {
	t.Helper()
	w := e.request("POST", "/api/auth/register", map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}, "")

	if w.Code != http.StatusCreated {
		t.Fatalf("register failed: status=%d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	token, ok := resp["token"].(string)
	if !ok || token == "" {
		t.Fatal("register response missing token")
	}
	return token
}

func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@test.example.com", prefix, os.Getpid())
}
