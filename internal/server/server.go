package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
      Error string `json:"error,omitempty"`
      Valid bool `json:"valid"`
    }

    w.Header().Set("Content-Type", "application/json")

    decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)

    if err != nil {
      log.Println(err)
      w.WriteHeader(500)
      return
    }

    if len(params.Body) > 140 {
      w.WriteHeader(400)
      data := response{ Error: "Chirp is too long", Valid: false }
      jsonData, err := json.Marshal(data)
      if err != nil {
        log.Println(err)
        w.WriteHeader(500)
        return
      }
      w.Write(jsonData)
      return
    }

    data := response { Valid: true }
    jsonData, err := json.Marshal(data)
    if err != nil {
      log.Println(err)
      w.WriteHeader(500)
      return
    }

    w.WriteHeader(200)
    w.Write(jsonData)
  })

  server := http.Server { Handler: mux, Addr: ":8080" } 
  log.Fatal(server.ListenAndServe())
}
