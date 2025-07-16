package main

import (
	"net/http"

	"encoding/json"

	"github.com/blacktag/chirpy-project/internal/auth"
	"github.com/blacktag/chirpy-project/internal/database"
	"github.com/google/uuid"
)



func (cfg *apiConfig) webhooksHandler(w http.ResponseWriter, r *http.Request) {
	type UpgradeEvent struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "have no api key")
		return
	}
	if apiKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	var req UpgradeEvent
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
	}
	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
	} 
	
	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user id or format")
		return
	}


	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), database.UpgradeUserToChirpyRedParams{
		ID: userID,
		IsChirpyRed: true,	
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to upgrade user")
		return
	}
	w.WriteHeader(http.StatusNoContent)

	

}