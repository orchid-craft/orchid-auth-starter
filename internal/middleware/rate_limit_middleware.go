package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	bunker "github.com/yudhiana/bunker/errors"
	httpUtil "github.com/yudhiana/web-api/internal/bootstrap/server/restful_server/http_util"
)

type clientRecord struct {
	count     int
	resetTime time.Time
}

// RateLimiterMiddleware membatasi jumlah permintaan per IP dalam rentang waktu tertentu
type RateLimiterMiddleware struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	clients map[string]*clientRecord
	stopCh  chan struct{}
}

// NewRateLimiterMiddleware membuat middleware rate limiting berbasis sliding/fixed window per IP
func NewRateLimiterMiddleware(limit int, window time.Duration) *RateLimiterMiddleware {
	rl := &RateLimiterMiddleware{
		limit:   limit,
		window:  window,
		clients: make(map[string]*clientRecord),
		stopCh:  make(chan struct{}),
	}

	// Goroutine pembersihan record expired secara periodik
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rl.cleanup()
			case <-rl.stopCh:
				return
			}
		}
	}()

	return rl
}

func (rl *RateLimiterMiddleware) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, rec := range rl.clients {
		if now.After(rec.resetTime) {
			delete(rl.clients, ip)
		}
	}
}

// Close menghentikan background ticker pembersihan
func (rl *RateLimiterMiddleware) Close() {
	select {
	case <-rl.stopCh:
		// already closed
	default:
		close(rl.stopCh)
	}
}

// Limit memeriksa apakah pemanggil telah melebihi batas request
func (rl *RateLimiterMiddleware) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		rl.mu.Lock()
		now := time.Now()
		rec, exists := rl.clients[ip]

		if !exists || now.After(rec.resetTime) {
			rl.clients[ip] = &clientRecord{
				count:     1,
				resetTime: now.Add(rl.window),
			}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if rec.count >= rl.limit {
			retryAfter := int(time.Until(rec.resetTime).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			rl.mu.Unlock()

			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			err := bunker.New(bunker.StatusTooManyRequests).SetMessage("terlalu banyak permintaan, silakan coba lagi beberapa saat lagi")
			httpUtil.ErrorResponse(w, err)
			return
		}

		rec.count++
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// getClientIP mengekstrak IP klien dari header X-Forwarded-For, X-Real-IP, atau RemoteAddr
func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
