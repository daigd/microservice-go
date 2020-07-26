package discover

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// InstanceInfo 定义服务实例信息结构体，用以服务注册时向注册中心提供相关信息
// json 的字段 必须大写，否则 consul 无法获取相关请求数据
type InstanceInfo struct {
	ID                string            `json:"id"`                // 服务ID
	Name              string            `json:"name"`              // 服务名称
	Service           string            `json:"service,omitempty"` // 服务发现时返回的服务名
	Tags              []string          `json:"tags,omitempty"`    // 服务标签，可用于服务过滤
	Address           string            `json:"address"`           // 服务实例Host
	Port              int               `json:"port"`              // 服务实例Port
	Meta              map[string]string `json:"meta,omitempty"`    // 服务元数据
	EnableTagOverride bool              `json:"enableTagOverride"` // 是否允许标签覆盖
	Check             HealthCheckInfo   `json:"check,omitempty"`   // 健康检查信息
	Weight            WeightInfo        `json:"weights,omitempty"` // 服务实例权重
}

// HealthCheckInfo 健康检查信息
type HealthCheckInfo struct {
	DeregisterCriticalServiceAfter string   `json:"deregisterCriticalServiceAfter"` // 多久之后注销服务
	Args                           []string `json:"ags,omitempty"`                  // 请求参数
	HTTP                           string   `json:"http"`                           // 健康检查地址
	Interval                       string   `json:"interval,omitempty"`             // Consul 主动进行健康检查间隔
	TTL                            string   `json:"ttl,omitempty"`                  // 服务实例主动提交健康检查间隔，与 Interval 只存其一
}

// WeightInfo 服务实例权重
type WeightInfo struct {
	Passing int `json:"passing"`
	Warning int `json:"warning"`
}

// MyConsulDiscoverClient 自定义 Consul 客户端
type MyConsulDiscoverClient struct {
	Host string // consul Host
	Port int    // consul Port
}

// NewMyConsulDiscoverClient 创建注册中心客户端，返回自定义 Consul 客户端
func NewMyConsulDiscoverClient(host string, port int) (client DiscoverClient, err error) {
	return &MyConsulDiscoverClient{Host: host, Port: port}, nil
}

// Register 服务注册，通过自定义 Consul 客户端向 Consul 注册指定服务实例
func (client *MyConsulDiscoverClient) Register(serviceName, serviceID, healthCheckURL, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {
	// 1 封装服务实例信息
	instanceInfo := &InstanceInfo{
		ID:      serviceID,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta:    meta,
		Check: HealthCheckInfo{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckURL,
			Interval:                       "15s",
		},
		Weight: WeightInfo{
			Passing: 10,
			Warning: 1,
		},
	}
	// 2 向 Consul 发送服务注册请求
	dataBytes, _ := json.Marshal(instanceInfo)
	// 2.1 构建请求
	req, err := http.NewRequest("PUT", "http://"+client.Host+":"+strconv.Itoa(client.Port)+"/v1/agent/service/register", bytes.NewReader(dataBytes))

	if err != nil {
		logger.Println("Register Service Request Error!")
	} else {
		// 2.2 设置请求参数格式，发送请求
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		httpClient := http.Client{}
		res, err := httpClient.Do(req)

		// 3 检查请求结果
		if err != nil {
			logger.Println("Register Service Error!")
			return false
		}
		res.Body.Close()
		if res.StatusCode == 200 {
			logger.Println("Register Service Success!")
			return true
		}
		logger.Printf("Register Service Response StatusCode is:%d\n", res.StatusCode)
	}
	return false
}

// DeRegister 服务注销
func (client *MyConsulDiscoverClient) DeRegister(instanceID string, logger *log.Logger) bool {
	return false
}

// DiscoverService 服务发现
func (client *MyConsulDiscoverClient) DiscoverService(instanceName string, logger *log.Logger) []interface{} {
	return nil
}
