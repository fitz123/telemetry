package telemetry

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var httpMetrics = NewScope("http")

type record struct {
	Version    string
	SType      string
	Origin     string
	Channel    string
	Ext        string
	CacheState string
}

func NewRecord(r *http.Request) *record {
	return &record{
		Version: chi.URLParam(r, "version"),
		SType:   chi.URLParam(r, "stype"),
		Origin:  chi.URLParam(r, "origin"),
		Channel: chi.URLParam(r, "channel"),
		Ext:     filepath.Ext(chi.URLParam(r, "ext")),
	}
}

func sample(start time.Time, r *http.Request, ww middleware.WrapResponseWriter) {
	status := ww.Status()
	if status == 0 { // TODO: see why we have status = 0 under test conditions (this came up during benchmarks for some of the requests)
		status = http.StatusOK
	}
	// Retrieve cache state from response headers
	cacheState := ww.Header().Get("X-Cache")
	if cacheState == "" {
		cacheState = "unknown"
	}
	record := NewRecord(r)
	labels := map[string]string{
		"status":  fmt.Sprintf("%d", status),
		"version": record.Version,
		"stype":   record.SType,
		"origin":  record.Origin,
		"channel": record.Channel,
		"ext":     record.Ext,
		"cache":   cacheState,
	}

	httpMetrics.RecordDuration("request", labels, start, time.Now().UTC())
	httpMetrics.RecordHit("requests", labels)
}
