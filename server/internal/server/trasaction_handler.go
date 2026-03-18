package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"webapplication/internal/models"
	"webapplication/internal/types"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func (s *Server) getAllTransactions(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserContextKey).(*types.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	db := s.db.GetDB()
	var transactions []models.Transaction
	result := db.Where("user_id = ?", claims.UserID).Find(&transactions)
	if result.Error != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (s *Server) getAllTransactionForAccount(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserContextKey).(*types.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())
	accountID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil || accountID == 0 {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	db := s.db.GetDB()

	var account models.Account
	err = db.Joins("JOIN user_accounts ON user_accounts.account_id = accounts.id").
		Where("accounts.id = ? AND user_accounts.user_id = ?", accountID, claims.UserID).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch account", http.StatusInternalServerError)
		return
	}

	var transactions []models.Transaction
	result := db.Where("account_id = ?", accountID).Find(&transactions)
	if result.Error != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (s *Server) createTransaction(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserContextKey).(*types.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	params := httprouter.ParamsFromContext(r.Context())
	accountID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil || accountID == 0 {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}
	db := s.db.GetDB()
	var account models.Account
	err = db.Joins("JOIN user_accounts ON user_accounts.account_id = accounts.id").
		Where("accounts.id = ? AND user_accounts.user_id = ?", accountID, claims.UserID).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch account", http.StatusInternalServerError)
		return
	}
	var req types.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.BaseCurrency == "" || req.Type == "" {
		http.Error(w, "base_currency and type are required", http.StatusBadRequest)
		return
	}
	tx := models.Transaction{
		Amount:           req.Amount,
		BaseCurrency:     req.BaseCurrency,
		Type:             req.Type,
		ShortDescription: req.ShortDescription,
		UserID:           claims.UserID,
		AccountID:        uint(accountID),
	}
	if err := db.Create(&tx).Error; err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}
