package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func DecodeJwtSubject(token string) (string, error) {
  jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
  }) 

  if err != nil {
    return "", err
  }

  if claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims); ok {
    return claims.Subject, nil
  }
  return "", nil
}

func IssueJwtToken(userId int, expires_in_seconds int) (string, error) {
  stringifiedUserId := fmt.Sprint(userId)

  if expires_in_seconds == 0 || expires_in_seconds > 86400 {
    expires_in_seconds = 86400
  }

  claims := jwt.RegisteredClaims {
    Issuer: "chirpy",
    IssuedAt: jwt.NewNumericDate(time.Now()),
    ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expires_in_seconds) * time.Second)),
    Subject: stringifiedUserId,
  }


  token := jwt.NewWithClaims(
    jwt.SigningMethodHS256,
    claims,
  )

  ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
  if err != nil {
    return "", err
  }
  
  return ss, nil
}
