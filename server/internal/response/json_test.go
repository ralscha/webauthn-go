package response

import (
	"errors"
	"net/http"
	"testing"
)

type failingResponseWriter struct {
	header http.Header
	writes int
}

func (w *failingResponseWriter) Header() http.Header {
	return w.header
}

func (w *failingResponseWriter) WriteHeader(_ int) {}

func (w *failingResponseWriter) Write(_ []byte) (int, error) {
	w.writes++
	return 0, errors.New("connection closed")
}

func TestJSONDoesNotRetryAfterResponseWriteFailure(t *testing.T) {
	w := &failingResponseWriter{header: make(http.Header)}

	if JSON(w, http.StatusOK, map[string]string{"status": "ok"}) {
		t.Fatal("JSON returned true after a write failure")
	}
	if w.writes != 1 {
		t.Fatalf("response writes = %d, want 1", w.writes)
	}
}
