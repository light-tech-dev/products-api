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
	ctx context.Context
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
		ctx: context.Background(),
	}
}

// Query يرجّع QuerySet لـ User.
func (s *Service) Query() *gormx.QuerySet[User] {
	return gormx.New[User]().WithContext(s.ctx)
}

func (s *Service) Register(req *RegisterRequest) (*UserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// فحص التكرار
	exists, err := s.Query().Q(gormx.QOr(
		gormx.Eq("username", req.Username),
		gormx.Eq("email", req.Email),
	)).Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	// تشفير
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

	if err := s.Query().Create(user); err != nil {
		return nil, err
	}

	return FromUser(user), nil
}

func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// جلب المستخدم
	user, err := s.Query().Find("username", req.Username)
	if err != nil {
		if gormx.IsNotFound(err) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// فحص كلمة المرور
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// توليد JWT
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  FromUser(user),
	}, nil
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
