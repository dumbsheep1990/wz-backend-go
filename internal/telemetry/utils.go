package telemetry

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SetSpanStatus 设置Span的状态
func SetSpanStatus(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "")
	}
}

// SetSpanAttributes 设置Span的属性
func SetSpanAttributes(span trace.Span, attrs map[string]string) {
	if span == nil || len(attrs) == 0 {
		return
	}

	var attributes []attribute.KeyValue
	for k, v := range attrs {
		attributes = append(attributes, attribute.String(k, v))
	}

	span.SetAttributes(attributes...)
}

// RecordException 记录异常到Span
func RecordException(span trace.Span, err error, attrs map[string]string) {
	if span == nil || err == nil {
		return
	}

	var attributes []attribute.KeyValue
	for k, v := range attrs {
		attributes = append(attributes, attribute.String(k, v))
	}

	span.RecordError(err, trace.WithAttributes(attributes...))
	span.SetStatus(codes.Error, err.Error())
}

// AddUserIDToSpan 添加用户ID到Span
func AddUserIDToSpan(span trace.Span, userID string) {
	if span == nil || userID == "" {
		return
	}

	span.SetAttributes(attribute.String("user.id", userID))
}

// AddTenantIDToSpan 添加租户ID到Span
func AddTenantIDToSpan(span trace.Span, tenantID string) {
	if span == nil || tenantID == "" {
		return
	}

	span.SetAttributes(attribute.String("tenant.id", tenantID))
}

// AddRequestIDToSpan 添加请求ID到Span
func AddRequestIDToSpan(span trace.Span, requestID string) {
	if span == nil || requestID == "" {
		return
	}

	span.SetAttributes(attribute.String("request.id", requestID))
}

// AddHTTPRequestAttributesToSpan 添加HTTP请求属性到Span
func AddHTTPRequestAttributesToSpan(span trace.Span, r *http.Request) {
	if span == nil || r == nil {
		return
	}

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.host", r.Host),
		attribute.String("http.user_agent", r.UserAgent()),
		attribute.String("http.remote_addr", r.RemoteAddr),
	)
}

// WithTraceContext 创建一个带有追踪上下文的函数包装器
func WithTraceContext(operation string, f func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx, span := StartSpan(ctx, operation)
		defer span.End()

		// 添加调用者函数信息
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			span.SetAttributes(
				attribute.String("code.function", funcName),
				attribute.String("code.filepath", file),
				attribute.Int("code.lineno", line),
			)
		}

		err := f(ctx)
		SetSpanStatus(span, err)
		return err
	}
}

// TraceFunctionWithName 跟踪函数调用并自动设置函数名称为Span名称
func TraceFunctionWithName(ctx context.Context) (context.Context, trace.Span) {
	pc, _, _, ok := runtime.Caller(1)
	fnName := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			fnName = fn.Name()
			// 提取短函数名称（去除包路径）
			if idx := strings.LastIndex(fnName, "."); idx >= 0 {
				fnName = fnName[idx+1:]
			}
		}
	}
	
	return StartSpan(ctx, fnName)
}

// TraceFunctionCall 跟踪函数调用，并自动添加函数信息和错误处理
func TraceFunctionCall(ctx context.Context, span *trace.Span) func(error) error {
	pc, file, line, ok := runtime.Caller(1)
	if ok && *span != nil {
		fnName := runtime.FuncForPC(pc).Name()
		(*span).SetAttributes(
			attribute.String("code.function", fnName),
			attribute.String("code.filepath", file),
			attribute.Int("code.lineno", line),
		)
	}

	return func(err error) error {
		if *span != nil {
			SetSpanStatus(*span, err)
		}
		return err
	}
}

// PrintTraceInfo 打印当前上下文的追踪信息（用于调试）
func PrintTraceInfo(ctx context.Context) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		log.Println("无有效的追踪上下文")
		return
	}

	log.Printf("追踪ID: %s, Span ID: %s, 是否采样: %t",
		spanCtx.TraceID().String(),
		spanCtx.SpanID().String(),
		spanCtx.IsSampled())
}

// FormatTraceURL 根据追踪ID和配置的后端系统生成可访问的链接
func FormatTraceURL(traceID string, backendType string, backendURL string) string {
	if traceID == "" {
		return ""
	}

	switch backendType {
	case "jaeger":
		return fmt.Sprintf("%s/trace/%s", strings.TrimSuffix(backendURL, "/"), traceID)
	case "zipkin":
		return fmt.Sprintf("%s/zipkin/traces/%s", strings.TrimSuffix(backendURL, "/"), traceID)
	}

	return ""
}
