package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zombfeed/Chirpy/internal/auth"
)

const (
	defaultExpSeconds = 3600
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode parameters", err)
		return
	}
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
		return
	}
	expSeconds := defaultExpSeconds
	if params.ExpiresInSeconds != 0 {
		expSeconds = params.ExpiresInSeconds
		if expSeconds <= 0 || expSeconds > defaultExpSeconds {
			expSeconds = defaultExpSeconds
		}
	}

	type response struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
		Token     string `json:"token"`
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized user", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
		Token:     token,
	})
}
