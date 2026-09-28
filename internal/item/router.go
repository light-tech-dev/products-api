package item

import (
	"github.com/gofiber/fiber/v2"
)

// Register يسجّل مسارات Item (nested تحت products).
func Register(router fiber.Router) {
	svc := NewService(nil)
	h := NewHandler(svc)

	// Nested: /products/:product_id/items
	products := router.Group("/products")
	products.Post("/:product_id/items", h.Create)
	products.Get("/:product_id/items", h.ListByProduct)

	// Standalone: /items/:id
	items := router.Group("/items")
	items.Delete("/:id", h.Delete)
}
