package api

import (
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/helpers"
)

func GetDirectory(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, path, err := BucketAndPathFromQuery(r)
		if err != nil {
			helpers.ClientError(w, err.Error(), http.StatusBadRequest)
			return
		}

		recursive := r.URL.Query().Get("recursive") == "true"

		contents, err := app.Store.ListDirectory(app, bucket, path, recursive)
		if err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}

		helpers.Json(w, http.StatusOK, contents)
	}
}

func GetDirectoryTree(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, path, err := BucketAndPathFromQuery(r)
		if err != nil {
			helpers.ClientError(w, err.Error(), http.StatusBadRequest)
			return
		}

		contents, err := app.Store.Directory(app, bucket, path)
		if err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}

		helpers.Json(w, http.StatusOK, contents)
	}
}
