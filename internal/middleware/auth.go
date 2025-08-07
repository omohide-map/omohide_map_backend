package middleware

import (
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	appErrors "github.com/omohide_map_backend/pkg/errors"
)

func JWTMiddleware(client *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return appErrors.MissingAuthHeader()
			}

			idToken := strings.TrimPrefix(authHeader, "Bearer ")
			if idToken == authHeader {
				return appErrors.InvalidAuthFormat()
			}

			token, err := client.VerifyIDToken(c.Request().Context(), idToken)
			if err != nil {
				return appErrors.InvalidToken(err.Error())
			}
			c.Set("userID", token.UID)
			c.Set("requestTime", time.Now())
			return next(c)
		}
	}
}
