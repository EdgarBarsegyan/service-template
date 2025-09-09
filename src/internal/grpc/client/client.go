package client

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	greeterpb "service-template/pkg/proto"
)

var (
	addr = flag.String("addr", "localhost:50051", "server address")
	name = flag.String("name", "World", "name to greet")
	ttl  = flag.Duration("ttl", 3*time.Second, "per-RPC deadline")
)

func SendMessage() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	c := greeterpb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), *ttl)
	defer cancel()

	resp, err := c.SayHello(ctx, &greeterpb.HelloRequest{Name: *name})
	if err != nil {
		st, _ := status.FromError(err)
		log.Fatalf("rpc error: %v", st.Message())
	}

	log.Printf("Greeting: %s", resp.GetMessage())
}
