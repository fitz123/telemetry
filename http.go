package telemetry

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5/middleware"
)

var httpMetrics = NewScope("http")

func sample(start time.Time, r *http.Request, ww middleware.WrapResponseWriter) {
	status := ww.Status()
	if status == 0 { // TODO: see why we have status = 0 under test conditions (this came up during benchmarks for some of the requests)
		status = http.StatusOK
	}
	// prometheus errors if string is not uft8 encoded
	if !utf8.ValidString(r.URL.Path) {
		return
	}
	// Retrieve path without filename and extension from the request URL
	path, ext := generalizePath(r.URL.Path)
	if ext == "" {
		ext = "none"
	}
	// Retrieve cache state from response headers
	cacheState := ww.Header().Get("X-Cache")
	if cacheState == "" {
		cacheState = "unknown"
	}
	labels := map[string]string{
		"endpoint": fmt.Sprintf("%s %s", r.Method, path),
		"status":   fmt.Sprintf("%d", status),
		"ext":      ext,
		"cache":    cacheState,
	}

	httpMetrics.RecordDuration("request", labels, start, time.Now().UTC())
	httpMetrics.RecordHit("requests", labels)
}

func generalizePath(path string) (string, string) {
	lastSlashIndex := strings.LastIndex(path, "/")
	lastDotIndex := strings.LastIndex(path, ".")
	if lastDotIndex > lastSlashIndex {
		ext := path[lastDotIndex+1:] // Get the extension without the dot
		return path[:lastSlashIndex], ext
	}
	return path, "" // Return the original path and an empty extension if no extension is found
}
