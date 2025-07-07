package trace

import (
	"crypto/rand"
	"fmt"
)

// 标准追踪头部常量
const (
	HeaderTraceID      = "x-trace-id"       // 追踪ID
	HeaderSpanID       = "x-span-id"        // 跨度ID
	HeaderParentSpanID = "x-parent-span-id" // 父跨度ID
)

// TraceInfo 追踪信息结构
type TraceInfo struct {
	TraceID      string `json:"trace_id"`
	SpanID       string `json:"span_id"`
	ParentSpanID string `json:"parent_span_id,omitempty"`
}

// generateTraceID 生成追踪ID
func generateTraceID() string {
	// 生成16字节随机数
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// generateSpanID 生成跨度ID
func generateSpanID() string {
	// 生成8字节随机数
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// NewTraceInfo 创建新的追踪信息
func NewTraceInfo() *TraceInfo {
	return &TraceInfo{
		TraceID: generateTraceID(),
		SpanID:  generateSpanID(),
	}
}

// NewChildSpan 创建子跨度
func (t *TraceInfo) NewChildSpan() *TraceInfo {
	return &TraceInfo{
		TraceID:      t.TraceID,
		SpanID:       generateSpanID(),
		ParentSpanID: t.SpanID,
	}
}
