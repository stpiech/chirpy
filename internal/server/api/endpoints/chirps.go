package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/stpiech/chirpy/internal/database"
	"github.com/stpiech/chirpy/internal/server/helpers"
)

func ShowChirp(w http.ResponseWriter, req *http.Request) {
  id := req.PathValue("chirpId")
  intId, err := strconv.Atoi(id)
  if err != nil {
    helpers.RespondWithError(w, 422, "Can't convert ID to int")
    return
  }

  data, err := database.Data()  
  if err != nil {
    helpers.RespondWithError(w, 500, "")
    return
  }

  for _, v := range data.Chirps {
    if v.Id == intId {
      jsonChirp, err := json.Marshal(v) 
      if err != nil {
        helpers.RespondWithError(w, 500, "")
        return
      }

      w.Write(jsonChirp)
      return
    }
  }

  w.WriteHeader(404)
}

func CreateChirp(w http.ResponseWriter, req *http.Request) {
  var params database.Chirp  

  err := json.NewDecoder(req.Body).Decode(&params)
  if err != nil {
    helpers.RespondWithError(w, 500, "Could not decode params")
    return
  }

  correctedBody, isValid := validateChirp(params) 

  if !isValid {
    helpers.RespondWithError(w, 422, "Chirp is not valid")
    return
  }

  params.Body = correctedBody
  record, err := database.WriteChirp(params)
  if err != nil {
    helpers.RespondWithError(w, 500, "")
    return
  }

  jsonRecord, err := json.Marshal(record)
  if err != nil {
    helpers.RespondWithError(w, 500, "") 
    return
  }

  w.WriteHeader(201)
  w.Write(jsonRecord)
}

func IndexChirps(w http.ResponseWriter, req *http.Request) {
  data, err := database.Data()  
  if err != nil {
    helpers.RespondWithError(w, 500, "")
    return
  }

  chirps := data.Chirps

  jsonChirps, err := json.Marshal(chirps)
  if err != nil {
    helpers.RespondWithError(w, 500, "") 
    return
  }

  w.Write(jsonChirps)
}

func validateChirp(params database.Chirp) (string, bool) {
    if len(params.Body) > 140 {
      return "", false
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
    
    return strings.Join(words, " "), true
  }
