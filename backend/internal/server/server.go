package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
)

func Serve(app *internal.App) error {
	router := Routes(app)
	server := &http.Server{
		Addr:     app.Config.Location(),
		Handler:  router,
		ErrorLog: slog.NewLogLogger(app.Log.Handler(), slog.LevelError),
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		app.Log.Error("internal/Serve: Failed to Start Go Server. Could not make TCP connection", slog.String("location", server.Addr))
		return err
	}

	app.Log.Info(fmt.Sprintf("Starting Go Server on %d", app.Config.Port))

	err = server.Serve(listener)

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.Log.Info("Go Server Stopped", slog.String("location", server.Addr))

	return nil
}
