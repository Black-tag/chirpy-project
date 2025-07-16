package main

import (
	"net/http"
	"time"

	"github.com/blacktag/chirpy-project/internal/auth"
)





func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}
	refreshToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	if refreshToken.RevokedAt.Valid || time.Now().After(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "refresh token expired or revoked")
		return
	}

	newAccessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create access token")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newAccessToken,
	})



}


func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
    tokenString, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Missing refresh token")
        return
    }

    err = cfg.db.RevokeRefreshToken(r.Context(), tokenString)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
        return
    }

    w.WriteHeader(http.StatusNoContent)
}