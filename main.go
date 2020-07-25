package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/daigd/microservice-go/config"
	"github.com/daigd/microservice-go/discover"
	"github.com/daigd/microservice-go/endpoint"
	"github.com/daigd/microservice-go/service"
	"github.com/daigd/microservice-go/transport"
	uuid "github.com/satori/go.uuid"
)

func main() {
	fmt.Println("Hello,微服务!")
	// 从命令行读取相关参数，没有则使用默认值，默认的服务地址：127.0.0.1：10086，默认的接口服务为：SayHello
	var (
		serviceHost = flag.String("servcie.host", "127.0.0.1", "service host")
		servicePort = flag.Int("service.port", 10086, "service port")
		serviceName = flag.String("service.name", "SayHello", "service name")
	)
	// 解析命令行参数
	flag.Parse()

	// 初始化上下文环境，Backgroud() 返回一个 空集上下文环境，常用于对上下文的初始化
	ctx := context.Background()

	// 声明服务发现客户端
	var discoverClient discover.DiscoverClient
	//TODO 未实现 DiscoverClient
	// 初始化 Service
	svcImpl := service.NewDiscoverServiceImpl(discoverClient)

	// 创建 SayHello Endpoint
	sayHelloEnpoint := endpoint.MakeSayHelloEndpoint(svcImpl)

	// 将所有的服务 endpoint 拼装起来
	endpts := endpoint.DiscoverEndpoint{
		SayHelloEndpoint: sayHelloEnpoint,
	}
	// 创建请求处理器 http.Handler
	handler := transport.MakeHTTPHandler(ctx, endpts, config.KitLogger)
	// 定义服务实例ID
	servicerID := *serviceName + ":" + uuid.NewV4().String()

	errChan := make(chan error)
	go func() {
		config.Logger.Printf("Server start at port:%s\n", strconv.Itoa(*servicePort))
		// 服务启动前先注册
		// 服务名默认是 SayHello,服务ID 由 服务名 + UUID 组成，声明健康检查的地址为 /health
		if !discoverClient.Register(*serviceName, servicerID, "/health", *serviceHost, *servicePort, nil, config.Logger) {
			config.Logger.Printf("Service %s register faild!\n", *serviceName)
			// 服务注册失败，程序退出
			os.Exit(-1)
		}
		h := handler
		// 按指定的地址启动服务并监听请求
		errChan <- http.ListenAndServe(*serviceHost+":"+strconv.Itoa(*servicePort), h)
	}()

	// 监控系统信号
	go func() {
		c := make(chan os.Signal, 1)
		// 监听的信号：interrupt，killed，terminated
		signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
		errChan <- fmt.Errorf("系统被迫退出:%s", <-c)
	}()

	err := <-errChan
	// 发现错误，注销服务
	discoverClient.DeRegister(servicerID, config.Logger)
	config.Logger.Println(err)
}
