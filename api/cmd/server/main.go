package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/afrodynamic/gochain/api/internal/adapters/bitcoin"
	"github.com/afrodynamic/gochain/api/internal/adapters/ethereum"
	"github.com/afrodynamic/gochain/api/internal/adapters/gochain"
	"github.com/afrodynamic/gochain/api/internal/core"
	"github.com/afrodynamic/gochain/api/internal/http/handlers"
	"github.com/afrodynamic/gochain/api/internal/http/openapi"
	"github.com/afrodynamic/gochain/api/internal/platform/config"
	"github.com/afrodynamic/gochain/api/internal/platform/mw"
	"github.com/afrodynamic/gochain/api/internal/service/wallet"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	nethttp_middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/rs/zerolog"
)

func splitCommaEnv(key string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func main() {
	cfg := config.Load()
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	reg := core.NewRegistry(gochain.NewAdapter(), ethereum.NewAdapter(), bitcoin.NewAdapter())
	def, ok := reg.Get(cfg.Chain)
	if !ok {
		log.Fatal().Str("chain", cfg.Chain).Msg("unknown chain")
	}
	svc := wallet.NewService(def, reg)

	r := chi.NewRouter()
	r.Use(mw.RequestLogger(&log))
	r.Use(mw.Recoverer(&log))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   splitCommaEnv("CORS_ORIGINS"),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: os.Getenv("CORS_CREDENTIALS") == "1",
		MaxAge:           300,
	}))

	sw, err := openapi.GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("swagger")
	}
	sw.Servers = nil
	r.Use(nethttp_middleware.OapiRequestValidator(sw))

	h := handlers.NewHandler(svc)
	openapi.HandlerFromMux(h, r)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Info().Str("addr", cfg.Addr).Str("chain", cfg.Chain).Strs("cors", splitCommaEnv("CORS_ORIGINS")).Msg("listening")
	go func() { _ = srv.ListenAndServe() }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Info().Msg("shutdown complete")
}
