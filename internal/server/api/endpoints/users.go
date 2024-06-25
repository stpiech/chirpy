package endpoints

import (
	"encoding/json"
	"net/http"

  "golang.org/x/crypto/bcrypt"
	"github.com/stpiech/chirpy/internal/database"
	"github.com/stpiech/chirpy/internal/server/helpers"
)

func LoginUser(w http.ResponseWriter, req *http.Request) { 
  var params database.User 
  err := json.NewDecoder(req.Body).Decode(&params)
  if err != nil {
    helpers.RespondWithError(w, 500, "Could not decode params")
    return
  }

  data, err := database.Data()  
  if err != nil {
    helpers.RespondWithError(w, 500, "")
    return
  }

  for _, v := range data.Users {
    if v.Email == params.Email {
      err = bcrypt.CompareHashAndPassword([]byte(v.Password), []byte(params.Password))
      if err != nil {
        helpers.RespondWithError(w, 401, "Wrong password")
        return
      }

      v.Password = ""
      jsonChirp, err := json.Marshal(v) 
      if err != nil {
        helpers.RespondWithError(w, 500, "")
        return
      }

      w.Write(jsonChirp)
      return
    }
  }

  helpers.RespondWithError(w, 401, "Email not found")
}

func CreateUser(w http.ResponseWriter, req *http.Request) { 
  var params database.User  
  err := json.NewDecoder(req.Body).Decode(&params)
  if err != nil {
    helpers.RespondWithError(w, 500, "Could not decode params")
    return
  }

  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
  if err != nil {
    helpers.RespondWithError(w, 500, "")
  }

  params.Password = string(hashedPassword)

  record, err := database.WriteUser(params)
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
