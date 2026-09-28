package product

import (
	"github.com/gofiber/fiber/v2"
)

// Register registers all product routes.
func Register(router fiber.Router) {
	svc := NewService()
	h := NewHandler(svc)

	products := router.Group("/products")

	products.Post("/", h.Create)
	products.Get("/", h.List)
	products.Get("/:id", h.GetByID)
	products.Put("/:id", h.Update)
	products.Delete("/:id", h.Delete)
}
