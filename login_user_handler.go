package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blacktag/chirpy-project/internal/auth"
	"github.com/blacktag/chirpy-project/internal/database"
	"github.com/google/uuid"
)

type loginResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	
}

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login handler secret:", cfg.secret)
	type loginreq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		
	}

	var lreq loginreq
	if err := json.NewDecoder(r.Body).Decode(&lreq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
	}
	// if lreq.ExpiresInSeconds == 0 || lreq.ExpiresInSeconds > 3600 {
	// 	lreq.ExpiresInSeconds = 3600

	// }
	

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
		time.Hour,
	)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to refresh token")
		return
	}

	refreshExpiresAt := time.Now().Add(60 * 24 * time.Hour)

	params := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: refreshExpiresAt,
		RevokedAt: sql.NullTime{},
		
	}

	err = cfg.db.CreateRefreshToken(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to store refreshed token")
		return
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token: token,
		RefreshToken:refreshToken,
		IsChirpyRed: user.IsChirpyRed,

		
	})
	
}