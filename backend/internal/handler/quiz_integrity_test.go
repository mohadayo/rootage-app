package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/lib/pq"
)

// startQuiz はクイズを開始し、セッションIDと出題された問題IDを返す
func (e *testEnv) startQuiz(t *testing.T, token string) (string, []string) {
	t.Helper()

	catW := e.request("GET", "/api/categories", nil, token)
	var categories []map[string]any
	json.NewDecoder(catW.Body).Decode(&categories)
	if len(categories) == 0 {
		t.Fatal("no categories available")
	}

	startW := e.request("POST", "/api/quiz/start", map[string]string{
		"category_id": categories[0]["id"].(string),
		"difficulty":  "beginner",
	}, token)
	if startW.Code != http.StatusOK {
		t.Fatalf("quiz start: status=%d body=%s", startW.Code, startW.Body.String())
	}

	var startResp map[string]any
	json.NewDecoder(startW.Body).Decode(&startResp)
	questions := startResp["questions"].([]any)
	ids := make([]string, len(questions))
	for i, q := range questions {
		ids[i] = q.(map[string]any)["id"].(string)
	}
	return startResp["session_id"].(string), ids
}

// 同じ問題に繰り返し回答してスコアを水増しできないこと
func TestQuizAnswer_RejectsDuplicateAnswer(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("quiz-dup")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "DupUser", email, "password123")
	sessionID, questionIDs := env.startQuiz(t, token)

	answer := map[string]any{
		"session_id":     sessionID,
		"question_id":    questionIDs[0],
		"selected_index": 0,
	}

	if w := env.request("POST", "/api/quiz/answer", answer, token); w.Code != http.StatusOK {
		t.Fatalf("1回目の回答: status=%d body=%s", w.Code, w.Body.String())
	}

	w := env.request("POST", "/api/quiz/answer", answer, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("2回目の回答: status = %d, want %d (body=%s)", w.Code, http.StatusBadRequest, w.Body.String())
	}

	finishW := env.request("POST", "/api/quiz/finish", map[string]string{"session_id": sessionID}, token)
	var finishResp map[string]any
	json.NewDecoder(finishW.Body).Decode(&finishResp)

	score := int(finishResp["score"].(float64))
	total := int(finishResp["total"].(float64))
	if score > total {
		t.Errorf("score = %d, total = %d: スコアが出題数を超えている", score, total)
	}
	if pct := finishResp["percentage"].(float64); pct > 100 {
		t.Errorf("percentage = %v, want <= 100", pct)
	}
}

// そのセッションで出題されていない問題には回答できず、正解も返さないこと
func TestQuizAnswer_RejectsQuestionOutsideSession(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("quiz-foreign")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "ForeignUser", email, "password123")
	sessionID, questionIDs := env.startQuiz(t, token)

	// 出題された10問に含まれない問題を1件選ぶ
	var foreignID string
	err := env.db.QueryRow(
		`SELECT id FROM questions WHERE NOT (id = ANY($1::uuid[])) LIMIT 1`,
		pq.Array(questionIDs),
	).Scan(&foreignID)
	if err != nil {
		t.Fatalf("出題外の問題を取得できませんでした: %v", err)
	}

	w := env.request("POST", "/api/quiz/answer", map[string]any{
		"session_id":     sessionID,
		"question_id":    foreignID,
		"selected_index": 0,
	}, token)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d (body=%s)", w.Code, http.StatusBadRequest, w.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if _, leaked := resp["correct_index"]; leaked {
		t.Error("出題外の問題に対して correct_index を返している")
	}
	if _, leaked := resp["explanation"]; leaked {
		t.Error("出題外の問題に対して explanation を返している")
	}
}

// 終了済みセッションを再度終了できないこと
func TestQuizFinish_RejectsSecondFinish(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("quiz-refin")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "RefinUser", email, "password123")
	sessionID, questionIDs := env.startQuiz(t, token)

	env.request("POST", "/api/quiz/answer", map[string]any{
		"session_id":     sessionID,
		"question_id":    questionIDs[0],
		"selected_index": 0,
	}, token)

	if w := env.request("POST", "/api/quiz/finish", map[string]string{"session_id": sessionID}, token); w.Code != http.StatusOK {
		t.Fatalf("1回目の終了: status=%d body=%s", w.Code, w.Body.String())
	}

	w := env.request("POST", "/api/quiz/finish", map[string]string{"session_id": sessionID}, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("2回目の終了: status = %d, want %d (body=%s)", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

// 選択肢の範囲外のインデックスを受け付けないこと
func TestQuizAnswer_RejectsOutOfRangeIndex(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("quiz-range")
	defer env.cleanup(t, email)

	token := env.registerUser(t, "RangeUser", email, "password123")
	sessionID, questionIDs := env.startQuiz(t, token)

	for _, idx := range []int{-1, 4, -999} {
		w := env.request("POST", "/api/quiz/answer", map[string]any{
			"session_id":     sessionID,
			"question_id":    questionIDs[0],
			"selected_index": idx,
		}, token)
		if w.Code != http.StatusBadRequest {
			t.Errorf("selected_index=%d: status = %d, want %d (body=%s)", idx, w.Code, http.StatusBadRequest, w.Body.String())
		}
	}
}
