package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/stpiech/chirpy/internal/database"
	"github.com/stpiech/chirpy/internal/server/helpers"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(w http.ResponseWriter, req *http.Request) { 
  type permittedParams struct {
    Email string `json:"email"`
    Password string `json:"password"`
    ExpiresInSeconds int `json:"expires_in_seconds"`
  }

  type responseBody struct {
    Id int `json:"id"`
    Email string `json:"email"`
    Token string `json:"token"`
  }
  
  var params permittedParams

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

      jwtToken, err := helpers.IssueJwtToken(v.Id, params.ExpiresInSeconds)

      responseBodyData := responseBody {
        Id: v.Id,
        Email: v.Email,
        Token: jwtToken,
      }

      jsonChirp, err := json.Marshal(responseBodyData) 
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

func UpdateUser(w http.ResponseWriter, req *http.Request) {
  type permittedParams struct {
    Email string `json:"email"`
    Password string `json:"password"`
  }

  type responseBody struct {
    Id int `json:"id"`
    Email string `json:"email"`
  }

  var params permittedParams
  err := json.NewDecoder(req.Body).Decode(&params)
  if err != nil {
    helpers.RespondWithError(w, 500, "Could not decode params")
  }
  
  jwtToken := req.Header.Get("Authorization")
  jwtToken = strings.ReplaceAll(jwtToken, "Bearer ", "")

  userId, err := helpers.DecodeJwtSubject(jwtToken)
  if err != nil {
    helpers.RespondWithError(w, 401, "")
  }

  intUserId, err := strconv.Atoi(userId)
  if err != nil {
    helpers.RespondWithError(w, 500, "")
  }

  password, err := encryptedPassword(params.Password)
  if err != nil {
    helpers.RespondWithError(w, 500, "")
  }

  newUser := database.User{
    Id: intUserId,
    Email: params.Email,
    Password: password,
  }
  
  updatedUser, err := database.UpdateUser(newUser) 
  if err != nil {
    helpers.RespondWithError(w, 500, "")
  }
  
  jsonRecord, err := json.Marshal(responseBody{Id: updatedUser.Id, Email: updatedUser.Email })
  if err != nil {
    helpers.RespondWithError(w, 500, "")
  }

  w.Write(jsonRecord)
}

func CreateUser(w http.ResponseWriter, req *http.Request) { 
  var params database.User  
  err := json.NewDecoder(req.Body).Decode(&params)
  if err != nil {
    helpers.RespondWithError(w, 500, "Could not decode params")
    return
  }

  password, err := encryptedPassword(params.Password)
  if err != nil {
    helpers.RespondWithError(w, 500, "")
  } 

  params.Password = password

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


func encryptedPassword(password string) (string, error) {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
  if err != nil {
    return "", err 
  }

  return string(hashedPassword), nil
}
