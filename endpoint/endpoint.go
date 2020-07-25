// Package endpoint 用于将请求参数转化为 service 可以接收的参数，并将 service 层的处理结果
// 封装成合适格式 作为响应返回
package endpoint

import (
	"context"

	"github.com/daigd/microservice-go/service"
	"github.com/go-kit/kit/endpoint"
)

// DiscoverEndpoint 服务发现 Endpoint
type DiscoverEndpoint struct {
	SayHelloEndpoint    endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

// SayHelloRequest 服务请求结构体
type SayHelloRequest struct {
}

// SayHelloResponse 服务响应结构体
type SayHelloResponse struct {
	Message string `json:"message"`
}

// MakeSayHelloEndpoint 构建打招呼 Endpoint
func MakeSayHelloEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		msg := svc.SayHello()
		return &SayHelloResponse{Message: msg}, nil
	}
}

// HealthCheckRequest 健康检查请求体
type HealthCheckRequest struct{}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 构建健康检查响应
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return &HealthCheckResponse{Status: status}, nil
	}
}
