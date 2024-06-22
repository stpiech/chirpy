package server

import (
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

  mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte("OK"))
  })

  mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte(fmt.Sprintf("Hits: %d", apiCfg.fileserverHits)))
  })

  mux.HandleFunc("/reset", func(w http.ResponseWriter, req *http.Request) {
    apiCfg.fileserverHits = 0
    w.WriteHeader(200)
  })

  server := http.Server { Handler: mux, Addr: ":8080" } 
  log.Fatal(server.ListenAndServe())
}
