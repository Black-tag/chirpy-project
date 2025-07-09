package main

import (
	"net/http"

	"github.com/google/uuid"
)




func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	ChirpIDStr := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(ChirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp iD format")
	}
	singleChirp, err := cfg.db.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w,http.StatusNotFound, "chirp not found" )
		return
	}
	singlechirp := Chirp{
		ID: singleChirp.ID,
		CreatedAt: singleChirp.CreatedAt,
		UpdatedAt: singleChirp.UpdatedAt,
		Body: singleChirp.Body,
		UserID: singleChirp.ID,
	}

	respondWithJSON(w, http.StatusOK, singlechirp)
}