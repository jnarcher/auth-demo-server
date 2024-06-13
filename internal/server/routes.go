package server

import (
	"auth-demo/internal/middleware"
	"auth-demo/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ApiError struct {
    Error string `json:"error"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func (s *Server) RegisterRoutes() http.Handler {
	r := http.NewServeMux()

	public := http.NewServeMux()
	public.HandleFunc("POST /login", getHandleFunc(s.handleLogin))
	public.HandleFunc("POST /signup", getHandleFunc(s.handleSignup))

    // TODO: put these routes in protected
    public.HandleFunc("GET /account/{id}", getHandleFunc(s.handleGetAccountById))
    public.HandleFunc("DELETE /account/{id}", getHandleFunc(s.handleDeleteAccountById))
    public.HandleFunc("GET /account", getHandleFunc(s.handleGetAccount))

	protected := http.NewServeMux()
	protected.HandleFunc("GET /hello", getHandleFunc(s.handleHello))

	r.Handle("/public/", http.StripPrefix("/public", public))
	r.Handle("/protected/", http.StripPrefix(
		"/protected", middleware.Auth(protected),
	))
	return r
}

func getHandleFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
            log.Printf("ERROR - %v\n", err)
			writeJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) error {
	resp := make(map[string]string)
	resp["message"] = "Hello World"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
        return writeJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	_, err = w.Write(jsonResp)
	return err 
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// // parse login request body
	// var loginRequest model.LoginRequest
	// if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
	// 	return writeJson(w, http.StatusBadRequest, fmt.Sprintf("Unable to parse request body: %v", err))
	// }

	//
	// // get acc from database
	// acc, err := s.db.GetAccount(loginRequest.User)
	// if err != nil {
	// 	// hlp.SendError(w, "Invalid credentials", http.StatusUnauthorized)
	// 	return writeJson(w, http.StatusUnauthorized, ApiError{Error: "Invalid credentials"})
	// }
	//
	// // authenticate password
	// if !auth.CheckPasswordHash(loginRequest.Pwd, acc.PwdHash) {
	// 	return writeJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Unable to parse token: %v", err)})
	// }
	//
	// // create token
	// tokenString, err := auth.CreateToken(acc.User, acc.PwdHash, acc.Role)
	// if err != nil {
	// 	return writeJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Unable to parse token: %v", err)})
	// }
	//
	// log.Printf("JWT created (%s): %s\n", loginRequest.User, tokenString)
	//
	// // set token cookie
	// cookie := &http.Cookie{
	// 	Name:     "token",
	// 	Value:    tokenString,
	// 	HttpOnly: true,
	// 	Secure:   true,
	// }
	// http.SetCookie(w, cookie)
	//
	// return writeJson(w, http.StatusOK, acc.Safe()
    return writeJson(w, http.StatusNotImplemented, ApiError{Error: "route not implemented"})
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) error {
    createAccReq := &model.CreateAccountRequest{}
    if err := json.NewDecoder(r.Body).Decode(&createAccReq); err != nil {
        return err
    }

    account, err := model.NewAccount(*createAccReq)
    if err != nil {
        return err
    }

    if err := s.db.CreateAccount(account); err != nil {
        return err
    }

    return writeJson(w, http.StatusOK, account)
}

func (s *Server) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
    id, err := getId(r)
    if err != nil {
        return err
    }

	acc, err := s.db.GetAccountById(id)
	if err != nil {
        return fmt.Errorf("No account found with id %d", id)
	}

	return writeJson(w, http.StatusOK, acc)
}

func (s *Server) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
    accounts, err := s.db.GetAccounts()
    if err != nil {
        return err
    }
    return writeJson(w, http.StatusOK, accounts)
}


func (s *Server) handleDeleteAccountById(w http.ResponseWriter, r *http.Request) error {
    id, err := getId(r)
    if err != nil {
        return err
    }

    if err := s.db.DeleteAccount(id); err != nil {
        return err
    }

    return writeJson(w, http.StatusOK, map[string]int{"deleted": id})
}

func getId(r *http.Request) (int, error) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return 0, fmt.Errorf("invalid account id `%s`", idStr)
    }
    return id, nil
}
