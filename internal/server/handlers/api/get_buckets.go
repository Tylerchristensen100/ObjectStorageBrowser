package api

import (
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/helpers"
)

func GetBuckets(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		buckets, err := app.Store.ListBuckets(app)
		if err != nil {
			http.Error(w, "Failed to list buckets: "+err.Error(), http.StatusInternalServerError)
			return
		}

		helpers.Json(w, http.StatusOK, buckets)
	}
}
