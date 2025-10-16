package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/Tylerchristensen100/object_browser/internal"
	"github.com/Tylerchristensen100/object_browser/internal/helpers"
)

func HealthCheck(app *internal.App) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		runtime.GC()

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		status := map[string]interface{}{
			"status": "ok",
			"memory": map[string]interface{}{
				"heapAlloc": fmt.Sprintf("%d KB", memStats.Alloc/1024),
				"numGC":     memStats.NumGC,
			},
			"goRoutines": runtime.NumGoroutine(),
			"goVersion":  runtime.Version(),
		}

		data, err := json.Marshal(status)
		if err != nil {
			helpers.ServerError(&app.Log, res, *req, err)
			return
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.Write(data)
		res.WriteHeader(http.StatusOK)
	}
}
