package main

import (
	"context"
	"log"
	"time"

	pb "github.com/clin211/grpc/service-types/go/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 建立连接
	conn, err := grpc.NewClient("localhost:6001",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := pb.NewUserServiceClient(conn)

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 发起 Unary RPC 调用
	resp, err := client.GetUserInfo(ctx, &pb.UserRequest{
		UserId: "123",
	})
	if err != nil {
		log.Fatalf("GetUserInfo failed: %v", err)
	}

	log.Printf("User Info: %+v", resp)
}
