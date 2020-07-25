// Package service 处理具体的业务逻辑
package service

import (
	"context"
	"errors"

	"github.com/daigd/microservice-go/config"
	"github.com/daigd/microservice-go/discover"
)

// ErrServiceNotExisted 服务不存在错误
var ErrServiceNotExisted = errors.New("Service not existed")

// Service 提供的服务接口
type Service interface {
	// 健康检查接口
	HealthCheck() bool
	// 打招呼接口
	SayHello() string
	// 服务发现接口
	DiscoverService(ctx context.Context, serviceName string) ([]interface{}, error)
}

// DiscoverSerivceImpl 服务实现
type DiscoverSerivceImpl struct {
	discoveClient discover.DiscoverClient
}

// NewDiscoverServiceImpl 服务实现
func NewDiscoverServiceImpl(discoverClient discover.DiscoverClient) Service {
	return &DiscoverSerivceImpl{discoveClient: discoverClient}
}

// HealthCheck 健康检查
func (service *DiscoverSerivceImpl) HealthCheck() bool {
	return true
}

// SayHello 打招呼
func (service *DiscoverSerivceImpl) SayHello() string {
	return "Hello I am a service"
}

// DiscoverService 服务发现
func (service *DiscoverSerivceImpl) DiscoverService(ctx context.Context, serviceName string) ([]interface{}, error) {
	instances := service.discoveClient.DiscoverService(serviceName, config.Logger)
	if instances == nil || len(instances) == 0 {
		return nil, ErrServiceNotExisted
	}
	return instances, nil
}
