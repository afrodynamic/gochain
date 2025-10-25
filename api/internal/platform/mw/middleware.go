package mw

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func RequestLogger(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			startTime := time.Now()

			nextHandler.ServeHTTP(responseWriter, request)

			logger.Info().
				Str("method", request.Method).
				Str("path", request.URL.Path).
				Dur("duration", time.Since(startTime)).
				Msg("HTTP request processed")
		})
	}
}

func Recoverer(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			defer func() {
				if recoveredValue := recover(); recoveredValue != nil {
					logger.Error().
						Interface("panic", recoveredValue).
						Str("method", request.Method).
						Str("path", request.URL.Path).
						Msg("panic recovered during request")

					http.Error(responseWriter, "internal server error", http.StatusInternalServerError)
				}
			}()

			nextHandler.ServeHTTP(responseWriter, request)
		})
	}
}
