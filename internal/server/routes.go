package server

import (
	"auth-demo/internal/auth"
	"auth-demo/internal/database"
	"auth-demo/internal/helpers"
	"auth-demo/internal/middleware"
	"auth-demo/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("POST /login", getHandleFunc(s.handleLogin))
	r.HandleFunc("POST /signup", getHandleFunc(s.handleSignup))

	protected := http.NewServeMux()
	protected.HandleFunc("GET /hello", getHandleFunc(s.handleHello))
	protected.HandleFunc("GET /account/{id}", getHandleFunc(s.handleGetAccountById))
	protected.HandleFunc("DELETE /account/{id}", getHandleFunc(s.handleDeleteAccountById))
	protected.HandleFunc("GET /account", getHandleFunc(s.handleGetAccount))

	r.Handle("/protected/", http.StripPrefix("/protected", middleware.WithAuth(protected)))
	return r
}

func getHandleFunc(fn model.ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Printf("ERROR - %v\n", err)
			helpers.WriteJson(w, http.StatusBadRequest, model.ApiError{Error: err.Error()})
		}
	}
}

func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) error {
	resp := make(map[string]string)

	pwd := r.Header.Get("account_pwd_hash")
	log.Println(pwd)

	resp["message"] = "Hello World"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return helpers.WriteJson(w, http.StatusInternalServerError, model.ApiError{Error: err.Error()})
	}
	_, err = w.Write(jsonResp)
	return err
}

func permissionDenied(w http.ResponseWriter) error {
	return helpers.WriteJson(
		w,
		http.StatusUnauthorized,
		model.ApiError{Error: "permission denied"},
	)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// parse login request body
	var loginRequest model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Println("Error - unable to decode login request")
		return permissionDenied(w)
	}

	// get acc from database
	acc, err := s.db.GetAccountByUser(loginRequest.User)
	if err != nil {
		log.Printf("Error - unable to get account from db: %+v\n", err)
		return permissionDenied(w)
	}

	// authenticate password
	if !auth.CheckPasswordHash(loginRequest.Pwd, acc.PwdHash) {
		log.Println("Error - could not authenticate password")
		return permissionDenied(w)
	}

	// create token
	if err := auth.SetAuthCookie(w, acc); err != nil {
		log.Printf("Unable to set auth cookie: %+v", err)
		return helpers.WriteJson(
			w,
			http.StatusInternalServerError,
			model.ApiError{Error: "server error"},
		)
	}
	log.Printf("JWT created for `%s`\n", acc.User)

	return helpers.WriteJson(w, http.StatusOK, acc)
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) error {
	createAccReq := &model.SignupRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createAccReq); err != nil {
		return err
	}

	account, err := database.NewAccount(*createAccReq)
	if err != nil {
		return err
	}

	if err := s.db.CreateAccount(account); err != nil {
		return err
	}

	if err := auth.SetAuthCookie(w, account); err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, account)
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

	return helpers.WriteJson(w, http.StatusOK, acc)
}

func (s *Server) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.db.GetAccounts()
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, accounts)
}

func (s *Server) handleDeleteAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}

	if err := s.db.DeleteAccount(id); err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]int{"deleted": id})
}

func getId(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid account id `%s`", idStr)
	}
	return id, nil
}
