package api

import "github.com/gofiber/fiber/v2"

// Schema — GET /api/v1/schema
//
// يرجّع schema كامل للـ API (شبيه بـ OpenAPI لكن مبسط).
func Schema(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"openapi": "3.0.0",
		"info": fiber.Map{
			"title":       "Products API",
			"version":     "1.0.0",
			"description": "REST API for managing products and items",
		},
		"servers": []fiber.Map{
			{"url": "http://localhost:8080", "description": "Development"},
		},
		"paths": fiber.Map{
			"/api/v1/": fiber.Map{
				"get": fiber.Map{
					"summary":     "API Info",
					"description": "Returns general API information",
					"tags":        []string{"info"},
				},
			},
			"/api/v1/endpoints": fiber.Map{
				"get": fiber.Map{
					"summary":     "List Endpoints",
					"description": "Returns all available endpoints",
					"tags":        []string{"info"},
				},
			},
			"/api/v1/auth/register": fiber.Map{
				"post": fiber.Map{
					"summary": "Register",
					"tags":    []string{"auth"},
					"requestBody": fiber.Map{
						"required": true,
						"content": fiber.Map{
							"application/json": fiber.Map{
								"schema": fiber.Map{
									"type": "object",
									"properties": fiber.Map{
										"username": fiber.Map{"type": "string"},
										"email":    fiber.Map{"type": "string"},
										"password": fiber.Map{"type": "string"},
									},
									"required": []string{"username", "email", "password"},
								},
							},
						},
					},
					"responses": fiber.Map{
						"201": fiber.Map{"description": "User created"},
						"409": fiber.Map{"description": "Username/email taken"},
					},
				},
			},
			"/api/v1/auth/login": fiber.Map{
				"post": fiber.Map{
					"summary": "Login",
					"tags":    []string{"auth"},
					"responses": fiber.Map{
						"200": fiber.Map{"description": "Login success"},
						"401": fiber.Map{"description": "Invalid credentials"},
					},
				},
			},
			"/api/v1/products": fiber.Map{
				"get": fiber.Map{
					"summary": "List Products",
					"tags":    []string{"products"},
					"parameters": []fiber.Map{
						{"name": "page", "in": "query", "schema": fiber.Map{"type": "integer"}},
						{"name": "per_page", "in": "query", "schema": fiber.Map{"type": "integer"}},
						{"name": "category", "in": "query", "schema": fiber.Map{"type": "string"}},
						{"name": "min_price", "in": "query", "schema": fiber.Map{"type": "number"}},
						{"name": "max_price", "in": "query", "schema": fiber.Map{"type": "number"}},
						{"name": "q", "in": "query", "schema": fiber.Map{"type": "string"}},
					},
					"security": []fiber.Map{{"bearerAuth": []string{}}},
				},
				"post": fiber.Map{
					"summary":  "Create Product",
					"tags":     []string{"products"},
					"security": []fiber.Map{{"bearerAuth": []string{}}},
				},
			},
			"/api/v1/products/{id}": fiber.Map{
				"get":    fiber.Map{"summary": "Get Product", "tags": []string{"products"}, "security": []fiber.Map{{"bearerAuth": []string{}}}},
				"put":    fiber.Map{"summary": "Update Product", "tags": []string{"products"}, "security": []fiber.Map{{"bearerAuth": []string{}}}},
				"delete": fiber.Map{"summary": "Delete Product", "tags": []string{"products"}, "security": []fiber.Map{{"bearerAuth": []string{}}}},
			},
			"/api/v1/products/{product_id}/items": fiber.Map{
				"get":  fiber.Map{"summary": "List Items", "tags": []string{"items"}, "security": []fiber.Map{{"bearerAuth": []string{}}}},
				"post": fiber.Map{"summary": "Create Item", "tags": []string{"items"}, "security": []fiber.Map{{"bearerAuth": []string{}}}},
			},
			"/api/v1/items/{id}": fiber.Map{
				"delete": fiber.Map{"summary": "Delete Item", "tags": []string{"items"}, "security": []fiber.Map{{"bearerAuth": []string{}}}},
			},
		},
		"components": fiber.Map{
			"securitySchemes": fiber.Map{
				"bearerAuth": fiber.Map{
					"type":   "http",
					"scheme": "bearer",
					"bearerFormat": "JWT",
				},
			},
		},
	})
}
