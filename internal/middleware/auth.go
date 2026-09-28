package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"products-api/internal/common"
	"products-api/internal/config"
)

func JWTAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return common.Unauthorized(c, "missing authorization")
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == auth {
			return common.Unauthorized(c, "invalid authorization format")
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			return common.Unauthorized(c, "invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return common.Unauthorized(c, "invalid claims")
		}

		c.Locals("user_id", uint(claims["user_id"].(float64)))
		c.Locals("username", claims["username"].(string))

		return c.Next()
	}
}
