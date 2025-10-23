package server

import (
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/server/handlers"
	"github.com/Tylerchristensen100/object_browser/internal/server/handlers/api"
	"github.com/Tylerchristensen100/object_browser/internal/server/handlers/docs"
	"github.com/Tylerchristensen100/object_browser/internal/server/middleware"
)

func Routes(app *internal.App) http.Handler {
	server := http.NewServeMux()

	auth := app.Auth

	server.HandleFunc("GET /docs", docs.Page(app))
	server.HandleFunc("GET /docs/config", docs.Config(app))
	server.HandleFunc("GET /healthz", handlers.HealthCheck(app))

	// API
	server.HandleFunc("GET /api/buckets", auth.Require(api.GetBuckets(app)))
	server.HandleFunc("GET /api/directory", auth.Require(api.GetDirectory(app)))
	server.HandleFunc("GET /api/directory/tree", auth.Require(api.GetDirectoryTree(app)))
	server.HandleFunc("DELETE /api/directory", auth.Require(api.DeleteDirectory(app)))
	server.HandleFunc("GET /api/object", auth.Require(api.GetObject(app)))
	server.HandleFunc("POST /api/object", auth.Require(api.PostObject(app)))
	server.HandleFunc("DELETE /api/object", auth.Require(api.DeleteObject(app)))

	server.HandleFunc("GET /api/user", auth.Require(api.GetUser(app)))

	server.HandleFunc("GET /login", handlers.Login(app))
	server.HandleFunc("GET /callback", handlers.Callback(app))

	// server.HandleFunc("GET /static/*", handlers.StaticFiles(app))

	return middleware.Stack(app, server)
}
