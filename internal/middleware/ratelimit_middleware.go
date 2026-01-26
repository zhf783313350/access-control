package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type RateLimitMiddleware struct {
	limiter *limit.TokenLimiter
}

func NewRateLimitMiddleware(limiter *limit.TokenLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: limiter,
	}
}

func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.Allow() {
			httpx.WriteJson(w, http.StatusTooManyRequests, map[string]string{
				"code":    "429",
				"message": "Too Many Requests - Please slow down",
			})
			return
		}
		next(w, r)
	}
}
