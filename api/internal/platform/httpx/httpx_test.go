package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusCreated, map[string]string{"x": "y"})
	if w.Code != http.StatusCreated {
		t.Fatalf("code=%d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("ct=%s", ct)
	}
	if !strings.Contains(w.Body.String(), `"x":"y"`) {
		t.Fatalf("body=%s", w.Body.String())
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequest(w, errors.New("bad"))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("code=%d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"error":"bad"`) {
		t.Fatalf("body=%s", w.Body.String())
	}
}
