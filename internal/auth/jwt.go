package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tugkanmeral/the-host-go/internal/config"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, username string) (string, error) {
	cfg := config.LoadConfig()

	secret := []byte(cfg.JWTSecret)

	expiration := 15 * time.Minute
	if cfg.JWTExpiration != "" {
		if d, err := time.ParseDuration(cfg.JWTExpiration); err == nil && d > 0 {
			expiration = d
		}
	}

	now := time.Now()

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func VerifyToken(tokenStr string) (*Claims, error) {
	cfg := config.LoadConfig()
	secret := []byte(cfg.JWTSecret)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GetUserId(tokenStr string) string {
	cfg := config.LoadConfig()
	secret := []byte(cfg.JWTSecret)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil || token == nil {
		return ""
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return ""
	}

	return claims.UserID
}
