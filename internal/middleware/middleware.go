package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}

func Apply(router http.Handler, stack Middleware) http.Handler {
    return stack(router)
}

func ApplyDefault(router http.Handler) http.Handler {
    return Apply(router, CreateStack(
        Cors,
        Logging,
    ))
}
