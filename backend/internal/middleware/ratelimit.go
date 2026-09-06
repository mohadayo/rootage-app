package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rootage-ses-quiz/backend/internal/pkg"
)

// rateLimiter は IP 単位の固定ウィンドウ・レートリミッタ。
// 総当たり・アカウント濫造・メール/OpenAI 濫用を抑える。
type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	max    int
	window time.Duration
}

// RateLimit は window あたり max 回までに制限するミドルウェアを返す。
func RateLimit(max int, window time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	rl := &rateLimiter{
		hits:   make(map[string][]time.Time),
		max:    max,
		window: window,
	}
	go rl.cleanupLoop()

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !rl.allow(clientIP(r)) {
				pkg.WriteError(w, http.StatusTooManyRequests, "リクエストが多すぎます。しばらく待ってから再度お試しください")
				return
			}
			next(w, r)
		}
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	now := time.Now()
	cutoff := now.Add(-rl.window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	kept := rl.hits[ip][:0]
	for _, t := range rl.hits[ip] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= rl.max {
		rl.hits[ip] = kept
		return false
	}
	rl.hits[ip] = append(kept, now)
	return true
}

// cleanupLoop は古くなった IP エントリを定期的に削除する（map の無限増加=DoSを防ぐ）。
func (rl *rateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		cutoff := time.Now().Add(-rl.window)
		rl.mu.Lock()
		for ip, times := range rl.hits {
			fresh := false
			for _, t := range times {
				if t.After(cutoff) {
					fresh = true
					break
				}
			}
			if !fresh {
				delete(rl.hits, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// clientIP は X-Forwarded-For の先頭ホップ（Render が設定）を優先して取得する。
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
