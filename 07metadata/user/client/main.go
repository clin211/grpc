package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	rpc "github.com/clin211/grpc/metadata/proto"
)

// generateTraceID 生成追踪ID
func generateTraceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().Unix())
}

// 演示1：基本元数据发送
func demonstrateBasicMetadata(client rpc.UserServiceClient) {
	fmt.Println("\n========== 演示1：基本元数据发送 ==========")

	// 创建基本元数据
	md := metadata.Pairs(
		"authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.sample.token",
		"user-agent", "grpc-client/1.0.0",
		"client-version", "1.2.0",
		"x-trace-id", generateTraceID(),
	)

	// 将元数据附加到context
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 发起RPC调用
	resp, err := client.GetUser(ctx, &rpc.GetUserRequest{
		UserId: "user_123",
	})

	if err != nil {
		log.Printf("调用失败: %v", err)
		return
	}

	fmt.Printf("调用成功: %v\n", resp)
}

// 演示2：动态添加元数据
func demonstrateDynamicMetadata(client rpc.UserServiceClient) {
	fmt.Println("\n========== 演示2：动态添加元数据 ==========")

	ctx := context.Background()

	// 第一步：创建基础元数据
	baseMD := metadata.Pairs(
		"authorization", "Bearer dynamic.token.example",
		"user-agent", "grpc-client/1.0.0",
	)
	ctx = metadata.NewOutgoingContext(ctx, baseMD)

	// 第二步：动态添加更多元数据
	md, _ := metadata.FromOutgoingContext(ctx)
	md.Set("x-request-id", generateRequestID())
	md.Append("x-custom-header", "value1", "value2", "value3")
	md.Set("x-client-ip", "192.168.1.100")

	// 添加权限信息
	md.Append("x-permission", "create", "read", "admin")

	// 更新context
	ctx = metadata.NewOutgoingContext(ctx, md)

	// 发起创建用户的调用
	resp, err := client.CreateUser(ctx, &rpc.CreateUserRequest{
		Username: "新用户",
		Email:    "newuser@example.com",
		Password: "password123",
	})

	if err != nil {
		log.Printf("调用失败: %v", err)
		return
	}

	fmt.Printf("调用成功: %v\n", resp)
}

