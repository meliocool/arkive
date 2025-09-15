package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/meliocool/arkive/internal/helper"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next httprouter.Handle, jwtSecret string) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		authHeader := strings.TrimSpace(request.Header.Get("Authorization"))
		if authHeader == "" {
			helper.WriteErr(writer, helper.ErrUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			helper.WriteErr(writer, helper.ErrUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		var claims Claims

		token, parseErr := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if parseErr != nil {
			log.Printf("Failed to parse token: %v", parseErr)
			helper.WriteErr(writer, helper.ErrUnauthorized)
			return
		}

		if !token.Valid {
			helper.WriteErr(writer, helper.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), ContextKeyUserID, claims.Subject)

		next(writer, request.WithContext(ctx), params)
	}
}
