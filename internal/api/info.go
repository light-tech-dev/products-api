package api

import (
	"github.com/gofiber/fiber/v2"
)

// Info — GET /api/v1/
//
// يرجّع معلومات عن الـ API.
func Info(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":        "Products API",
		"version":     "1.0.0",
		"description": "REST API for managing products and items",
		"status":      "active",
		"documentation": fiber.Map{
			"endpoints": "/api/v1/endpoints",
			"schema":    "/api/v1/schema",
			"health":    "/health",
		},
		"stats": fiber.Map{
			"total_endpoints": 18,
		},
	})
}
