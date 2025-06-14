package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rx3lixir/agg-api/token"
)

type authKey struct{}

func GetAuthMiddleWareFunc(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read the authorization header
			// Verify the token
			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error verifying token: %v", err), http.StatusUnauthorized)
				return
			}

			// Pass the payload/claims down the context
			ctx := context.WithValue(r.Context(), authKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAdminAuthMiddleWareFunc(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read the authorization header
			// Verify the token
			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error verifying token: %v", err), http.StatusUnauthorized)
				return
			}

			// Checking whether the user is admin
			if !claims.IsAdmin {
				http.Error(w, "user is not an admin", http.StatusUnauthorized)
				return
			}

			// Pass the payload/claims down the context
			ctx := context.WithValue(r.Context(), authKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifyClaimsFromAuthHeader(r *http.Request, tokenMaker *token.JWTMaker) (*token.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("Authorization header is missing")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	token := fields[1]

	claims, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
