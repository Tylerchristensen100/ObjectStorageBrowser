package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/constants"
	"github.com/Tylerchristensen100/object_browser/internal/helpers"
)

func GetObject(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, path, err := BucketAndPathFromQuery(r)
		if err != nil {
			helpers.ClientError(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, err := app.Store.GetFile(app, bucket, path)
		if err != nil {
			if errors.Is(err, constants.ErrObjectNotFound) {
				helpers.ClientError(w, constants.ErrObjectNotFound.Error(), http.StatusNotFound)
				return
			}
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}

		w.WriteHeader(http.StatusOK)

		w.Header().Set("Content-Disposition", "attachment; filename=\""+path+"\"")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprint(len(file)))

		if _, err := w.Write(file); err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}
	}
}
