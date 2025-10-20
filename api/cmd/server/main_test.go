package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/afrodynamic/gochain/api/internal/adapters/bitcoin"
	"github.com/afrodynamic/gochain/api/internal/adapters/ethereum"
	"github.com/afrodynamic/gochain/api/internal/adapters/gochain"
	"github.com/afrodynamic/gochain/api/internal/core"
	"github.com/afrodynamic/gochain/api/internal/http/handlers"
	"github.com/afrodynamic/gochain/api/internal/http/openapi"
	"github.com/afrodynamic/gochain/api/internal/platform/mw"
	"github.com/afrodynamic/gochain/api/internal/service/wallet"
	"github.com/go-chi/chi/v5"
	nethttp_middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestServer_BootAndShutdown(t *testing.T) {
	log := zerolog.Nop()
	reg := core.NewRegistry(gochain.NewAdapter(), ethereum.NewAdapter(), bitcoin.NewAdapter())
	svc := wallet.NewService(gochain.NewAdapter(), reg)

	r := chi.NewRouter()
	r.Use(mw.RequestLogger(&log))
	r.Use(mw.Recoverer(&log))

	sw, err := openapi.GetSwagger()
	assert.NoError(t, err)
	sw.Servers = nil
	r.Use(nethttp_middleware.OapiRequestValidator(sw))
	h := handlers.NewHandler(svc)
	openapi.HandlerFromMux(h, r)

	srv := &http.Server{Addr: ":0", Handler: r}
	go srv.ListenAndServe()

	done := make(chan struct{})
	go func() {
		defer close(done)
		time.Sleep(200 * time.Millisecond)
		srv.Shutdown(context.Background())
	}()
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("shutdown timeout")
	}
}
