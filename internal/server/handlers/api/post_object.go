package api

import (
	"net/http"
	"strings"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/helpers"
)

func PostObject(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, path, err := BucketAndPathFromQuery(r)
		if err != nil {
			helpers.ClientError(w, err.Error(), http.StatusBadRequest)
			return
		}

		multipartFile, header, err := r.FormFile("file")
		if err != nil {
			helpers.ClientError(w, "Failed to get file from form: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer multipartFile.Close()

		filename := header.Filename
		filename = strings.TrimPrefix(filename, "/")
		filename = strings.ReplaceAll(filename, " ", "_")
		filename = strings.ReplaceAll(filename, "\\", "/")

		fullPath := path + filename

		exists, err := app.Store.FileExists(app, bucket, fullPath)
		if err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}

		if exists {
			helpers.ClientError(w, "File already exists", http.StatusConflict)
			return
		}

		info, err := app.Store.SaveFile(app, bucket, multipartFile, header, fullPath)
		if err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}

		helpers.Json(w, http.StatusOK, info)
	}
}
