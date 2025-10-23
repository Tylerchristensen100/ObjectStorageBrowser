package api

import (
	"net/http"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/helpers"
)

func GetUser(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"mock":       true,
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john.doe@example.com",
			"roles":      []string{"user", "admin"},
		}

		w.WriteHeader(http.StatusOK)

		helpers.Json(w, http.StatusOK, data)
	}
}
