package bloghttp

import (
	"context"
	authpb "micro-blog/proto/auth"
	"net/http"
	"strings"
)

type Authenticator interface {
	ValidateToken(ctx context.Context, token string) (*authpb.ValidateTokenResponse, error)
}

func AuthMiddleware(authenticator Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(
					w,
					"missing authorization header",
					http.StatusUnauthorized)
				return
			}
			token := strings.Trim(authHeader, "Bearer ")

			resp, err := authenticator.ValidateToken(r.Context(), token)
			if err != nil {
				http.Error(
					w,
					"validate token error",
					http.StatusInternalServerError,
				)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, resp.UserId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
