
package main

import (
	"encoding/json"
	"net/http"
	"log"
	"time"
	"github.com/google/uuid"
	
)
type createUserRequest struct {
	
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
        ID        uuid.UUID `json:"id"`
        Email     string    `json:"email"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }


// handler to create users
func (cfg *apiConfig)createUserHandler(w http.ResponseWriter, r *http.Request){
	// method check
	if r.Method != "POST"{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	
	w.Header().Set("Content-Type", "application/json")
	
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err!=nil {
		respondWithError(w, http.StatusBadRequest,"Invalid Request")
		return
	}
	// validate the input 
	if req.Email == ""  {
		respondWithError(w, http.StatusBadRequest, "Inavalid request")
		return
	}

	

	user, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("❌ DB error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldnt create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, createUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}