package main

import (
	"auth-demo/internal/server"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env: %v", err)
    }

	server := server.NewServer()

	log.Printf("Listening on localhost:%s\n", os.Getenv("PORT"))
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
