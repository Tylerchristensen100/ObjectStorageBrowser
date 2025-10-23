package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/helpers"
)

func DeleteDirectory(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, path, err := BucketAndPathFromQuery(r)
		if err != nil {
			helpers.ClientError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !strings.HasSuffix(path, "/") {
			helpers.ClientError(w, "Invalid directory path", http.StatusBadRequest)
			return
		}

		deleted, err := app.Store.DeleteDirectory(app, bucket, path)
		if err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}

		helpers.Json(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("`%s` deleted successfully", path), "deleted_files": fmt.Sprintf("%d", deleted)})
	}
}
