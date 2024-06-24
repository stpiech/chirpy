package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type apiConfig struct {
  fileserverHits int
}

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    next.ServeHTTP(w, req)
    cfg.fileserverHits = cfg.fileserverHits + 1
  })
}

func Listen() {
  mux := http.NewServeMux()
  apiCfg := apiConfig{ fileserverHits: 0 }

  mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./internal/server/static")))))

  mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte("OK"))
  })

  mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(200)
    html := "<html><body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>"
    w.Write([]byte(fmt.Sprintf(html, apiCfg.fileserverHits)))
  })

  mux.HandleFunc("/api/reset", func(w http.ResponseWriter, req *http.Request) {
    apiCfg.fileserverHits = 0
    w.WriteHeader(200)
  })

  mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, req *http.Request) {
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
      respondWithError(w, 500, "Something went wrong")
      return
    }

    if len(params.Body) > 140 {
      respondWithError(w, 400, "Chirp is too long")
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
    respondWithJSON(w, 200, data)
  })

  server := http.Server { Handler: mux, Addr: ":8080" } 
  log.Fatal(server.ListenAndServe())
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
  type errorType struct {
    Message string `json:"message"`
  }
  
  respondWithJSON(w, code, errorType { msg })
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)

  jsonMsg, err := json.Marshal(payload)
  
  if err != nil {
    w.WriteHeader(500)
    return
  }

  w.Write(jsonMsg)
}
