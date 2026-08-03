package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// 環境変数が設定されていなければ管理者は作られないこと
func TestEnsureAdminUser_SkipsWhenUnset(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("admin-unset")
	defer env.cleanup(t, email)

	if err := env.authSvc.EnsureAdminUser(context.Background(), "", ""); err != nil {
		t.Fatalf("EnsureAdminUser: %v", err)
	}

	var count int
	if err := env.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Errorf("users = %d, want 0", count)
	}
}

// 設定されていれば管理者が作成され、そのアカウントでログインできること
func TestEnsureAdminUser_CreatesAdmin(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("admin-create")
	defer env.cleanup(t, email)

	if err := env.authSvc.EnsureAdminUser(context.Background(), email, "admin-password"); err != nil {
		t.Fatalf("EnsureAdminUser: %v", err)
	}

	w := env.request("POST", "/api/auth/login", map[string]string{
		"email":    email,
		"password": "admin-password",
	}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("login: status=%d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	user := resp["user"].(map[string]any)
	if user["role"] != "admin" {
		t.Errorf("role = %v, want admin", user["role"])
	}
}

// 既にユーザーが存在する場合はパスワードを上書きしないこと
func TestEnsureAdminUser_DoesNotOverwriteExisting(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("admin-exists")
	defer env.cleanup(t, email)

	env.registerUser(t, "既存ユーザー", email, "original-password")

	if err := env.authSvc.EnsureAdminUser(context.Background(), email, "overwrite-password"); err != nil {
		t.Fatalf("EnsureAdminUser: %v", err)
	}

	w := env.request("POST", "/api/auth/login", map[string]string{
		"email":    email,
		"password": "original-password",
	}, "")
	if w.Code != http.StatusOK {
		t.Errorf("元のパスワードでログインできない: status=%d body=%s", w.Code, w.Body.String())
	}

	w = env.request("POST", "/api/auth/login", map[string]string{
		"email":    email,
		"password": "overwrite-password",
	}, "")
	if w.Code == http.StatusOK {
		t.Error("既存アカウントのパスワードが上書きされている")
	}
}

// 短すぎるパスワードは拒否されること
func TestEnsureAdminUser_RejectsShortPassword(t *testing.T) {
	env := setupTestEnv(t)
	email := uniqueEmail("admin-short")
	defer env.cleanup(t, email)

	if err := env.authSvc.EnsureAdminUser(context.Background(), email, "short"); err == nil {
		t.Error("8文字未満のパスワードを受け付けている")
	}
}
