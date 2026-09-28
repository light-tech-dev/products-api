package auth

import (
	"errors"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════
// RegisterRequest
// ═══════════════════════════════════════════════

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return errors.New("username is required")
	}
	if len(r.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(r.Email, "@") {
		return errors.New("invalid email")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

// ═══════════════════════════════════════════════
// LoginRequest
// ═══════════════════════════════════════════════

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Username == "" || r.Password == "" {
		return errors.New("username and password are required")
	}
	return nil
}

// ═══════════════════════════════════════════════
// UserResponse
// ═══════════════════════════════════════════════

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

func FromUser(u *User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

// ═══════════════════════════════════════════════
// LoginResponse
// ═══════════════════════════════════════════════

type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}
