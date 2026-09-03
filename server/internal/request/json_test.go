package request

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSONRejectsOversizedBody(t *testing.T) {
	body := `{"value":"` + strings.Repeat("x", int(MaxBodyBytes)) + `"}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	var dst struct {
		Value string `json:"value"`
	}
	err := DecodeJSON(recorder, req, &dst)
	if err == nil || !strings.Contains(err.Error(), "must not be larger") {
		t.Fatalf("error = %v, want body-size error", err)
	}
}

func TestDecodeJSONRejectsMultipleValues(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"value":"first"} {"value":"second"}`))
	recorder := httptest.NewRecorder()

	var dst struct {
		Value string `json:"value"`
	}
	err := DecodeJSON(recorder, req, &dst)
	if err == nil || !strings.Contains(err.Error(), "single JSON value") {
		t.Fatalf("error = %v, want multiple-value error", err)
	}
}
