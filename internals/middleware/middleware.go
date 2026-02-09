package middleware

import (
	"strings"

	"tenant-Dynamin-DB/internals/token"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "authorization header missing"})
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid authorization format"})
		}

		tokenStr := strings.TrimPrefix(auth, prefix)

		claims := &token.Claims{}

		tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})

		if err != nil || !tok.Valid {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid or expired token"})
		}

		// store into context
		c.Locals("user_id", claims.UserID)
		c.Locals("tenant", claims.Tenant)

		return c.Next()
	}
}
