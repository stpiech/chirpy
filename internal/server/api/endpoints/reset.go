package endpoints

import "net/http"

func Reset(w http.ResponseWriter, req *http.Request, fileserverHits *int) {
    *fileserverHits = 0
    w.WriteHeader(200)
}
