package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"webapplication/auth"
	"webapplication/internal/models"
	"webapplication/internal/types"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = os.Getenv("JWT_SECRET_KEY")

func generateToken(userID uint, login string) (string, int64, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &types.Claims{
		UserID: userID,
		Login:  login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.Unix(), nil
}

// generateRefreshToken creates a random refresh token
func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *Server) saveRefreshTokenAsync(userID uint, refreshToken string) {
	go func() {
		token := models.Token{
			UserID:       userID,
			RefreshToken: refreshToken,
		}

		if err := s.db.GetDB().Create(&token).Error; err != nil {
			log.Printf("Failed to save refresh token for user %d: %v", userID, err)
		} else {
			log.Printf("Refresh token saved successfully for user %d", userID)
		}
	}()
}

// SignInHandler handles user sign-in requests
func (s *Server) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req types.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.HashedPassword == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}
	var user models.User
	result := s.db.GetDB().Where("login = ?", req.Login).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := auth.ComparePassword(user.HashPassword, req.HashedPassword); !err {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenString, expiresAt, err := generateToken(user.ID, user.Login)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := generateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}
	s.saveRefreshTokenAsync(user.ID, refreshToken)
	response := types.SignInResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
		Login:     user.Login,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req types.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.HashedPassword == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}
	hashedPassword, err := auth.HashPassword(req.HashedPassword)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}
	var existing models.User
	if err := s.db.GetDB().Where("login = ?", req.Login).First(&existing).Error; err == nil {
		http.Error(w, "Login already exists", http.StatusConflict)
		return
	}

	user := models.User{
		Login:        req.Login,
		HashPassword: hashedPassword,
	}
	if err := s.db.GetDB().Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	defaultAccount := models.Account{
		Name:         "Default Account",
		BaseCurrence: "USD",
		Balance:      0,
		Users:        []*models.User{&user},
	}
	if err := s.db.GetDB().Create(&defaultAccount).Error; err != nil {
		log.Printf("Failed to create default account for user %d: %v", user.ID, err)
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}
	s.saveRefreshTokenAsync(user.ID, refreshToken)

	tokenString, expiresAt, err := generateToken(user.ID, user.Login)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.SignInResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
		Login:     user.Login,
	})
}

// RefreshTokenHandler handles refresh token requests and issues new access tokens
func (s *Server) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req types.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// Find refresh token in database
	var token models.Token
	result := s.db.GetDB().Where("refresh_token = ?", req.RefreshToken).First(&token)
	if result.Error != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// Get user associated with the token
	var user models.User
	result = s.db.GetDB().First(&user, token.UserID)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Generate new access token
	tokenString, expiresAt, err := generateToken(user.ID, user.Login)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}
	s.db.GetDB().Model(&token).Update("refresh_token", newRefreshToken)

	response := types.SignInResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
		Login:     user.Login,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
