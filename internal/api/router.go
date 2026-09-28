package api

import "github.com/gofiber/fiber/v2"

// Register يسجّل مسارات API Info.
func Register(router fiber.Router) {
	api := router.Group("/")

	api.Get("/", Info)
	api.Get("/endpoints", Endpoints)
	api.Get("/schema", Schema)
}
