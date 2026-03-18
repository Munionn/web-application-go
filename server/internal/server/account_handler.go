package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"webapplication/internal/models"
	"webapplication/internal/types"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) getAllAccountsForUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserContextKey).(*types.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db := s.db.GetDB()
	var accounts []models.Account

	result := db.Joins("JOIN user_accounts ON user_accounts.account_id = accounts.id").
		Where("user_accounts.user_id = ?", claims.UserID).
		Find(&accounts)
	if result.Error != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (s *Server) getAccountForUserById(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserContextKey).(*types.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())
	accountID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	db := s.db.GetDB()
	var account models.Account

	result := db.Joins("JOIN user_accounts ON user_accounts.account_id = accounts.id").
		Where("accounts.id = ? AND user_accounts.user_id = ?", accountID, claims.UserID).
		First(&account)
	if result.Error != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (s *Server) createAccount(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserContextKey).(*types.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req types.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Account name is required", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := s.db.GetDB().First(&user, claims.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	account := models.Account{
		Name:  req.Name,
		Users: []*models.User{&user},
	}
	if err := s.db.GetDB().Create(&account).Error; err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (s *Server) updateAccount(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserContextKey).(*types.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	var account models.Account
	if err := s.db.GetDB().
		Joins("JOIN user_accounts ON user_accounts.account_id = accounts.id").
		Where("accounts.id = ? AND user_accounts.user_id = ?", id, claims.UserID).
		First(&account).Error; err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	var req types.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.BaseCurrence != nil {
		updates["base_currence"] = *req.BaseCurrence
	}
	if req.Balance != nil {
		updates["balance"] = *req.Balance
	}

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	if err := s.db.GetDB().Model(&account).Updates(updates).Error; err != nil {
		http.Error(w, "Failed to update account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.UserContextKey).(*types.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	var account models.Account
	if err := s.db.GetDB().
		Joins("JOIN user_accounts ON user_accounts.account_id = accounts.id").
		Where("accounts.id = ? AND user_accounts.user_id = ?", id, claims.UserID).
		First(&account).Error; err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	if err := s.db.GetDB().Delete(&account).Error; err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
