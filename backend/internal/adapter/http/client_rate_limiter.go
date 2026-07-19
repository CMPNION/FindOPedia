package http

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type ClientLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	b        int
}

func NewClientLimiter(r rate.Limit, b int) *ClientLimiter {
	return &ClientLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

func (cl *ClientLimiter) GetLimiter(clientID string) *rate.Limiter {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	limiter, exists := cl.limiters[clientID]
	if !exists {
		limiter = rate.NewLimiter(cl.r, cl.b)
		cl.limiters[clientID] = limiter
	}
	return limiter
}

func (cl *ClientLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var key string
			if userID, ok := userIDFromCtx(r); ok {
				key = fmt.Sprintf("uid:%d", userID)
			} else {
				key = "ip:" + r.RemoteAddr
			}
			if !cl.GetLimiter(key).Allow() {
				writeError(w, http.StatusTooManyRequests, "too many requests")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
