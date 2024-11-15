package main

import (
	"context"
	"fmt"
	"github.com/TariqueNasrullah/otel-practice/internal/delivery/grpc/book"
	"github.com/TariqueNasrullah/otel-practice/otel"
	"github.com/TariqueNasrullah/otel-practice/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	shutdownFunc, err := otel.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = shutdownFunc(context.Background())
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":8080"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterBookServiceServer(grpcServer, book.NewService())

	srvErr := make(chan error, 1)
	go func() {
		log.Printf("server listening at %v", lis.Addr())
		srvErr <- grpcServer.Serve(lis)
	}()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	metricSrv := &http.Server{
		Handler:     mux,
		Addr:        ":8081",
		BaseContext: func(listener net.Listener) context.Context { return ctx },
	}

	metricErr := make(chan error, 1)
	go func() {
		log.Printf("metric server listening at %v", metricSrv.Addr)
		metricErr <- metricSrv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		return
	case err = <-metricErr:
		return
	case <-ctx.Done():
		stop()
	}

	grpcServer.GracefulStop()

	err = metricSrv.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
