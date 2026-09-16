package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func Secret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me"
	}
	return []byte(secret)
}

func Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return Secret(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipAuth(c) {
				return next(c)
			}

			header := c.Request().Header.Get(echo.HeaderAuthorization)
			if !strings.HasPrefix(header, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status": http.StatusUnauthorized,
					"error":  "missing bearer token",
				})
			}

			claims, err := Parse(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status": http.StatusUnauthorized,
					"error":  "invalid or expired token",
				})
			}

			c.Set("email", claims.Email)
			c.Set("user_id", claims.UserID)
			return next(c)
		}
	}
}

func UserID(c echo.Context) string {
	if value, ok := c.Get("user_id").(string); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}

	header := c.Request().Header.Get(echo.HeaderAuthorization)
	if strings.HasPrefix(header, "Bearer ") {
		if claims, err := Parse(strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))); err == nil {
			if claims.UserID != "" {
				return claims.UserID
			}
			if claims.Subject != "" {
				return claims.Subject
			}
		}
	}

	return strings.TrimSpace(c.QueryParam("user_id"))
}

func skipAuth(c echo.Context) bool {
	path := c.Request().URL.Path
	if path == "/health" || strings.HasPrefix(path, "/uploads") {
		return true
	}

	method := c.Request().Method
	if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
		return true
	}

	switch path {
	case "/api/login", "/api/register", "/api/forgot", "/api/reset":
		return true
	}
	return false
}
