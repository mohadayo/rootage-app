package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rootage-ses-quiz/backend/internal/pkg"
)

const testSecret = "test-secret"

func TestAuth_ValidToken(t *testing.T) {
	token, _ := pkg.GenerateToken("user-1", "test@example.com", "user", 0, testSecret)

	called := false
	handler := Auth(testSecret, nil, func(w http.ResponseWriter, r *http.Request) {
		called = true
		if uid := GetUserID(r.Context()); uid != "user-1" {
			t.Errorf("UserID = %q, want %q", uid, "user-1")
		}
		if role := GetRole(r.Context()); role != "user" {
			t.Errorf("Role = %q, want %q", role, "user")
		}
		w.WriteHeader(http.StatusOK)
	})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler(w, r)

	if !called {
		t.Error("handler should have been called")
	}
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAuth_TokenVersionMatch(t *testing.T) {
	token, _ := pkg.GenerateToken("user-1", "test@example.com", "user", 3, testSecret)

	called := false
	tv := func(ctx context.Context, userID string) (int, error) { return 3, nil }
	handler := Auth(testSecret, tv, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, r)

	if !called || w.Code != http.StatusOK {
		t.Errorf("matching token_version should pass: called=%v code=%d", called, w.Code)
	}
}

func TestAuth_TokenVersionMismatch(t *testing.T) {
	// パスワード変更後（DB側の token_version が進んだ後）の古いトークンは失効する。
	token, _ := pkg.GenerateToken("user-1", "test@example.com", "user", 3, testSecret)

	tv := func(ctx context.Context, userID string) (int, error) { return 4, nil }
	handler := Auth(testSecret, tv, func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for stale token_version")
	})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_NoHeader(t *testing.T) {
	handler := Auth(testSecret, nil, func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_InvalidFormat(t *testing.T) {
	handler := Auth(testSecret, nil, func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	handler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	handler := Auth(testSecret, nil, func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	handler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_WrongSecret(t *testing.T) {
	token, _ := pkg.GenerateToken("user-1", "test@example.com", "user", 0, "other-secret")

	handler := Auth(testSecret, nil, func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAdmin_AdminRole(t *testing.T) {
	called := false
	handler := Admin(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.WithValue(context.Background(), RoleKey, "admin")
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler(w, r)

	if !called {
		t.Error("handler should have been called for admin")
	}
}

func TestAdmin_UserRole(t *testing.T) {
	handler := Admin(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for user role")
	})

	ctx := context.WithValue(context.Background(), RoleKey, "user")
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestAdmin_NoRole(t *testing.T) {
	handler := Admin(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called without role")
	})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestGetUserID_Empty(t *testing.T) {
	ctx := context.Background()
	if uid := GetUserID(ctx); uid != "" {
		t.Errorf("GetUserID should return empty for context without value, got %q", uid)
	}
}

func TestGetRole_Empty(t *testing.T) {
	ctx := context.Background()
	if role := GetRole(ctx); role != "" {
		t.Errorf("GetRole should return empty for context without value, got %q", role)
	}
}
