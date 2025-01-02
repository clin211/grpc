package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	lb "github.com/clin211/grpc/load-balance/rpc"
)

var (
	addrs = []string{":50051", ":50052", ":50053"}
)

type HelloServer struct {
	lb.UnimplementedHelloServiceServer
	addr string
}

func (s *HelloServer) SyaHello(ctx context.Context, req *lb.HelloRequest) (*lb.HelloResponse, error) {
	message := fmt.Sprintf("Hello %s , form %s", req.GetName(), s.addr)
	return &lb.HelloResponse{Message: message}, nil
}

func startServer(addr string) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	lb.RegisterHelloServiceServer(s, &HelloServer{addr: addr})
	log.Printf("server listening at %v", addr)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	var wg sync.WaitGroup

	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			startServer(addr)
		}(addr)
	}

	wg.Wait()
}
