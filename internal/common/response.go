package common

import (
	"github.com/gofiber/fiber/v2"
)

// ═══════════════════════════════════════════════
// Response — موحّد لكل الـ API
// ═══════════════════════════════════════════════

// Data يرجّع data فقط.
func Data(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"data": data,
	})
}

// OK يرجّع 200 مع data.
func OK(c *fiber.Ctx, data any) error {
	return Data(c, fiber.StatusOK, data)
}

// Created يرجّع 201.
func Created(c *fiber.Ctx, data any) error {
	return Data(c, fiber.StatusCreated, data)
}

// NoContent يرجّع 204.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error يرجّع خطأ.
func Error(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": msg,
	})
}

// BadRequest يرجّع 400.
func BadRequest(c *fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusBadRequest, msg)
}

// NotFound يرجّع 404.
func NotFound(c *fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusNotFound, msg)
}

// Unauthorized يرجّع 401.
func Unauthorized(c *fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusUnauthorized, msg)
}

// Forbidden يرجّع 403.
func Forbidden(c *fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusForbidden, msg)
}

// Conflict يرجّع 409.
func Conflict(c *fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusConflict, msg)
}

// Internal يرجّع 500.
func Internal(c *fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusInternalServerError, msg)
}

// Paginated يرجّع نتيجة pagination.
func Paginated(c *fiber.Ctx, items any, total int64, page, perPage int) error {
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return c.JSON(fiber.Map{
		"items":       items,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
		"has_next":    page < totalPages,
		"has_prev":    page > 1,
	})
}
