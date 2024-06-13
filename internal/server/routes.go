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
	resp["message"] = "Hello World"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
        return helpers.WriteJson(w, http.StatusInternalServerError, model.ApiError{Error: err.Error()})
	}
	_, err = w.Write(jsonResp)
	return err 
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// parse login request body
	var loginRequest model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		return helpers.WriteJson(w, http.StatusBadRequest, fmt.Sprintf("Unable to parse request body: %v", err))
	}

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
    return helpers.WriteJson(w, http.StatusNotImplemented, model.ApiError{Error: "route not implemented"})
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
