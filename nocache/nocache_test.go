package nocache

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func Test_NoCacheHeaders(t *testing.T) {
	responseHeaders := map[string]string{
		"Cache-Control": "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
		"Pragma":        "no-cache",
		"Expires":       time.Unix(0, 0).UTC().Format(http.TimeFormat),
	}

	m := gin.New()
	m.Use(NoCache())

	recorder := httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	for _, k := range etagHeaders {
		r.Header.Add(k, "value")
	}
	m.ServeHTTP(recorder, r)

	for key, value := range responseHeaders {
		if recorder.Header()[key][0] != value {
			t.Errorf("Missing header: %s", key)
		}
	}
}
