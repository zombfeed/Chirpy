package main

import (
	"encoding/json"
	"net/http"

	"github.com/zombfeed/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password", err)
		return
	}
	pwValid, err := auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not validate password", err)
		return
	}
	if !pwValid {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, convertDBToUser(user))
}
