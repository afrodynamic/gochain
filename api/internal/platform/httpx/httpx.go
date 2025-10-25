package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(responseWriter http.ResponseWriter, statusCode int, value any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	_ = json.NewEncoder(responseWriter).Encode(value)
}

func BadRequest(responseWriter http.ResponseWriter, err error) {
	JSON(responseWriter, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
}

func Panicf(format string, arguments ...any) {
	panic(formatString(format, arguments...))
}

func formatString(format string, arguments ...any) string {
	return fmt.Sprintf(format, arguments...)
}
