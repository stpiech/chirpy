package server

import (
	"log"
	"net/http"
)

func Listen() {
  mux := http.NewServeMux()
  mux.Handle("/", http.FileServer(http.Dir("./internal/server/static")))
  server := http.Server { Handler: mux, Addr: ":8080" } 
  log.Fatal(server.ListenAndServe())
}