// 演示3：接收响应元数据
func demonstrateReceiveMetadata(client rpc.UserServiceClient) {
	fmt.Println("\n========== 演示3：接收响应元数据 ==========")

	// 准备接收元数据的变量
	var header, trailer metadata.MD

	// 创建请求元数据
	md := metadata.Pairs(
		"authorization", "Bearer receive.metadata.token",
		"user-agent", "grpc-client/1.0.0",
		"x-trace-id", generateTraceID(),
		"x-device-id", "device_12345",
		"x-client-ip", "203.0.113.1",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 发起调用，同时指定接收元数据
	resp, err := client.Login(
		ctx,
		&rpc.LoginRequest{
			Username: "admin",
			Password: "123456",
		},
		grpc.Header(&header),   // 接收头部元数据
		grpc.Trailer(&trailer), // 接收尾部元数据
	)

	if err != nil {
		fmt.Printf("调用失败: %v\n", err)

		// 即使调用失败，也可能收到尾部元数据
		if len(trailer) > 0 {
			fmt.Println("收到的尾部元数据（错误情况）:")
			for key, values := range trailer {
				fmt.Printf("  %s: %v\n", key, values)
			}
		}
		return
	}

	// 处理头部元数据
	fmt.Println("=== 收到的头部元数据 ===")
	for key, values := range header {
		fmt.Printf("  %s: %v\n", key, values)
	}

	// 处理响应数据
	fmt.Printf("登录响应: %v\n", resp)

	// 处理尾部元数据
	fmt.Println("=== 收到的尾部元数据 ===")
	for key, values := range trailer {
		fmt.Printf("  %s: %v\n", key, values)
	}
}

// 演示4：错误处理和元数据
func demonstrateErrorWithMetadata(client rpc.UserServiceClient) {
	fmt.Println("\n========== 演示4：错误处理和元数据 ==========")

	var header, trailer metadata.MD

	// 故意发送错误的认证信息
	md := metadata.Pairs(
		"authorization", "Invalid token format", // 错误的格式
		"user-agent", "grpc-client/1.0.0",
		"x-trace-id", generateTraceID(),
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := client.GetUser(
		ctx,
		&rpc.GetUserRequest{UserId: "user_456"},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)

	if err != nil {
		// 解析gRPC状态错误
		if st, ok := status.FromError(err); ok {
			fmt.Printf("错误状态码: %s\n", st.Code())
			fmt.Printf("错误消息: %s\n", st.Message())
		}

		// 检查是否有尾部元数据
		if len(trailer) > 0 {
			fmt.Println("收到的尾部元数据（错误情况）:")
			for key, values := range trailer {
				fmt.Printf("  %s: %v\n", key, values)
			}
		}
	}
}

// 演示5：多值元数据
func demonstrateMultiValueMetadata(client rpc.UserServiceClient) {
	fmt.Println("\n========== 演示5：多值元数据 ==========")

	// 创建多值元数据
	md := metadata.New(map[string]string{
		"authorization": "Bearer multivalue.token.example",
		"user-agent":    "grpc-client/1.0.0",
		"x-trace-id":    generateTraceID(),
	})

	// 添加多个相同键的值
	md.Append("x-supported-format", "json", "protobuf", "xml")
	md.Append("x-feature-flag", "new-ui", "advanced-search", "real-time-updates")
	md.Append("x-user-role", "user", "moderator")

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 发起调用
	resp, err := client.GetUser(ctx, &rpc.GetUserRequest{
		UserId: "user_789",
	})

	if err != nil {
		log.Printf("调用失败: %v", err)
		return
	}

	fmt.Printf("调用成功: %v\n", resp)
}

// 演示6：二进制元数据
func demonstrateBinaryMetadata(client rpc.UserServiceClient) {
	fmt.Println("\n========== 演示6：二进制元数据 ==========")

	// 创建文本元数据
	md := metadata.Pairs(
		"authorization", "Bearer binary.metadata.token",
		"user-agent", "grpc-client/1.0.0",
		"x-trace-id", generateTraceID(),
	)

	// 添加二进制数据（键名必须以 -bin 结尾）
	binaryData := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x47, 0x52, 0x50, 0x43} // "Hello gRPC"
	signature := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	md.Set("custom-data-bin", string(binaryData))
	md.Set("signature-bin", string(signature))

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 发起调用
	resp, err := client.GetUser(ctx, &rpc.GetUserRequest{
		UserId: "user_binary",
	})

	if err != nil {
		log.Printf("调用失败: %v", err)
		return
	}

	fmt.Printf("调用成功: %v\n", resp)
}

// 演示7：链路追踪信息传递
func demonstrateTracingMetadata(client rpc.UserServiceClient) {
	fmt.Println("\n========== 演示7：链路追踪信息传递 ==========")

	// 模拟上游服务传递的追踪信息
	traceID := generateTraceID()
	spanID := generateTraceID()
	parentSpanID := generateTraceID()

	md := metadata.Pairs(
		"authorization", "Bearer tracing.token.example",
		"user-agent", "grpc-client/1.0.0",
		// 追踪相关的元数据
		"x-trace-id", traceID,
		"x-span-id", spanID,
		"x-parent-span-id", parentSpanID,
		"x-request-id", generateRequestID(),
		// 其他上下文信息
		"x-user-id", "user_12345",
		"x-session-id", "session_abcdef",
		"x-correlation-id", fmt.Sprintf("corr_%d", time.Now().Unix()),
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	fmt.Printf("发起请求 - TraceID: %s, SpanID: %s\n", traceID, spanID)

	// 发起调用
	resp, err := client.GetUser(ctx, &rpc.GetUserRequest{
		UserId: "user_tracing",
	})

	if err != nil {
		log.Printf("请求失败 [%s]: %v", traceID, err)
		return
	}

	fmt.Printf("请求成功 [%s]: %v\n", traceID, resp)
}

func main() {
	// 连接到gRPC服务器
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接服务器失败: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := rpc.NewUserServiceClient(conn)

	fmt.Println("gRPC元数据客户端示例启动")
	fmt.Println("连接到服务器: localhost:8080")

	// 依次演示各种元数据发送方式
	demonstrateBasicMetadata(client)
	time.Sleep(1 * time.Second)

	demonstrateDynamicMetadata(client)
	time.Sleep(1 * time.Second)

	demonstrateReceiveMetadata(client)
	time.Sleep(1 * time.Second)

	demonstrateErrorWithMetadata(client)
	time.Sleep(1 * time.Second)

	demonstrateMultiValueMetadata(client)
	time.Sleep(1 * time.Second)

	demonstrateBinaryMetadata(client)
	time.Sleep(1 * time.Second)

	demonstrateTracingMetadata(client)

	fmt.Println("\n所有元数据演示完成!")
}
