package auth

import (
	"github.com/gofiber/fiber/v2"

	"products-api/internal/config"
	"products-api/internal/middleware"
)

// Register registers auth routes.
func Register(router fiber.Router, cfg *config.Config) {
	svc := NewService(cfg)
	h := NewHandler(svc)

	auth := router.Group("/auth")

	// Public
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)

	// Protected
	auth.Get("/me", middleware.JWTAuth(cfg), h.Me)
}
