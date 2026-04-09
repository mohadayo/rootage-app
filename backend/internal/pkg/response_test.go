package pkg

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"status": "ok"}

	WriteJSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["status"] != "ok" {
		t.Errorf("body status = %q, want %q", result["status"], "ok")
	}
}

func TestWriteJSON_DifferentStatus(t *testing.T) {
	w := httptest.NewRecorder()
	WriteJSON(w, http.StatusCreated, map[string]int{"id": 1})

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, http.StatusBadRequest, "invalid input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "invalid input" {
		t.Errorf("error = %q, want %q", result["error"], "invalid input")
	}
}

func TestReadJSON(t *testing.T) {
	body := `{"name":"test","value":42}`
	r := httptest.NewRequest("POST", "/", strings.NewReader(body))

	var dst struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	err := ReadJSON(r, &dst)
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}
	if dst.Name != "test" {
		t.Errorf("Name = %q, want %q", dst.Name, "test")
	}
	if dst.Value != 42 {
		t.Errorf("Value = %d, want %d", dst.Value, 42)
	}
}

func TestReadJSON_InvalidJSON(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("not json"))
	var dst map[string]any
	err := ReadJSON(r, &dst)
	if err == nil {
		t.Error("ReadJSON should fail with invalid JSON")
	}
}

func TestReadJSON_EmptyBody(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(""))
	var dst map[string]any
	err := ReadJSON(r, &dst)
	if err == nil {
		t.Error("ReadJSON should fail with empty body")
	}
}
