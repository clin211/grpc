package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	rpc "github.com/clin211/grpc/metadata/trace/proto"
	"github.com/clin211/grpc/metadata/trace/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// extractTraceInfo 从元数据中提取追踪信息
func extractTraceInfo(md metadata.MD) *trace.TraceInfo {
	traceID := getFirstValue(md, trace.HeaderTraceID)
	spanID := getFirstValue(md, trace.HeaderSpanID)
	parentSpanID := getFirstValue(md, trace.HeaderParentSpanID)

	if traceID == "" {
		// 如果没有追踪信息，创建新的
		return trace.NewTraceInfo()
	}

	return &trace.TraceInfo{
		TraceID:      traceID,
		SpanID:       spanID,
		ParentSpanID: parentSpanID,
	}
}

// UserServer 用户服务实现
type UserServer struct {
	rpc.UnimplementedProfileServiceServer
}

// GetProfile 获取用户资料（支持链路追踪）
func (s *UserServer) GetProfile(ctx context.Context, req *rpc.GetProfileRequest) (*rpc.GetProfileResponse, error) {
	// 业务逻辑处理的开始时间
	startTime := time.Now()
	// 提取追踪信息
	md, _ := metadata.FromIncomingContext(ctx)
	traceInfo := extractTraceInfo(md)

	// 创建追踪日志记录器
	tracer := trace.NewTraceLogger("UserService")
	// 记录请求开始
	tracer.LogRequest(traceInfo, "GetUser", fmt.Sprintf("UserID: %s", req.GetUserId()))

	log.Printf("[追踪] 收到GetProfile请求 - TraceID: %s, SpanID: %s, UserID: %s",
		traceInfo.TraceID, traceInfo.SpanID, req.GetUserId())

	// 模拟业务逻辑：获取基本用户信息
	userInfo := &rpc.GetProfileResponse{
		Profile: &rpc.ProfileInfo{
			UserId:    req.GetUserId(),
			Nickname:  "张三",
			AvatarUrl: "https://example.com/avatar.jpg",
			Bio:       "这是一个示例用户",
			Location:  "北京",
			Website:   "https://example.com",
			Interests: []string{"编程", "阅读", "旅行"},
		},
	}

	tracer.LogResponse(traceInfo, "Get Profile", fmt.Sprintf("userInfo: %v", userInfo), time.Since(startTime), nil)

	log.Printf("[追踪] GetProfile处理完成 - TraceID: %s, SpanID: %s",
		traceInfo.TraceID, traceInfo.SpanID)

	return userInfo, nil
}

// getFirstValue 获取元数据中的第一个值
func getFirstValue(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	profileService := &UserServer{}
	rpc.RegisterProfileServiceServer(grpcServer, profileService)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
