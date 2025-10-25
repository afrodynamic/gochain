package httpapi

import (
	"net"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func CreateH2CHandler(httpApplication http.Handler, _ *grpc.Server) http.Handler {
	http2Server := &http2.Server{}

	return h2c.NewHandler(httpApplication, http2Server)
}

func CreateTCPListener(address string) net.Listener {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		panic(err)
	}

	return listener
}
