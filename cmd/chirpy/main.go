package main

import (
	"github.com/joho/godotenv"
	"github.com/stpiech/chirpy/internal/server"
)

func main() {
  godotenv.Load()
  server.Listen()
}
