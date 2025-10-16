package templates

import (
	"net/http"
	"text/template"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/helpers"
	"github.com/Tylerchristensen100/object_browser/internal/object_store"
	"github.com/Tylerchristensen100/object_browser/internal/server/handlers/api"
)

func DirectoryTemplate(app *internal.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const templateName = "directory"
		tmpl, err := template.New(templateName).ParseFiles("internal/templates/directory.html")
		if err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}
		bucket, path, err := api.BucketAndPathFromQuery(r)
		if err != nil {
			helpers.ClientError(w, err.Error(), http.StatusBadRequest)
			return
		}

		recursive := r.URL.Query().Get("recursive") == "true"

		listing, err := app.Store.ListDirectory(app, bucket, path, recursive)
		if err != nil {
			helpers.ServerError(app.Logger(), w, *r, err)
			return
		}

		dir, err := app.Store.Directory(app, bucket, path)

		data := struct {
			Listings  []object_store.Listing
			Directory *object_store.DirectoryItem
		}{
			Listings:  listing,
			Directory: dir,
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
