package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/blacktag/chirpy-project/internal/auth"
	"github.com/blacktag/chirpy-project/internal/database"
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
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coudnt decode parametrs")
		return
	}

	// _, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	// if err == nil {
	// 	respondWithError(w, http.StatusConflict, " User alredy exist")
	// 	return
	// }

	

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:  req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("❌ DB error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldnt create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, createUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}