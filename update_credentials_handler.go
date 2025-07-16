package main

import (
	"encoding/json"
	"net/http"

	"github.com/blacktag/chirpy-project/internal/auth"
	"github.com/blacktag/chirpy-project/internal/database"
)



func (cfg *apiConfig) updateCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing/malformed access token")
		return
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid acccess token")
		return
	}
	type updateReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	hashedpassword, err :=auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot hash the new password")
		return
	}

	err = cfg.db.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		Email: req.Email,
		HashedPassword: hashedpassword,
		ID:  userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update user credentials")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to fetch the updated user")
		return 
	}
	respondWithJSON(w, http.StatusOK, loginResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})

}