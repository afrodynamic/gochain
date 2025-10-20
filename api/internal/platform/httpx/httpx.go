package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func BadRequest(w http.ResponseWriter, err error) {
	JSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
}

func Panicf(format string, v ...any) {
	panic(fmtSprintf(format, v...))
}

func fmtSprintf(format string, v ...any) string {
	return fmt.Sprintf(format, v...)
}
