package test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	service "grpc_demo/proto"
	"log"
	"net"
	"testing"
	"time"
)

// server is used to implement helloworld.GreeterServer.
type Hello struct {
	service.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *Hello) SayHello(ctx context.Context, in *service.HelloRequest) (*service.HelloReply, error) {
	log.Printf("Received: %v", in.GetTitle())
	return &service.HelloReply{Msg: "Hello " + in.GetTitle()}, nil
}

func TestHelloServer(T *testing.T) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":50051"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service.RegisterGreeterServer(s, &Hello{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func TestHelloClient(T *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &service.HelloRequest{Title: " my name is tom"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMsg())
}
