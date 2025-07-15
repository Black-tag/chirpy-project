package main

import (
	"encoding/json"
	"net/http"
	"github.com/blacktag/chirpy-project/internal/auth"
	"time"
	"github.com/google/uuid"
	"fmt"

	
	
)

type loginResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Token      string    `json:"token"`
	
}

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login handler secret:", cfg.secret)
	type loginreq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	var lreq loginreq
	if err := json.NewDecoder(r.Body).Decode(&lreq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
	}
	if lreq.ExpiresInSeconds == 0 || lreq.ExpiresInSeconds > 3600 {
		lreq.ExpiresInSeconds = 3600

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
	token, err := auth.MakeJWT(
		user.ID,
		cfg.secret,
		time.Duration(lreq.ExpiresInSeconds)*time.Second,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token: token,

		
	})
	
}