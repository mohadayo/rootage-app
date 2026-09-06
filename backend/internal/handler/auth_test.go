package handler

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestRegister_Success(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("reg-ok")
	defer env.cleanup(t, email)

	w := env.request("POST", "/api/auth/register", map[string]string{
		"name":     "Test User",
		"email":    email,
		"password": "password123",
	}, "")

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusCreated, w.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["token"] == nil || resp["token"] == "" {
		t.Error("response should contain token")
	}
	user := resp["user"].(map[string]any)
	if user["email"] != email {
		t.Errorf("email = %v, want %v", user["email"], email)
	}
	if user["role"] != "user" {
		t.Errorf("role = %v, want user", user["role"])
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("reg-dup")
	defer env.cleanup(t, email)

	// 1回目
	env.request("POST", "/api/auth/register", map[string]string{
		"name": "User1", "email": email, "password": "password123",
	}, "")

	// 2回目
	w := env.request("POST", "/api/auth/register", map[string]string{
		"name": "User2", "email": email, "password": "password456",
	}, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("reg-short")
	defer env.cleanup(t, email)

	w := env.request("POST", "/api/auth/register", map[string]string{
		"name": "User", "email": email, "password": "12345",
	}, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("reg-weak")
	defer env.cleanup(t, email)

	// 長さは足りるが、よくある弱いパスワードは拒否する
	w := env.request("POST", "/api/auth/register", map[string]string{
		"name": "User", "email": email, "password": "password",
	}, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLogin_CaseInsensitiveEmail(t *testing.T) {
	env := setupTestEnv(t)
	base := uniqueEmail("case")
	mixed := "Mixed-" + base // 大文字を含むアドレスで登録
	defer env.cleanup(t, mixed)

	env.registerUser(t, "Case", mixed, "password123")

	// 小文字で入力してもログインできる（正規化されているため）
	w := env.request("POST", "/api/auth/login", map[string]string{
		"email": "mixed-" + base, "password": "password123",
	}, "")
	if w.Code != http.StatusOK {
		t.Errorf("case-insensitive login status = %d, want %d (body=%s)", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestRegister_MissingFields(t *testing.T) {
	env := setupTestEnv(t)

	cases := []struct {
		name string
		body map[string]string
	}{
		{"no email", map[string]string{"name": "User", "password": "password123"}},
		{"no password", map[string]string{"name": "User", "email": "x@test.com"}},
		{"no name", map[string]string{"email": "x@test.com", "password": "password123"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := env.request("POST", "/api/auth/register", tc.body, "")
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestLogin_Success(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("login-ok")
	defer env.cleanup(t, email)

	env.registerUser(t, "Test", email, "password123")

	w := env.request("POST", "/api/auth/login", map[string]string{
		"email": email, "password": "password123",
	}, "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["token"] == nil || resp["token"] == "" {
		t.Error("response should contain token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("login-wrong")
	defer env.cleanup(t, email)

	env.registerUser(t, "Test", email, "password123")

	w := env.request("POST", "/api/auth/login", map[string]string{
		"email": email, "password": "wrongpassword",
	}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestLogin_NonexistentUser(t *testing.T) {
	env := setupTestEnv(t)

	w := env.request("POST", "/api/auth/login", map[string]string{
		"email": "nonexistent@test.com", "password": "password123",
	}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	env := setupTestEnv(t)

	r := env.request("POST", "/api/auth/login", nil, "")
	if r.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", r.Code, http.StatusBadRequest)
	}
}
