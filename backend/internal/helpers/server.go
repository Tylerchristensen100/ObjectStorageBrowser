package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func Json(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("helpers/Json: Failed to encode JSON response", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func BucketFromQuery(r *http.Request) *string {
	bucket := r.URL.Query().Get("bucket")
	if bucket == "" {
		return nil
	}
	return &bucket
}

func PathFromQuery(r *http.Request) *string {
	path := r.URL.Query().Get("path")
	if path == "" {
		return nil
	}
	return &path
}

func ServerError(logger *slog.Logger, res http.ResponseWriter, req http.Request, err error) {
	var (
		method = req.Method
		uri    = req.URL.RequestURI()
		// trace  = string(debug.Stack())
	)
	logger.Error("helpers/ServerError: "+err.Error(), slog.String("method", method), slog.String("uri", uri))
	// logger.Error("helpers/ServerError: "+err.Error(), slog.String("method", method), slog.String("uri", uri), slog.String("stack_trace", trace))
	data, err := writeHTTPError(http.StatusInternalServerError, err.Error(), &req.Host)
	if err != nil {
		slog.Error("helpers/ServerError: Failed to marshal error response", slog.String("error", err.Error()))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusInternalServerError)
	if _, err := res.Write(data); err != nil {
		slog.Error("helpers/ServerError: Failed to write error response", slog.String("error", err.Error()))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func ClientError(res http.ResponseWriter, message string, status int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	data, err := writeHTTPError(status, message, nil)
	if err != nil {
		slog.Error("helpers/ClientError: Failed to marshal error response", slog.String("error", err.Error()))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := res.Write(data); err != nil {
		slog.Error("helpers/ClientError: Failed to write error response", slog.String("error", err.Error()))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

const DocsRoute = "/docs"

func writeHTTPError(status int, message string, host *string) ([]byte, error) {
	e := errorResponse{
		StatusCode: status,
		Message:    message,
	}
	if host != nil {
		e.Docs = *host + DocsRoute
	} else {
		e.Docs = "." + DocsRoute
	}

	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type errorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Docs       string `json:"docs"`
}
