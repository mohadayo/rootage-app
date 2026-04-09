package middleware

import (
	"net/http"

	"github.com/rootage-ses-quiz/backend/internal/pkg"
)

func Admin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := GetRole(r.Context())
		if role != "admin" {
			pkg.WriteError(w, http.StatusForbidden, "管理者権限が必要です")
			return
		}
		next.ServeHTTP(w, r)
	}
}
