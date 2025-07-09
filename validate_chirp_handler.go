package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/blacktag/chirpy-project/internal/database"
	"github.com/google/uuid"
)

// type chirpResponse struct {
//     CleanedBody string `json:"cleaned_body"`
// }

// type chirprequest struct {
// 	Body string `json:"body"`
// }

// type errorResponse struct {
// 	Error string `json:"error"`
// }

// type validResponse struct {
// 	Valid bool `json:"valid"`
// }

type chirpInput struct {
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
	
}



func (cfg *apiConfig)createChirpHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return

	}

	// var req chirprequest
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&req)
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "Invalid Request Body")
	// 	return
	// }

	// // validate chirp 

	
	var req chirpInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	err := validateChirp(req.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cleanedBody := filterProfanity(req.Body)
	

	// success response
	// respondWithJSON(w, http.StatusOK, map[string]string{
	// 	"cleaned_body": cleanedBody,
	// })

	params := database.CreateChirpParams{
		Body : cleanedBody, 
		UserID: req.UserID,

	}
	
	chirp, err := cfg.db.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,"cannot create chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated,Chirp{
		ID:        chirp.ID,
    	CreatedAt: chirp.CreatedAt,
    	UpdatedAt: chirp.UpdatedAt,
    	Body:      chirp.Body,
    	UserID:    chirp.UserID,
	})
	
}

func validateChirp(body string) error {
    const maxChirpLength = 140
    if len(body) > maxChirpLength {
        return fmt.Errorf("chirp is too long")
    }
    return nil
}

func filterProfanity(text string) string {
    // Define profane words (case-insensitive)
    profaneWords := map[string]bool{
        "kerfuffle": true,
        "sharbert":  true,
        "fornax":    true,
    }

    // Split text into words while preserving punctuation as separate tokens
    words := strings.Fields(text)
    
    for i, word := range words {
        lowerWord := strings.ToLower(word)
        // Check if word (without punctuation) is profane
        if profaneWords[lowerWord] {
            words[i] = "****"
        }
    }
    
    return strings.Join(words, " ")
}