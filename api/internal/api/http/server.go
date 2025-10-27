package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	chainv1 "github.com/afrodynamic/gochain/api/proto/chain/v1"
	walletv1 "github.com/afrodynamic/gochain/api/proto/wallet/v1"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")

		if origin != "" {
			responseWriter.Header().Set("Access-Control-Allow-Origin", origin)
			responseWriter.Header().Set("Vary", "Origin")
			responseWriter.Header().Set("Access-Control-Allow-Credentials", "true")

			requestedHeaders := request.Header.Get("Access-Control-Request-Headers")
			if requestedHeaders == "" {
				requestedHeaders = "Content-Type, Authorization, X-Requested-With, X-Grpc-Web, Grpc-Timeout, X-User-Agent, Grpc-Encoding, Grpc-Accept-Encoding"
			}

			responseWriter.Header().Set("Access-Control-Allow-Headers", requestedHeaders)
			responseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			responseWriter.Header().Set("Access-Control-Expose-Headers", "Grpc-Status, Grpc-Message")
		}

		if request.Method == http.MethodOptions {
			responseWriter.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}

func healthHandler() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		})
	}
}

func NewHandler(chainService chainv1.ChainServer, walletService walletv1.WalletServer) (http.Handler, *grpc.Server, error) {
	grpcServer := grpc.NewServer()

	chainv1.RegisterChainServer(grpcServer, chainService)
	walletv1.RegisterWalletServer(grpcServer, walletService)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	reflection.Register(grpcServer)

	grpcWebServer := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(string) bool { return true }))

	baseContext := context.Background()
	gatewayMux := runtime.NewServeMux()

	if err := chainv1.RegisterChainHandlerServer(baseContext, gatewayMux, chainService); err != nil {
		return nil, nil, err
	}

	if err := walletv1.RegisterWalletHandlerServer(baseContext, gatewayMux, walletService); err != nil {
		return nil, nil, err
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("/health", healthHandler())
	rootMux.Handle("/", http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		contentType := request.Header.Get("Content-Type")

		if request.ProtoMajor == 2 && strings.HasPrefix(contentType, "application/grpc") {
			grpcServer.ServeHTTP(responseWriter, request)

			return
		}

		if grpcWebServer.IsGrpcWebRequest(request) || grpcWebServer.IsAcceptableGrpcCorsRequest(request) || grpcWebServer.IsGrpcWebSocketRequest(request) {
			grpcWebServer.ServeHTTP(responseWriter, request)

			return
		}

		gatewayMux.ServeHTTP(responseWriter, request)
	}))

	_ = os.Setenv("GRPC_GO_REQUIRE_HANDSHAKE", "off")

	return enableCORS(rootMux), grpcServer, nil
}
