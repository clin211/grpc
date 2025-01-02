package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	lb "github.com/clin211/grpc/load-balance/rpc"
)

const (
	exampleScheme      = "example"
	exampleServiceName = "lb.example.lin"
)

var addrs = []string{"localhost:50051", "localhost:50052", "localhost:50053"}

func main() {
	address := exampleScheme + ":///" + exampleServiceName
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	fmt.Println("没有用负载均衡")
	rpcHandler(conn)

	lbConn, err := grpc.NewClient(
		address,
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer lbConn.Close()

	fmt.Println("用负载均衡：")
	rpcHandler(lbConn)
}

func rpcHandler(conn *grpc.ClientConn) {
	c := lb.NewHelloServiceClient(conn)
	for i := 0; i < 10; i++ {
		resp, err := c.SayHello(context.TODO(), &lb.HelloRequest{Name: "clina"})
		if err != nil {
			log.Fatalf("could not greet: %s", err)
		}

		fmt.Printf("resp : %v", resp.Message)
	}
}

type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: addrs,
		},
	}
	r.start()
	return r, nil
}

func (*exampleResolverBuilder) Scheme() string {
	return exampleScheme
}

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
