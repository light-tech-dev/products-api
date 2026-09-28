package auth

import (
	"github.com/gofiber/fiber/v2"

	"products-api/internal/config"
	"products-api/internal/middleware"
)

// Register يسجّل مسارات Auth.
//
//     path('auth/register/', views.Register.as_view()),
//     path('auth/login/', views.Login.as_view()),
//     path('auth/me/', views.Me.as_view()),
// ]
func Register(router fiber.Router, cfg *config.Config) {
	svc := NewService(cfg)
	h := NewHandler(svc)

	auth := router.Group("/auth")

	// Public
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)

	// Protected (يحتاج auth middleware)
	auth.Get("/me", middleware.JWTAuth(cfg), h.Me)
}
