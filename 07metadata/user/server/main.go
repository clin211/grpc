package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	rpc "github.com/clin211/grpc/metadata/proto"
)

// UserServer 实现用户服务
type UserServer struct {
	rpc.UnimplementedUserServiceServer
}

// getMetadataValue 获取元数据的第一个值
func getMetadataValue(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// printRequestMetadata 打印请求元数据
func printRequestMetadata(md metadata.MD) {
	fmt.Println("========== 收到的请求元数据 ==========")
	for key, values := range md {
		fmt.Printf("  %s: %v\n", key, values)
	}
	fmt.Println("=====================================")
}

// GetUser 获取用户信息
func (s *UserServer) GetUser(ctx context.Context, req *rpc.GetUserRequest) (*rpc.GetUserResponse, error) {
	startTime := time.Now()

	// 从context中提取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("警告：没有接收到元数据")
	} else {
		printRequestMetadata(md)

		// 验证认证信息
		authToken := getMetadataValue(md, "authorization")
		if authToken == "" {
			return nil, status.Error(codes.Unauthenticated, "缺少认证信息")
		}

		if !strings.HasPrefix(authToken, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "无效的认证格式")
		}

		// 获取其他元数据
		userAgent := getMetadataValue(md, "user-agent")
		clientVersion := getMetadataValue(md, "client-version")
		traceID := getMetadataValue(md, "x-trace-id")

		log.Printf("处理GetUser请求 - UserID: %s, TraceID: %s, ClientVersion: %s, UserAgent: %s",
			req.GetUserId(), traceID, clientVersion, userAgent)
	}

	// 发送头部元数据
	header := metadata.Pairs(
		"server-version", "1.0.0",
		"server-instance", "user-service-01",
		"processing-start", startTime.Format(time.RFC3339),
	)

	if err := grpc.SendHeader(ctx, header); err != nil {
		log.Printf("发送头部元数据失败: %v", err)
	}

	// 模拟业务逻辑处理
	time.Sleep(100 * time.Millisecond)

	// 构造响应
	response := &rpc.GetUserResponse{
		UserId:    req.GetUserId(),
		Username:  "张三",
		Email:     "zhangsan@example.com",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	// 设置尾部元数据
	processingTime := time.Since(startTime)
	trailer := metadata.Pairs(
		"processing-time", processingTime.String(),
		"records-found", "1",
		"cache-hit", "false",
	)
	grpc.SetTrailer(ctx, trailer)

	log.Printf("GetUser请求处理完成，耗时: %v", processingTime)

	return response, nil
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *rpc.CreateUserRequest) (*rpc.CreateUserResponse, error) {
	startTime := time.Now()

	// 提取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "缺少元数据")
	}

	printRequestMetadata(md)

	// 验证权限
	permissions := md.Get("x-permission")
	hasCreatePermission := false
	for _, perm := range permissions {
		if perm == "create" || perm == "admin" {
			hasCreatePermission = true
			break
		}
	}

	if !hasCreatePermission {
		return nil, status.Error(codes.PermissionDenied, "权限不足，无法创建用户")
	}

	// 获取请求来源信息
	clientIP := getMetadataValue(md, "x-client-ip")
	requestID := getMetadataValue(md, "x-request-id")

	log.Printf("处理CreateUser请求 - Username: %s, ClientIP: %s, RequestID: %s",
		req.GetUsername(), clientIP, requestID)

	// 发送头部元数据
	header := metadata.Pairs(
		"server-version", "1.0.0",
		"operation", "create_user",
	)
	grpc.SendHeader(ctx, header)

	// 模拟创建用户
	time.Sleep(200 * time.Millisecond)
	userID := fmt.Sprintf("user_%d", time.Now().Unix())

	// 设置尾部元数据
	trailer := metadata.Pairs(
		"processing-time", time.Since(startTime).String(),
		"new-user-id", userID,
		"operation-result", "success",
	)
	grpc.SetTrailer(ctx, trailer)

	return &rpc.CreateUserResponse{
		UserId:  userID,
		Message: "用户创建成功",
	}, nil
}

// Login 用户登录
func (s *UserServer) Login(ctx context.Context, req *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	startTime := time.Now()

	// 提取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "缺少元数据")
	}

	printRequestMetadata(md)

	// 获取客户端信息
	userAgent := getMetadataValue(md, "user-agent")
	clientIP := getMetadataValue(md, "x-client-ip")
	deviceID := getMetadataValue(md, "x-device-id")

	log.Printf("处理Login请求 - Username: %s, ClientIP: %s, DeviceID: %s, UserAgent: %s",
		req.GetUsername(), clientIP, deviceID, userAgent)

	// 发送头部元数据
	header := metadata.Pairs(
		"server-version", "1.0.0",
		"auth-method", "password",
	)
	grpc.SendHeader(ctx, header)

	// 模拟登录验证
	time.Sleep(150 * time.Millisecond)

	// 简单的用户名密码验证
	if req.GetUsername() == "admin" && req.GetPassword() == "123456" {
		// 登录成功
		token := fmt.Sprintf("token_%d", time.Now().Unix())
		userID := "user_001"

		// 设置成功的尾部元数据
		trailer := metadata.Pairs(
			"processing-time", time.Since(startTime).String(),
			"login-result", "success",
			"session-created", "true",
		)
		grpc.SetTrailer(ctx, trailer)

		return &rpc.LoginResponse{
			Token:   token,
			UserId:  userID,
			Message: "登录成功",
		}, nil
	} else {
		// 登录失败
		trailer := metadata.Pairs(
			"processing-time", time.Since(startTime).String(),
			"login-result", "failed",
			"error-reason", "invalid_credentials",
		)
		grpc.SetTrailer(ctx, trailer)

		return nil, status.Error(codes.Unauthenticated, "用户名或密码错误")
	}
}

func main() {
	// 监听端口
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}

	// 创建gRPC服务器
	server := grpc.NewServer()

	// 注册用户服务
	rpc.RegisterUserServiceServer(server, &UserServer{})

	log.Println("gRPC服务器启动成功，监听端口 :8080")
	log.Println("等待客户端连接...")

	// 启动服务器
	if err := server.Serve(lis); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
