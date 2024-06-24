package endpoints

import "net/http"

func Healthz(w http.ResponseWriter, req *http.Request) {
  w.WriteHeader(200)
  w.Write([]byte("OK"))
}
