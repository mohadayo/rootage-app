package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rootage-ses-quiz/backend/internal/pkg"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
	RoleKey   contextKey = "role"
)

// TokenVersionFunc はユーザーの現在のトークン世代を返す。
// パスワード変更・ロール変更で古いトークンを失効させるために使う。
// middleware をリポジトリ実装に結合させないよう関数として注入する。
type TokenVersionFunc func(ctx context.Context, userID string) (int, error)

func Auth(secret string, tokenVersion TokenVersionFunc, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			pkg.WriteError(w, http.StatusUnauthorized, "認証が必要です")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			pkg.WriteError(w, http.StatusUnauthorized, "無効な認証ヘッダーです")
			return
		}

		claims, err := pkg.ValidateToken(parts[1], secret)
		if err != nil {
			pkg.WriteError(w, http.StatusUnauthorized, "無効なトークンです")
			return
		}

		// トークン世代を DB の現在値と突合し、パスワード変更後の古いトークンを失効させる。
		if tokenVersion != nil {
			current, err := tokenVersion(r.Context(), claims.UserID)
			if err != nil || current != claims.TokenVersion {
				pkg.WriteError(w, http.StatusUnauthorized, "無効なトークンです")
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDKey).(string); ok {
		return v
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if v, ok := ctx.Value(RoleKey).(string); ok {
		return v
	}
	return ""
}
