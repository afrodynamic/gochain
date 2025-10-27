package main

import (
	"log"
	"net/http"
	"os"

	"github.com/afrodynamic/gochain/api/internal/adapter"
	"github.com/afrodynamic/gochain/api/internal/adapter/bitcoin"
	"github.com/afrodynamic/gochain/api/internal/adapter/ethereum"
	goadapter "github.com/afrodynamic/gochain/api/internal/adapter/gochain"
	grpcapi "github.com/afrodynamic/gochain/api/internal/api/grpc"
	httpapi "github.com/afrodynamic/gochain/api/internal/api/http"
	"github.com/afrodynamic/gochain/api/internal/chain/gochain"
	"github.com/afrodynamic/gochain/api/internal/consensus/pow"
	"github.com/afrodynamic/gochain/api/internal/storage/pebble"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	dataPath := os.Getenv("GOCHAIN_DATA_PATH")

	store, err := pebble.New(dataPath)

	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	bc := gochain.New(pow.New(8), store)

	reg := adapter.NewRegistry()
	reg.Register("gochain", goadapter.NewAdapter(bc))
	reg.Register("ethereum", ethereum.NewAdapter())
	reg.Register("bitcoin", bitcoin.NewAdapter())
	adp, _ := reg.Get("gochain")

	cs := grpcapi.NewChain(bc)
	ws := grpcapi.NewWallet(adp)

	handler, gs, err := httpapi.NewHandler(cs, ws)

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/demo/blocks", withDemoCORS(newBlocksHandler(bc)))
	mux.Handle("/demo/transactions", withDemoCORS(newTransactionsHandler(bc)))
	mux.Handle("/", handler)

	log.Printf("listening on :%s (REST+gRPC-Web+health), chain=gochain", port)
	log.Fatal(http.Serve(httpapi.CreateTCPListener(":"+port), httpapi.CreateH2CHandler(mux, gs)))
}
