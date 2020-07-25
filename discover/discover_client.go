// Package discover 自定义服务注册与发现接口，用以屏蔽具体的服务注册与发现组件实现，方便替换组件
package discover

import "log"

// DiscoverClient 服务注册与发现客户端接口
type DiscoverClient interface {
	// Register 服务注册
	// serviceName 服务名
	// serviceID 服务ID
	// healthCheckURL 健康检查地址
	// instanceHost 服务实例地址
	// instancePort 服务实例端口
	// meta 服务实例元数据
	Register(serviceName, serviceID, healthCheckURL, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool

	// DeRegister 服务注销
	DeRegister(instanceID string, logger *log.Logger) bool

	// DiscoverService 服务发现
	DiscoverService(instanceName string, logger *log.Logger) []interface{}
}
