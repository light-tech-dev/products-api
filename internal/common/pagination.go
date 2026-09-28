package common

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ═══════════════════════════════════════════════
// Pagination — استخراج من Query
// ═══════════════════════════════════════════════

type PaginationParams struct {
	Page    int
	PerPage int
}

// GetPagination يستخرج page و per_page من Query.
func GetPagination(c *fiber.Ctx) PaginationParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
	}
}

// GetFloat يستخرج float من Query.
func GetFloat(c *fiber.Ctx, key string) (float64, bool) {
	v := c.Query(key)
	if v == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(v, 64)
	return f, err == nil
}

// GetInt يستخرج int من Query.
func GetInt(c *fiber.Ctx, key string) (int, bool) {
	v := c.Query(key)
	if v == "" {
		return 0, false
	}
	i, err := strconv.Atoi(v)
	return i, err == nil
}

// GetUintParam يستخرج uint من Params.
func GetUintParam(c *fiber.Ctx, key string) (uint, error) {
	v := c.Params(key)
	if v == "" {
		return 0, ErrInvalidID
	}
	id, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, ErrInvalidID
	}
	return uint(id), nil
}
