package middleware

import (
	"net/http"
	"strings"
)

// BodyLimit はリクエストボディのサイズを制限するミドルウェア。
// 上限なしのボディ読み込みでプロセスを OOM に追い込む攻撃を防ぐ。
// ファイルアップロード系のパスだけ大きめの上限を許可する。
func BodyLimit(defaultLimit, uploadLimit int64) func(http.Handler) http.Handler {
	uploadPaths := []string{
		"/api/admin/documents",        // 文書アップロード
		"/api/admin/questions/import", // CSV 一括インポート
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limit := defaultLimit
			for _, p := range uploadPaths {
				if strings.HasPrefix(r.URL.Path, p) {
					limit = uploadLimit
					break
				}
			}
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}
