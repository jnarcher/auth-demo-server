package server

import (
	"auth-demo/internal/database"
	"auth-demo/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	port int
    db database.DB
}

func NewServer(port int, db database.DB) *Server {
	return &Server{
		port: port,
        db: db,
	}
}

func (s *Server) Run() {
    router := s.RegisterRoutes()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      middleware.ApplyDefault(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Listening on port %d\n", s.port)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
