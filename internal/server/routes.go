package server

import (
	hlp "auth-demo/internal/helpers"
	"auth-demo/internal/middleware"
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := http.NewServeMux()

	public := http.NewServeMux()
	public.HandleFunc("POST /login", s.loginHandler)

	protected := http.NewServeMux()
	protected.HandleFunc("GET /hello", s.helloWorldHandler)

	r.Handle("/public/", http.StripPrefix("/public", public))
	r.Handle("/protected/", http.StripPrefix(
		"/protected", middleware.Auth(protected),
	))
	return r
}

func (s *Server) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}
	_, _ = w.Write(jsonResp)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type LoginRequest struct {
        User string `json:"user"`
		Pwd  string `json:"pwd"`
	}

	var loginRequest LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if loginRequest.User == "test" && loginRequest.Pwd == "123456" {
		tokenString, err := middleware.CreateToken(
			loginRequest.User,
			loginRequest.Pwd,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("No username found"))
			return
		}

		w.WriteHeader(http.StatusOK)
        log.Printf("JWT created (%s): %s\n", loginRequest.User, tokenString)
		return
	} else {
		log.Println(loginRequest)
		hlp.SendError(w, "Invalid credentials", http.StatusUnauthorized)
	}
}
