package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/object_store"
	"github.com/Tylerchristensen100/object_browser/internal/server"
	"github.com/joho/godotenv"
)

const configPath = "./etc/config.yaml"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	app := loadConfig()
	client, err := object_store.LoadConfig(app, app.Secrets.ObjectStore.Endpoint, app.Secrets.ObjectStore.AccessKey, app.Secrets.ObjectStore.SecretKey)
	if err != nil {
		slog.Error("main/main: Failed to initialize S3 client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	app.Store = client

	err = server.Serve(app)
	if err != nil {
		slog.Error("main/main: Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func loadConfig() *internal.App {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	config, err := internal.ConfigFromYaml(logger, configPath)
	if err != nil {
		slog.Error("main/loadConfig: " + err.Error())
		os.Exit(1)
	}

	slog.SetDefault(logger)
	app := &internal.App{
		Log:     *logger,
		Ctx:     context.Background(),
		Store:   nil,
		Config:  config,
		Secrets: internal.LoadSecrets(),
	}

	app.Auth = internal.LoadOauthConfig(app)

	return app
}
