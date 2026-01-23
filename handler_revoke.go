package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/zombfeed/Chirpy/internal/auth"
	"github.com/zombfeed/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "no auth header found", err)
		return
	}
	dbToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "refresh token not found", err)
		return
	}

	cfg.dbQueries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		Token:     dbToken.Token,
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: time.Now().UTC(),
	})
	w.WriteHeader(http.StatusNoContent)
}
