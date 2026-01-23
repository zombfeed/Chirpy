package main

import (
	"net/http"
	"time"

	"github.com/zombfeed/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "no auth header found", err)
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token does not exist or is expired", err)
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create new access token", err)
	}
	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}
