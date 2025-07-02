package main

import (
	"context"
	"log"
	"net"

	pb "github.com/clin211/grpc/service-types/go/rpc"
	"google.golang.org/grpc"
)

// UserService 实现
type userService struct {
	pb.UnimplementedUserServiceServer
}

// GetUserInfo 实现 Unary RPC
func (s *userService) GetUserInfo(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Printf("Received request for user ID: %s", req.UserId)

	// 模拟数据库查询
	user := &pb.UserResponse{
		UserId:   req.UserId,
		Username: "clin",
		Email:    "7674254@qq.com",
		Age:      18,
	}

	return user, nil
}

func main() {
	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册服务
	pb.RegisterUserServiceServer(server, &userService{})

	// 监听端口
	lis, err := net.Listen("tcp", ":6001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Server started on :6001")

	// 启动服务
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
