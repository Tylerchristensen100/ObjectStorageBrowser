package middleware

import (
	"net/http"
)

func headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Vary", "origin")
		res.Header().Add("Vary", "Access-Control-Request-Method")
		res.Header().Add("Vary", "Access-Control-Request-Headers")

		res.Header().Set("cross-origin-allow-when-cross-origin", "*")

		res.Header().Set("referrer-policy", "origin-when-cross-origin")
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		res.Header().Add("Server", "Go - http.ServeMux")
		next.ServeHTTP(res, req)
	})
}
