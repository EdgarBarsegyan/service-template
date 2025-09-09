package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	greeterpb "service-template/pkg/proto"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	msg := "Hello " + in.GetName()
	return &greeterpb.HelloReply{Message: msg}, nil
}

func RunGrpcServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	greeterpb.RegisterGreeterServer(srv, &server{})

	// Опционально: health-check
	hs := health.NewServer()
	healthpb.RegisterHealthServer(srv, hs)

	log.Println("gRPC listening on :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
