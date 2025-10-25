package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
)

func TestRequestLogger(t *testing.T) {
	logger := zerolog.Nop()

	handler := RequestLogger(&logger)(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("unexpected status code: %d", responseRecorder.Code)
	}
}

func TestRecoverer(t *testing.T) {
	logger := zerolog.Nop()

	handler := Recoverer(&logger)(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		panic("unexpected panic")
	}))

	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", responseRecorder.Code)
	}
}
