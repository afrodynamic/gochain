package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	chainv1 "github.com/afrodynamic/gochain/api/proto/chain/v1"
	walletv1 "github.com/afrodynamic/gochain/api/proto/wallet/v1"
)

func newHealthHandler(grpcAddress string) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, _ *http.Request) {
		contextWithTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		grpcClient, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			responseWriter.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(responseWriter).Encode(map[string]any{
				"status":    "unhealthy",
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			})
			return
		}

		defer grpcClient.Close()

		healthClient := healthpb.NewHealthClient(grpcClient)
		_, err = healthClient.Check(contextWithTimeout, &healthpb.HealthCheckRequest{Service: ""})
		status := "ok"
		statusCode := http.StatusOK

		if err != nil {
			status = "unhealthy"
			statusCode = http.StatusServiceUnavailable
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(statusCode)
		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"status":    status,
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		})
	}
}

func NewMux(grpcAddress string) (http.Handler, error) {
	baseContext := context.Background()
	gatewayMux := runtime.NewServeMux()
	grpcOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := chainv1.RegisterChainHandlerFromEndpoint(baseContext, gatewayMux, grpcAddress, grpcOptions); err != nil {
		return nil, err
	}

	if err := walletv1.RegisterWalletHandlerFromEndpoint(baseContext, gatewayMux, grpcAddress, grpcOptions); err != nil {
		return nil, err
	}

	mainMux := http.NewServeMux()
	mainMux.Handle("/", gatewayMux)
	mainMux.Handle("/health", newHealthHandler(grpcAddress))

	return mainMux, nil
}

func Listen(httpAddress, grpcAddress string) error {
	mux, err := NewMux(grpcAddress)

	if err != nil {
		return err
	}

	return http.ListenAndServe(httpAddress, mux)
}
