package server

import (
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/server/handlers"
	"github.com/Tylerchristensen100/object_browser/internal/server/handlers/api"
	"github.com/Tylerchristensen100/object_browser/internal/server/handlers/docs"
	"github.com/Tylerchristensen100/object_browser/internal/server/handlers/templates"
	"github.com/Tylerchristensen100/object_browser/internal/server/middleware"
)

func Routes(app *internal.App) http.Handler {
	server := http.NewServeMux()

	server.HandleFunc("GET /docs", docs.Page(app))
	server.HandleFunc("GET /docs/config", docs.Config(app))
	server.HandleFunc("GET /healthz", handlers.HealthCheck(app))

	// API
	server.HandleFunc("GET /api/buckets", api.GetBuckets(app))
	server.HandleFunc("GET /api/directory", api.GetDirectory(app))
	server.HandleFunc("GET /api/directory/tree", api.GetDirectoryTree(app))
	server.HandleFunc("DELETE /api/directory", api.DeleteDirectory(app))
	server.HandleFunc("GET /api/object", api.GetObject(app))
	server.HandleFunc("POST /api/object", api.PostObject(app))
	server.HandleFunc("DELETE /api/object", api.DeleteObject(app))

	// Templates
	server.HandleFunc("GET /templates/directory", templates.DirectoryTemplate(app))

	server.HandleFunc("GET /static/*", handlers.StaticFiles(app))

	return middleware.Stack(app, server)
}
