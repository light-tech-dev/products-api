package product

import (
	"github.com/gofiber/fiber/v2"

	"products-api/internal/common"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Create — POST /api/v1/products
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return common.BadRequest(c, "invalid body")
	}

	product, err := h.svc.Create(c.Context(), &req)
	if err != nil {
		if err == ErrSKUAlreadyExists {
			return common.Conflict(c, err.Error())
		}
		return common.BadRequest(c, err.Error())
	}

	return common.Created(c, product)
}

// GetByID — GET /api/v1/products/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := common.GetUintParam(c, "id")
	if err != nil {
		return common.BadRequest(c, "invalid id")
	}

	product, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return common.NotFound(c, err.Error())
	}

	return common.OK(c, product)
}

// List — GET /api/v1/products
func (h *Handler) List(c *fiber.Ctx) error {
	p := common.GetPagination(c)

	filters := map[string]any{
		"category":  c.Query("category"),
		"search":    c.Query("q"),
		"is_active": true,
	}
	if minPrice, ok := common.GetFloat(c, "min_price"); ok {
		filters["min_price"] = minPrice
	}
	if maxPrice, ok := common.GetFloat(c, "max_price"); ok {
		filters["max_price"] = maxPrice
	}

	result, err := h.svc.List(c.Context(), p.Page, p.PerPage, filters)
	if err != nil {
		return common.Internal(c, err.Error())
	}

	return c.JSON(result)
}

// Update — PUT /api/v1/products/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := common.GetUintParam(c, "id")
	if err != nil {
		return common.BadRequest(c, "invalid id")
	}

	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return common.BadRequest(c, "invalid body")
	}

	product, err := h.svc.Update(c.Context(), id, &req)
	if err != nil {
		return common.BadRequest(c, err.Error())
	}

	return common.OK(c, product)
}

// Delete — DELETE /api/v1/products/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := common.GetUintParam(c, "id")
	if err != nil {
		return common.BadRequest(c, "invalid id")
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return common.BadRequest(c, err.Error())
	}

	return common.NoContent(c)
}
