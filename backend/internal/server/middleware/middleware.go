package middleware

import (
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
)

func Stack(app *internal.App, server http.Handler) http.Handler {
	return logging(app, cors(app, headers(server)))
}
