package cmd

import (
	"encoding/json"
	"net/http"
	"regexp"
	"unicode"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func isValidUsername(username string) bool {
	if len(username) == 0 {
		return false
	}

	// Must not start with a digit
	if unicode.IsDigit(rune(username[0])) {
		return false
	}

	// Must contain at least one letter
	hasLetter := false
	for _, r := range username {
		if unicode.IsLetter(r) {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return false
	}

	// Can’t be just special characters
	onlySpecials := regexp.MustCompile(`^[^a-zA-Z0-9]+$`)
	return !onlySpecials.MatchString(username)
}
