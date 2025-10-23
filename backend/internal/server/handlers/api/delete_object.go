package api

import (
	"fmt"
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/helpers"
)

func DeleteObject(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, path, err := BucketAndPathFromQuery(r)
		if err != nil {
			helpers.ClientError(w, err.Error(), http.StatusBadRequest)
			return
		}

		exists, err := app.Store.FileExists(app, bucket, path)
		if err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}

		if !exists {
			helpers.ClientError(w, "File does not exist", http.StatusNotFound)
			return
		}

		err = app.Store.DeleteFile(app, bucket, path)
		if err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}

		helpers.Json(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("`%s` deleted successfully", path)})
	}
}

func BucketAndPathFromQuery(r *http.Request) (string, string, error) {
	bucket := r.URL.Query().Get("bucket")
	path := r.URL.Query().Get("path")

	if bucket == "" {
		return "", "", fmt.Errorf("missing bucket")
	}

	return bucket, path, nil
}
