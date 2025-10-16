package docs

import (
	"net/http"
	"path/filepath"

	"github.com/Tylerchristensen100/object_browser/internal"
)

func Page(app *internal.App) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		html, err := filepath.Abs("etc/swagger/index.html")
		if err != nil {
			app.Log.Error("docs/Page: " + err.Error())
		}

		http.ServeFile(res, req, html)
	}
}

func Config(app *internal.App) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		yaml, err := filepath.Abs("etc/swagger/openapi.yaml")
		if err != nil {
			app.Log.Error("docs/Config: " + err.Error())
		}
		http.ServeFile(res, req, yaml)
	}
}
