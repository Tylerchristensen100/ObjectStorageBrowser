package handlers

import (
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
	"golang.org/x/oauth2"
)

func Login(app *internal.App) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		url := app.Auth.Config.AuthCodeURL("state", oauth2.AccessTypeOnline)
		http.Redirect(res, req, url, http.StatusFound)
	}
}

func Callback(app *internal.App) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		if code == "" {
			http.Error(res, "Missing Authorization code", http.StatusBadRequest)
			return
		}

		token, err := app.Auth.Token(code)
		if err != nil {
			app.Log.Error("handlers/Callback: Failed to exchange token", "error", err)
			http.Error(res, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = app.Auth.SetAuthCookie(token, res)
		if err != nil {
			app.Log.Error("handlers/Callback: Failed to set auth cookie", "error", err)
			http.Error(res, "Failed to set auth cookie: "+err.Error(), http.StatusInternalServerError)
			return
		}

		path := req.URL.Query().Get("url")
		if path == "" {
			path = "/"
		}

		http.Redirect(res, req, path, http.StatusFound)
	}
}
