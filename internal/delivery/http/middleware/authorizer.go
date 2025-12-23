package middleware

import (
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID        string          `json:"sub"`
	Username      string          `json:"username"`
	Email         string          `json:"email"`
	Authorization []Authorization `json:"authorization"`
}

type Authorization struct {
	App         string   `json:"app"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
