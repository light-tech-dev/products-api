package common

import "errors"

// ═══════════════════════════════════════════════
// Common Errors — مشتركة بين كل الـ apps
// ═══════════════════════════════════════════════

var (
	// General
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrBadRequest   = errors.New("bad request")

	// Validation
	ErrValidation = errors.New("validation failed")
	ErrInvalidID  = errors.New("invalid id")
)
