package handlers

import (
	"net/http"
	"os"

	"github.com/Tylerchristensen100/object_browser/internal"
)

func StaticFiles(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		http.ServeFileFS(w, r, os.DirFS("./static"), r.URL.Path)
	}
}
