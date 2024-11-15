package main

import (
	"fmt"
	"github.com/TariqueNasrullah/otel-practice/internal/delivery/grpc/book"
	"github.com/TariqueNasrullah/otel-practice/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":8080"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterBookServiceServer(grpcServer, book.NewService())

	log.Printf("server listening at %v", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
