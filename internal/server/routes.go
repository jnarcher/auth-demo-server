package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
    r := http.NewServeMux()

    r.HandleFunc("/", s.helloWorldHandler)

    r.HandleFunc("POST /login", s.loginHandler)

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
    type LoginRequest struct {
        User string
        Pwd string
    }
    var loginRequest LoginRequest

    if err := json.NewDecoder(r.Body).Decode((&loginRequest)); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Login: %+v\n", loginRequest)
    w.WriteHeader(http.StatusUnauthorized)
}
