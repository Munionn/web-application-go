package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"webapplication/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	db := s.db.GetDB()
	var users []models.User

	result := db.Where("active = ?", true).Find(&users)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	go func() {
		log.Printf("User fetched: %d users", len(users))
	}()
}

// getUserHandler retrieves a single user by ID
func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := s.db.GetDB()
	var user models.User

	result := db.First(&user, id)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// createUserHandler creates a new user
func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := s.db.GetDB()
	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	go func() {
		account := models.Account{
			Name:  "Default Account",
			Users: []*models.User{&user},
		}
		if err := s.db.GetDB().Create(&account).Error; err != nil {
			log.Printf("Failed to create default account for user %d: %v", user.ID, err)
		} else {
			log.Printf("Default account created for user %d", user.ID)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := s.db.GetDB()
	var user models.User

	// First, find the user
	if err := db.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Decode update data
	var updateData models.User
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the user
	result := db.Model(&user).Updates(updateData)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// deleteUserHandler soft deletes a user
func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := s.db.GetDB()
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	result := db.Delete(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
