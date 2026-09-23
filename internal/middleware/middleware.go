package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yukay/CRM/internal/models"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func parseAndVerifyJWT(tokenStr string, secretKey []byte) (int, error) {

	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("Токен невалиден")
	}

	return claims.UserId, nil

}

func AuthMiddleware(secretkey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "отсутствует заголовок авторизации", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "невалидный заголовок авторизации", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			userID, err := parseAndVerifyJWT(token, secretkey)
			if err != nil {
				http.Error(w, "невалидный или просроченный токен", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
