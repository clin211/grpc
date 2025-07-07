package trace

import (
	"fmt"
	"log"
	"time"
)

// TraceLogger 追踪日志记录器
type TraceLogger struct {
	serviceName string
}

// NewTraceLogger 创建追踪日志记录器
func NewTraceLogger(serviceName string) *TraceLogger {
	return &TraceLogger{serviceName: serviceName}
}

// LogRequest 记录请求开始
func (tl *TraceLogger) LogRequest(traceInfo *TraceInfo, method, details string) {
	log.Printf("[%s] 请求开始 - TraceID: %s, SpanID: %s, Method: %s, Details: %s",
		tl.serviceName, traceInfo.TraceID, traceInfo.SpanID, method, details)
}

// LogResponse 记录请求结束
func (tl *TraceLogger) LogResponse(traceInfo *TraceInfo, method, data string, duration time.Duration, err error) {
	status := "SUCCESS"
	if err != nil {
		status = fmt.Sprintf("ERROR: %v", err)
	}

	log.Printf("[%s] 请求结束 - TraceID: %s, SpanID: %s, Method: %s, Data: %s, Duration: %v, Status: %s",
		tl.serviceName, traceInfo.TraceID, traceInfo.SpanID, method, data, duration, status)
}

// LogDownstreamCall 记录下游服务调用
func (tl *TraceLogger) LogDownstreamCall(traceInfo *TraceInfo, targetService, method string) {
	log.Printf("[%s] 调用下游服务 - TraceID: %s, SpanID: %s, Target: %s, Method: %s",
		tl.serviceName, traceInfo.TraceID, traceInfo.SpanID, targetService, method)
}
