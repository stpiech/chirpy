package server

import (
	"net/http"
)

func Listen() {
  mux := http.NewServeMux()
  server := http.Server { Handler: mux, Addr: ":8080" } 
  server.ListenAndServe()
}
