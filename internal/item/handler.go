package item

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

// Create — POST /api/v1/products/:product_id/items
func (h *Handler) Create(c *fiber.Ctx) error {
	productID, err := common.GetUintParam(c, "product_id")
	if err != nil {
		return common.BadRequest(c, "invalid product id")
	}

	var req CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return common.BadRequest(c, "invalid body")
	}

	item, err := h.svc.Create(productID, &req)
	if err != nil {
		return common.BadRequest(c, err.Error())
	}

	return common.Created(c, item)
}

// ListByProduct — GET /api/v1/products/:product_id/items
func (h *Handler) ListByProduct(c *fiber.Ctx) error {
	productID, err := common.GetUintParam(c, "product_id")
	if err != nil {
		return common.BadRequest(c, "invalid product id")
	}

	items, err := h.svc.ListByProduct(productID)
	if err != nil {
		return common.Internal(c, err.Error())
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"count": len(items),
	})
}

// Delete — DELETE /api/v1/items/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := common.GetUintParam(c, "id")
	if err != nil {
		return common.BadRequest(c, "invalid id")
	}

	if err := h.svc.Delete(id); err != nil {
		return common.Internal(c, err.Error())
	}

	return common.NoContent(c)
}
