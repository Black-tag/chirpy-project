package main

import (
	"encoding/json"
	"net/http"
	"github.com/blacktag/chirpy-project/internal/auth"
	"time"
	"github.com/google/uuid"

	
	
)

type loginResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
}

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	type loginreq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var lreq loginreq
	if err := json.NewDecoder(r.Body).Decode(&lreq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
	}
	

	user, err := cfg.db.GetUserByEmail(r.Context(), lreq.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	erri := auth.CheckPassworhash(lreq.Password, user.HashedPassword)
	if erri != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		
	})
	
}