package bloghttp

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Authenticator interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
}

func AuthMiddleware(authenticator Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(
					w,
					"missing authorization header",
					http.StatusUnauthorized)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			allowed, err := authenticator.ValidateToken(
				ctx,
				token,
			)
			if err != nil {
				log.Printf("failed to check authentication: %v", err)
				http.Error(w,
					"internal server error",
					http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w,
					"unauthorized",
					http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
