package product

import (
	"github.com/gofiber/fiber/v2"
)

// Register يسجّل مسارات Product.
//
//     path('products/', views.ProductListCreate.as_view()),
//     path('products/<int:pk>/', views.ProductDetail.as_view()),
// ]
func Register(router fiber.Router) {
	svc := NewService(nil)
	h := NewHandler(svc)

	products := router.Group("/products")

	products.Post("/", h.Create)
	products.Get("/", h.List)
	products.Get("/:id", h.GetByID)
	products.Put("/:id", h.Update)
	products.Delete("/:id", h.Delete)
}
