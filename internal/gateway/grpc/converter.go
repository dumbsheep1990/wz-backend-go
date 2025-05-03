package grpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
)

// HTTPToGRPC 处理HTTP到gRPC的转换
func HTTPToGRPC(c *gin.Context, client *Client, methodName string, timeout int) {
	// 获取请求体
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无法读取请求体",
			"error":   err.Error(),
		})
		return
	}

	// 构建gRPC方法全路径
	fullMethod := fmt.Sprintf("/%s/%s", client.ServiceName, methodName)

	// 创建元数据上下文
	ctx := createMetadataContext(c, timeout)

	// 创建消息缓冲区
	var resp []byte

	// 准备发送的数据（如果是GET请求或者空请求体，创建空的JSON对象）
	if len(body) == 0 || c.Request.Method == http.MethodGet {
		body = []byte("{}")
	}

	// 调用gRPC方法
	respHeaders := metadata.MD{}
	respTrailers := metadata.MD{}

	err = client.Conn.Invoke(ctx, fullMethod, body, &resp, 
		grpc.Header(&respHeaders), grpc.Trailer(&respTrailers))

	// 处理响应
	handleGRPCResponse(c, respHeaders, resp, err)
}

// createMetadataContext 创建带元数据的上下文
func createMetadataContext(c *gin.Context, timeoutSeconds int) context.Context {
	// 创建上下文
	ctx := c.Request.Context()

	// 创建元数据
	md := metadata.New(nil)

	// 从HTTP头传递关键信息到gRPC元数据
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			// 转换为小写，因为gRPC元数据键是大小写敏感的
			key := strings.ToLower(k)
			md.Set(key, v...)
		}
	}

	// 传递查询参数
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			key := fmt.Sprintf("query-%s", strings.ToLower(k))
			md.Set(key, v...)
		}
	}

	// 传递路径参数
	for _, param := range c.Params {
		key := fmt.Sprintf("param-%s", strings.ToLower(param.Key))
		md.Set(key, param.Value)
	}

	// 特别处理常用的元数据
	if userID, exists := c.Get("userID"); exists {
		md.Set("x-user-id", fmt.Sprintf("%v", userID))
	}
	if tenantID, exists := c.Get("tenantID"); exists {
		md.Set("x-tenant-id", fmt.Sprintf("%v", tenantID))
	}
	if role, exists := c.Get("role"); exists {
		md.Set("x-user-role", fmt.Sprintf("%v", role))
	}

	// 添加请求ID用于跟踪
	if requestID, exists := c.Get("requestID"); exists {
		md.Set("x-request-id", fmt.Sprintf("%v", requestID))
	}

	// 创建带元数据的上下文
	return metadata.NewOutgoingContext(ctx, md)
}

// handleGRPCResponse 处理gRPC响应
func handleGRPCResponse(c *gin.Context, headers metadata.MD, respBytes []byte, err error) {
	// 处理错误情况
	if err != nil {
		// 从gRPC错误中提取状态
		st := status.Convert(err)
		httpStatus := grpcCodeToHTTPStatus(st.Code())
		
		// 构建错误响应
		errorResp := gin.H{
			"code":    int(st.Code()),
			"message": st.Message(),
		}

		// 添加详情信息
		if len(st.Details()) > 0 {
			details := make([]interface{}, 0, len(st.Details()))
			for _, detail := range st.Details() {
				if msg, ok := detail.(proto.Message); ok {
					jsonBytes, err := protojson.Marshal(msg)
					if err == nil {
						var jsonObj map[string]interface{}
						if err := json.Unmarshal(jsonBytes, &jsonObj); err == nil {
							details = append(details, jsonObj)
						}
					}
				}
			}
			if len(details) > 0 {
				errorResp["details"] = details
			}
		}

		c.JSON(httpStatus, errorResp)
		return
	}

	// 处理成功响应
	
	// 将gRPC响应头传递到HTTP头
	for k, v := range headers {
		for _, value := range v {
			c.Header(k, value)
		}
	}

	// 如果不是有效的JSON，进行转换
	if !isValidJSON(respBytes) {
		// 尝试将Protobuf响应转换为JSON
		var buf bytes.Buffer
		buf.WriteString("{\"data\":\"")
		buf.Write(respBytes)
		buf.WriteString("\"}")
		respBytes = buf.Bytes()
	}

	// 设置内容类型并返回响应
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", respBytes)
}

// isValidJSON 检查字节切片是否是有效的JSON
func isValidJSON(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}

// grpcCodeToHTTPStatus 将gRPC状态码转换为HTTP状态码
func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
