package endpoints

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/stpiech/chirpy/internal/server/helpers"
)

func ValidateChirp(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
      Body string `json:"body"`
    }

    type response struct {    
      CleanedBody string `json:"cleaned_body"`
    }

    decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)

    if err != nil {
     helpers.RespondWithError(w, 500, "Something went wrong")
      return
    }

    if len(params.Body) > 140 {
      helpers.RespondWithError(w, 400, "Chirp is too long")
      return
    }

    bannedWords := map[string]bool {
      "kerfuffle": true,
      "sharbert": true,
      "fornax": true,
    }

    words := strings.Split(params.Body, " ")  
    for i, word := range words {
      _, wordBanned := bannedWords[strings.ToLower(word)]

      if wordBanned {
        words[i] = "****"
      } 
    }

    data := response { CleanedBody: strings.Join(words, " ") }
    helpers.RespondWithJSON(w, 200, data)
  }
