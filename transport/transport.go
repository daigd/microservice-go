// Package transport 指定向外暴露服务的提供方式，如本项目中指定为 Http 协议
package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"

	"github.com/daigd/microservice-go/endpoint"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHTTPHandler 构建请求处理器
func MakeHTTPHandler(ctx context.Context, endpoint endpoint.DiscoverEndpoint, logger log.Logger) http.Handler {
	// 创建一个服务路由器
	r := mux.NewRouter()

	ops := []kithttp.ServerOption{
		// 服务异常处理器函数
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		// 服务异常编码处理器函数，异常编码器使用自定义的 errorEncoder
		kithttp.ServerErrorEncoder(errorEncoder),
	}
	// 声明并初始化 GET 请求的 say-hello 接口，初始化需要用到对应的 SayHelloEndpoint, Request请求解码器，Response响应编码器
	r.Methods("GET").Path("/say-hello").Handler(kithttp.NewServer(
		endpoint.SayHelloEndpoint,
		decodeSayHelloRequest,
		encodeJSONResponse,
		ops...,
	))
	// 声明并初始化 /health 请求
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoint.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJSONResponse,
		ops...))
	return r
}

// 对say-hello 请求参数进行转换（解码）
func decodeSayHelloRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return &endpoint.SayHelloRequest{}, nil
}

// 对 health 请求参数进行转换
func decodeHealthCheckRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return &endpoint.HealthCheckRequest{}, nil
}

// 将接口返回结果编码成 json 格式 返回
func encodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(response)
	return nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	// 对服务异常信息指定返回格式，将异常码指定成 500
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
