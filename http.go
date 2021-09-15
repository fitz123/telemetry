package telemetry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

var httpMetrics = NewNamespace("http")

func sample(start time.Time, r *http.Request, ww middleware.WrapResponseWriter) {
	status := ww.Status()
	if status == 0 { // TODO: see why we have status = 0 under test conditions (this came up during benchmarks for some of the requests)
		status = http.StatusOK
	}

	labels := map[string]string{
		"endpoint": fmt.Sprintf("%s %s", r.Method, r.URL.Path),
		"status":   fmt.Sprintf("%d", status),
	}

	httpMetrics.RecordDuration("request", labels, start, time.Now().UTC())
	httpMetrics.RecordHit("requests", labels)
}
