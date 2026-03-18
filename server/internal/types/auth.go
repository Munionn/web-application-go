package types

import "github.com/golang-jwt/jwt/v5"

type AuthRequest struct {
	Login          string `json:"login"`
	HashedPassword string `json:"password"`
}

type SignInResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	UserID    uint   `json:"user_id"`
	Login     string `json:"login"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
