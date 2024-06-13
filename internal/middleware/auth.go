package middleware

import (
	"auth-demo/internal/auth"
	"auth-demo/internal/helpers"
	"auth-demo/internal/model"
	"errors"
	// "fmt"
	"net/http"
)

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				_ = helpers.WriteJson(
					w,
					http.StatusBadRequest,
					model.ApiError{Error: "Auth cookie not found"},
				)
			default:
				_ = helpers.WriteJson(
					w,
					http.StatusInternalServerError,
					model.ApiError{Error: "server error"},
				)
			}
			return
		}

		tokenString := c.Value

		_, err = auth.VerifyToken(tokenString)
		if err != nil {
			_ = helpers.WriteJson(
				w,
				http.StatusUnauthorized,
				model.ApiError{Error: "Invalid auth token"},
			)
			return
		}

		// id, ok := claims["id"].(int)
		// if !ok {
		//           _ = helpers.WriteJson(
		//               w,
		//               http.StatusUnauthorized,
		//               model.ApiError{Error: "Unable to parse id from auth token"},
		//           )
		// 	return
		// }
		//
		// user, ok := claims["user"].(string)
		// if !ok {
		//           _ = helpers.WriteJson(
		//               w,
		//               http.StatusUnauthorized,
		//               model.ApiError{Error: "Unable to parse user auth token"},
		//           )
		// 	return
		// }
		//
		// pwdHash, ok := claims["pwdHash"].(string)
		// if !ok {
		//           _ = helpers.WriteJson(
		//               w,
		//               http.StatusUnauthorized,
		//               model.ApiError{Error: "Unable to parse pwd auth token"},
		//           )
		// 	return
		// }
		//
		//       r.Header.Add("account_id", fmt.Sprintf("%d", id))
		//       r.Header.Add("account_user", user)
		//       r.Header.Add("account_pwdHash", pwdHash)
		next.ServeHTTP(w, r)
	})
}
