package helpers

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
  type errorType struct {
    Message string `json:"message"`
  }
  
  RespondWithJSON(w, code, errorType { msg })
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)

  jsonMsg, err := json.Marshal(payload)
  
  if err != nil {
    w.WriteHeader(500)
    return
  }

  w.Write(jsonMsg)
}

