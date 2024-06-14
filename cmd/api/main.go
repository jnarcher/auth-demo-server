package main

import (
	"auth-demo/internal/database"
	"auth-demo/internal/server"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	store, err := database.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	server := server.NewServer(port, store)
	server.Run()
}
