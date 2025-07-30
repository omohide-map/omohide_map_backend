// Package middleware provides middleware for the application.
package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type SupabaseClaims struct {
	jwt.RegisteredClaims
	UserMetadata map[string]interface{} `json:"user_metadata"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone"`
	Role         string                 `json:"role"`
}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization format",
				})
			}

			jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
			if jwtSecret == "" {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "JWT secret not configured",
				})
			}

			token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			if claims, ok := token.Claims.(*SupabaseClaims); ok && token.Valid {
				c.Set("user_id", claims.Subject)
				c.Set("user_email", claims.Email)
				c.Set("user_claims", claims)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token claims",
			})
		}
	}
}
