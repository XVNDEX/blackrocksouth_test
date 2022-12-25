package main

import (
	"log"
	"net"

	"github.com/XVNDEX/blackrocksouth_test/internal/repository"
	"github.com/XVNDEX/blackrocksouth_test/internal/service"
	"github.com/XVNDEX/blackrocksouth_test/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	creds, err := repository.NewCredentials()
	if err != nil {
		log.Fatalf("failed to create new users: %v", err)
	}

	authServer := service.NewAuthServer(creds)
	serv := grpc.NewServer()
	pb.RegisterAuthServiceServer(serv, authServer)

	lis, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	if err := serv.Serve(lis); err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}
