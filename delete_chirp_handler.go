package main

import (
	"net/http"

	"github.com/blacktag/chirpy-project/internal/auth"
	"github.com/google/uuid"
)




func (cfg *apiConfig) deleteChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid access token")
		return
	}
	

	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unaiuthorized")
		return
	}
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp does not exists")
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you can only delete your own chirp")
		return
	}

	err = cfg.db.DeletChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot delete chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}