package auth

import (
	"context"
	"errors"
	"time"

	"github.com/abdallah-elngar/gormx"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"products-api/internal/config"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailTaken         = errors.New("email already taken")
)

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Query(ctx context.Context) *gormx.QuerySet[User] {
	return gormx.New[User]().WithContext(ctx)
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	exists, err := s.Query(ctx).Q(gormx.QOr(
		gormx.Eq("username", req.Username),
		gormx.Eq("email", req.Email),
	)).Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
		IsActive: true,
	}

	if err := s.Query(ctx).Create(user); err != nil {
		return nil, err
	}

	return FromUser(user), nil
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := s.Query(ctx).Find("username", req.Username)
	if err != nil {
		if gormx.IsNotFound(err) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  FromUser(user),
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.Query(ctx).Get(id)
	if err != nil {
		return nil, err
	}
	return FromUser(user), nil
}

func (s *Service) generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Duration(s.cfg.JWT.Hours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}
