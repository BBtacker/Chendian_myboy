package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT 签名密钥（默认占位；生产环境必须通过 SetJWTSecret 从配置注入，
// 且各服务（auth 签名 / gateway 校验）的 Secret 必须保持一致）
var jwtSecret = []byte("change-me-jwt-secret")

// SetJWTSecret 设置JWT密钥（可选，默认使用内置密钥）
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// Claims JWT 载荷
type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint64, username, name, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Name:     name,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "facetest-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
