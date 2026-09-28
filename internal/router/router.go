package router

import (
	"github.com/gofiber/fiber/v2"

	"products-api/internal/api"
	"products-api/internal/auth"
	"products-api/internal/config"
	"products-api/internal/item"
	"products-api/internal/middleware"
	"products-api/internal/product"
)

// Setup يربط كل المسارات.
func Setup(app *fiber.App, cfg *config.Config) {
	// ═══ Health (public) ═══
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	// ═══ API v1 ═══
	apiV1 := app.Group("/api/v1")

	// ═══ API Info (public) ═══
	api.Register(apiV1)

	// ═══ Auth (public + protected) ═══
	auth.Register(apiV1, cfg)

	// ═══ Protected Routes ═══
	authMW := middleware.JWTAuth(cfg)
	protected := apiV1.Group("", authMW)

	product.Register(protected)
	item.Register(protected)
}
