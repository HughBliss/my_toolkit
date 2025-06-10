package grpcerver

import (
	"fmt"
	zfg "github.com/chaindead/zerocfg"
	"github.com/hughbliss/my_toolkit/tracer"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

var (
	listenGroup = zfg.NewGroup("listen")
	listenHost  = zfg.Str("host", "0.0.0.0", "LISTEN_HOST", zfg.Group(listenGroup))
	listenPort  = zfg.Uint32("port", 11000, "LISTEN_PORT", zfg.Group(listenGroup))
)

func Init() *grpc.Server {
	return grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Timeout:           15 * time.Second,
			MaxConnectionAge:  5 * time.Minute,
			Time:              15 * time.Minute,
		}),

		grpc.StatsHandler(tracer.ServerTracePropagator()),
		grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(tracer.Provider))),
	)
}

func Listener() (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *listenHost, *listenPort))
	if err != nil {
		return nil, err
	}
	return listener, nil
}
