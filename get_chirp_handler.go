package main

import (
	"net/http"
	"sort"

	"github.com/blacktag/chirpy-project/internal/database"
	"github.com/google/uuid"
)



func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authorID := r.URL.Query().Get("author_id")

	var dbChirps []database.Chirp
	var err error

	if authorID != "" {
		userID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author id format")
			return
		}
		dbChirps, err = cfg.db.GetChirpsByUserID(r.Context(), userID)
	}else {
		dbChirps, err = cfg.db.GetChirps(r.Context())
	}
	
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
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "asc" {
		sort.Slice(chirps, func(i int, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})

	} else if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})

	}
	respondWithJSON(w, http.StatusOK,chirps)
}