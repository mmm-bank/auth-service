package main

import (
	"github.com/mmm-bank/auth-service/http"
	"github.com/mmm-bank/auth-service/storage"
	"log"
	"os"
)

func main() {
	addr := ":8080"
	p := storage.NewUserPostgres(os.Getenv("POSTGRES_URL"))
	r := storage.NewSessionRedis(os.Getenv("REDIS_ADDR"))
	jwtKey := os.Getenv("JWT_SECRET")
	server := http.NewAuthService(p, r, jwtKey)

	log.Printf("Authentication server is running on port %s...", addr[1:])
	if err := http.CreateAndRunServer(server, addr); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
