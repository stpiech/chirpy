package server

import (
	"log"
	"net/http"

	"github.com/stpiech/chirpy/internal/server/api/endpoints"
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

  mux.HandleFunc("GET /api/healthz", endpoints.Healthz)
  mux.HandleFunc("GET /api/reset", func(w http.ResponseWriter, req *http.Request) { endpoints.Reset(w, req, &apiCfg.fileserverHits) })

  mux.HandleFunc("POST /api/chirps", endpoints.CreateChirp)
  mux.HandleFunc("GET /api/chirps", endpoints.IndexChirps)
  mux.HandleFunc("GET /api/chirps/{chirpId}", endpoints.ShowChirp)

  mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, req *http.Request) { endpoints.Metrics(w, req, &apiCfg.fileserverHits) })

  server := http.Server { Handler: mux, Addr: ":8080" } 
  log.Fatal(server.ListenAndServe())
}

