package main

import (
	"net/http"
	
	
)







func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps")
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, c := range dbChirps{
		chirps[i] = Chirp{
			ID: c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body: c.Body,
			UserID: c.UserID,
		}
	}
	respondWithJSON(w, http.StatusOK,chirps)
}