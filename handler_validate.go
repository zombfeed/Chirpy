package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	filtered := profanityFilter(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: filtered})
}

func profanityFilter(body string) string {
	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	splitBody := strings.Split(body, " ")
	for i, word := range splitBody {
		lword := strings.ToLower(word)
		if _, ok := profanity[lword]; ok {
			splitBody[i] = "****"
		}
	}
	return strings.Join(splitBody, " ")
}
