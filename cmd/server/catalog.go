package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"strings"

	"github.com/XVNDEX/blackrocksouth_test/internal/data"
	"github.com/XVNDEX/blackrocksouth_test/internal/repository"
	"github.com/XVNDEX/blackrocksouth_test/internal/service"
	"github.com/XVNDEX/blackrocksouth_test/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	catalog, err := repository.NewCatalog()
	if err != nil {
		log.Fatalf("failed to create catalog: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(data.Path("server_cert.pem"), data.Path("server_key.pem"))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
		// Enable TLS for all incoming connections.
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	catalogServer := service.NewCatalogServer(catalog)
	catalogGrpcServer := grpc.NewServer(opts...)
	pb.RegisterCatalogServer(catalogGrpcServer, catalogServer)

	log.Println("starting server on :9090 port")
	lis, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	if err := catalogGrpcServer.Serve(lis); err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}
}

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

// valid validates the authorization.
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "some-secret-token"
}
