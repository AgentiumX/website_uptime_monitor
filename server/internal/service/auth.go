package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"uptime-monitor/server/internal/config"
	"uptime-monitor/server/internal/repository"
)

// Claims holds the JWT claims used for authentication.
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthService handles authentication logic.
type AuthService struct {
	cfg *config.ServerConfig
}

// NewAuthService creates a new AuthService with the given server configuration.
func NewAuthService(cfg *config.ServerConfig) *AuthService {
	return &AuthService{cfg: cfg}
}

// Login verifies credentials and returns a JWT token on success.
func (s *AuthService) Login(username, password string) (string, error) {
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", errors.New("生成令牌失败")
	}

	return tokenStr, nil
}

// ValidateToken parses and validates a JWT token string, returning its claims.
func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("认证令牌无效")
	}
	return claims, nil
}
