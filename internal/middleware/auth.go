package middleware

import (
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(client *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			idToken := strings.TrimPrefix(authHeader, "Bearer ")
			if idToken == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization format",
				})
			}

			token, err := client.VerifyIDToken(c.Request().Context(), idToken)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}
			c.Set("userID", token.UID)
			c.Set("requestTime", time.Now())
			return next(c)
		}
	}
}
