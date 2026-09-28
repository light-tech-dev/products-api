package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"products-api/internal/common"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Register — POST /api/v1/auth/register
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return common.BadRequest(c, "invalid body")
	}

	user, err := h.svc.Register(&req)
	if err != nil {
		if strings.Contains(err.Error(), "taken") {
			return common.Conflict(c, err.Error())
		}
		return common.BadRequest(c, err.Error())
	}

	return common.Created(c, user)
}

// Login — POST /api/v1/auth/login
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return common.BadRequest(c, "invalid body")
	}

	result, err := h.svc.Login(&req)
	if err != nil {
		return common.Unauthorized(c, err.Error())
	}

	return common.OK(c, result)
}

// Me — GET /api/v1/auth/me
func (h *Handler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	user, err := h.svc.Query().Get(userID)
	if err != nil {
		return common.NotFound(c, "user not found")
	}

	return common.OK(c, FromUser(user))
}
