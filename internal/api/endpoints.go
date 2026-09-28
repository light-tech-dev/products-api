package api

import "github.com/gofiber/fiber/v2"

// ═══════════════════════════════════════════════
// Endpoint يمثل endpoint واحد
// ═══════════════════════════════════════════════

type Endpoint struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Auth        bool     `json:"auth_required"`
	Tags        []string `json:"tags"`
}

// Endpoints — GET /api/v1/endpoints
//
// يرجّع قائمة كل الـ endpoints.
func Endpoints(c *fiber.Ctx) error {
	endpoints := []Endpoint{
		// ═══ API Info ═══
		{"GET", "/api/v1/", "API information", false, []string{"info"}},
		{"GET", "/api/v1/endpoints", "List all endpoints", false, []string{"info"}},
		{"GET", "/api/v1/schema", "API schema", false, []string{"info"}},
		{"GET", "/health", "Health check", false, []string{"info"}},

		// ═══ Auth ═══
		{"POST", "/api/v1/auth/register", "Register new user", false, []string{"auth"}},
		{"POST", "/api/v1/auth/login", "Login and get JWT", false, []string{"auth"}},
		{"GET", "/api/v1/auth/me", "Get current user", true, []string{"auth"}},

		// ═══ Products ═══
		{"POST", "/api/v1/products", "Create product", true, []string{"products"}},
		{"GET", "/api/v1/products", "List products (with pagination & filters)", true, []string{"products"}},
		{"GET", "/api/v1/products/:id", "Get product by ID", true, []string{"products"}},
		{"PUT", "/api/v1/products/:id", "Update product", true, []string{"products"}},
		{"DELETE", "/api/v1/products/:id", "Delete product", true, []string{"products"}},

		// ═══ Items (nested) ═══
		{"POST", "/api/v1/products/:product_id/items", "Create item for product", true, []string{"items"}},
		{"GET", "/api/v1/products/:product_id/items", "List items of product", true, []string{"items"}},
		{"DELETE", "/api/v1/items/:id", "Delete item", true, []string{"items"}},
	}

	// تجميع حسب tag
	grouped := make(map[string][]Endpoint)
	for _, e := range endpoints {
		for _, tag := range e.Tags {
			grouped[tag] = append(grouped[tag], e)
		}
	}

	return c.JSON(fiber.Map{
		"total":     len(endpoints),
		"grouped":   grouped,
		"flat":      endpoints,
	})
}
