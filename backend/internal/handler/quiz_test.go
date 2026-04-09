package handler

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListCategories_Authenticated(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("cat-auth")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "Test", email, "password123")

	w := env.request("GET", "/api/categories", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var categories []map[string]any
	json.NewDecoder(w.Body).Decode(&categories)
	if len(categories) == 0 {
		t.Error("should return at least one category")
	}
}

func TestListCategories_Unauthenticated(t *testing.T) {
	env := setupTestEnv(t)

	w := env.request("GET", "/api/categories", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestQuizFlow_StartAnswerFinish(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("quiz-flow")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "QuizUser", email, "password123")

	// カテゴリ一覧取得
	catW := env.request("GET", "/api/categories", nil, token)
	var categories []map[string]any
	json.NewDecoder(catW.Body).Decode(&categories)
	if len(categories) == 0 {
		t.Fatal("no categories available")
	}
	categoryID := categories[0]["id"].(string)

	// クイズ開始
	startW := env.request("POST", "/api/quiz/start", map[string]string{
		"category_id": categoryID,
		"difficulty":  "beginner",
	}, token)
	if startW.Code != http.StatusOK {
		t.Fatalf("quiz start: status=%d body=%s", startW.Code, startW.Body.String())
	}

	var startResp map[string]any
	json.NewDecoder(startW.Body).Decode(&startResp)
	sessionID := startResp["session_id"].(string)
	questions := startResp["questions"].([]any)

	if sessionID == "" {
		t.Fatal("session_id should not be empty")
	}
	if len(questions) == 0 {
		t.Fatal("should have at least 1 question")
	}

	// 全問回答
	correctCount := 0
	for _, q := range questions {
		qMap := q.(map[string]any)
		ansW := env.request("POST", "/api/quiz/answer", map[string]any{
			"session_id":     sessionID,
			"question_id":    qMap["id"],
			"selected_index": 0,
		}, token)
		if ansW.Code != http.StatusOK {
			t.Fatalf("answer: status=%d body=%s", ansW.Code, ansW.Body.String())
		}

		var ansResp map[string]any
		json.NewDecoder(ansW.Body).Decode(&ansResp)
		if ansResp["is_correct"].(bool) {
			correctCount++
		}
	}

	// 終了
	finishW := env.request("POST", "/api/quiz/finish", map[string]string{
		"session_id": sessionID,
	}, token)
	if finishW.Code != http.StatusOK {
		t.Fatalf("finish: status=%d body=%s", finishW.Code, finishW.Body.String())
	}

	var finishResp map[string]any
	json.NewDecoder(finishW.Body).Decode(&finishResp)

	score := int(finishResp["score"].(float64))
	total := int(finishResp["total"].(float64))

	if score != correctCount {
		t.Errorf("score = %d, want %d", score, correctCount)
	}
	if total != len(questions) {
		t.Errorf("total = %d, want %d", total, len(questions))
	}

	// ダッシュボード確認
	dashW := env.request("GET", "/api/users/me/stats", nil, token)
	if dashW.Code != http.StatusOK {
		t.Fatalf("dashboard: status=%d body=%s", dashW.Code, dashW.Body.String())
	}

	var dashResp map[string]any
	json.NewDecoder(dashW.Body).Decode(&dashResp)
	if int(dashResp["total_sessions"].(float64)) != 1 {
		t.Errorf("total_sessions = %v, want 1", dashResp["total_sessions"])
	}
	if dashResp["achievements"] == nil {
		t.Error("achievements should not be nil")
	}

	// 復習確認
	reviewW := env.request("GET", "/api/quiz/review", nil, token)
	if reviewW.Code != http.StatusOK {
		t.Fatalf("review: status=%d body=%s", reviewW.Code, reviewW.Body.String())
	}
}

func TestQuizStart_InvalidCategory(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("quiz-badcat")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "Test", email, "password123")

	w := env.request("POST", "/api/quiz/start", map[string]string{
		"category_id": "nonexistent-id",
		"difficulty":  "beginner",
	}, token)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestQuizAnswer_WrongSession(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("quiz-badsess")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "Test", email, "password123")

	w := env.request("POST", "/api/quiz/answer", map[string]any{
		"session_id":     "nonexistent-session",
		"question_id":    "nonexistent-question",
		"selected_index": 0,
	}, token)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestQuizFinish_WrongSession(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("quiz-finbad")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "Test", email, "password123")

	w := env.request("POST", "/api/quiz/finish", map[string]string{
		"session_id": "nonexistent-session",
	}, token)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestQuizStats_NoData(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("quiz-nostats")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "Test", email, "password123")

	w := env.request("GET", "/api/quiz/stats", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if int(resp["total_sessions"].(float64)) != 0 {
		t.Errorf("total_sessions = %v, want 0", resp["total_sessions"])
	}
}

func TestDashboard_NoData(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("dash-empty")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "Test", email, "password123")

	w := env.request("GET", "/api/users/me/stats", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if int(resp["total_sessions"].(float64)) != 0 {
		t.Errorf("total_sessions = %v, want 0", resp["total_sessions"])
	}
}
