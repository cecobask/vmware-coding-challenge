package middleware

import (
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type RequestLoggerMiddleware struct {
	logger *zap.Logger
}

func NewRequestLoggerMiddleware(logger *zap.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{
		logger: logger,
	}
}

func (rlm *RequestLoggerMiddleware) Handle(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestTime := time.Now()
		rwr := chiMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			rlm.logger.Info(
				"request info",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", rwr.Status()),
				zap.Duration("duration", time.Since(requestTime)),
			)
		}()
		next.ServeHTTP(rwr, r)
	}
	return http.HandlerFunc(fn)
}
