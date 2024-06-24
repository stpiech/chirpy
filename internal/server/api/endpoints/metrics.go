package endpoints

import (
	"fmt"
	"net/http"
)

func Metrics(w http.ResponseWriter, req *http.Request, filServerHits *int) {
    w.WriteHeader(200)
    html := "<html><body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>"
    w.Write([]byte(fmt.Sprintf(html, *filServerHits)))
  } 
