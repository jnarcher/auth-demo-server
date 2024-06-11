package server

import (
	"auth-demo/internal/middleware"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newServer := &Server{
		port: port,
	}

    router := newServer.RegisterRoutes()

    stack := middleware.CreateStack(
        middleware.Logging,
        middleware.Cors,
    )

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      stack(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

    return server
}
