package middleware

import (
	"net/http"
	"strings"

	"github.com/Tylerchristensen100/object_browser/internal"
)

func cors(app *internal.App, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")

		if origin != "" {
			for _, allowedOrigin := range app.Config.TrustedOrigins {
				if origin == allowedOrigin || strings.HasSuffix(origin, allowedOrigin) {
					res.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		if req.Method == http.MethodOptions {
			res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			res.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(res, req)
	})
}
