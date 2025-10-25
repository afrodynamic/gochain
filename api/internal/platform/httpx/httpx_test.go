package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSON(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	JSON(responseRecorder, http.StatusCreated, map[string]string{"key": "value"})

	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf("unexpected status code: %d", responseRecorder.Code)
	}

	contentType := responseRecorder.Header().Get("Content-Type")

	if contentType != "application/json" {
		t.Fatalf("unexpected content type: %s", contentType)
	}

	if !strings.Contains(responseRecorder.Body.String(), `"key":"value"`) {
		t.Fatalf("unexpected body: %s", responseRecorder.Body.String())
	}
}

func TestBadRequest(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	BadRequest(responseRecorder, errors.New("bad request"))

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", responseRecorder.Code)
	}

	if !strings.Contains(responseRecorder.Body.String(), `"error":"bad request"`) {
		t.Fatalf("unexpected body: %s", responseRecorder.Body.String())
	}
}
