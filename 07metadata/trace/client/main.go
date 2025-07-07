package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	rpc "github.com/clin211/grpc/metadata/trace/proto"
	"github.com/clin211/grpc/metadata/trace/trace"
)

// TracingClient 支持链路追踪的客户端
type TracingClient struct {
	client rpc.ProfileServiceClient
}

// GetUserWithTracing 带追踪的获取用户信息
func (tc *TracingClient) GetUserWithTracing(ctx context.Context, userID string, traceInfo *trace.TraceInfo) (*rpc.GetProfileResponse, error) {
	// 将追踪信息添加到元数据
	md := metadata.Pairs(
		trace.HeaderTraceID, traceInfo.TraceID,
		trace.HeaderSpanID, traceInfo.SpanID,
	)

	if traceInfo.ParentSpanID != "" {
		md.Append(trace.HeaderParentSpanID, traceInfo.ParentSpanID)
	}

	// 创建带追踪信息的context
	tracingCtx := metadata.NewOutgoingContext(ctx, md)

	log.Printf("[追踪] 发起GetUser调用 - TraceID: %s, SpanID: %s",
		traceInfo.TraceID, traceInfo.SpanID)

	return tc.client.GetProfile(tracingCtx, &rpc.GetProfileRequest{UserId: userID})
}

// 实际使用示例
func main() {
	// 建立gRPC连接
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 创建追踪信息
	traceInfo := trace.NewTraceInfo()

	client := &TracingClient{client: rpc.NewProfileServiceClient(conn)}

	// 发起带追踪的调用
	resp, err := client.GetUserWithTracing(context.Background(), "user123", traceInfo)
	if err != nil {
		log.Printf("[追踪] 调用失败 - TraceID: %s, Error: %v", traceInfo.TraceID, err)
		return
	}

	log.Printf("[追踪] 调用成功 - TraceID: %s, Profile: %v",
		traceInfo.TraceID, resp.GetProfile())
}
