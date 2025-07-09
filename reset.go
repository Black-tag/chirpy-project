package main

import (
	
	"net/http"
	
)
//  handler to reset metrics 
func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)

	if cfg.platform != "dev" {
        respondWithJSON(w, http.StatusForbidden, map[string]string{
            "error": "This endpoint is only available in development mode",
        })
        return
    }


	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset database")
		return 

	}
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Database reset successfully",
	})
}