package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
)
// type chirpResponse struct {
//     CleanedBody string `json:"cleaned_body"`
// }

type chirprequest struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// type validResponse struct {
// 	Valid bool `json:"valid"`
// }



func (cfg *apiConfig)validateChirpHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")

	}

	var req chirprequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}

	// validate chirp 

	err = validateChirp(req.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cleanedBody := filterProfanity(req.Body)
	

	// success response
	respondWithJSON(w, http.StatusOK, map[string]string{
		"cleaned_body": cleanedBody,
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