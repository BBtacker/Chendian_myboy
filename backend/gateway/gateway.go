package main

import (
	"flag"
	"fmt"

	"faceTest/backend/gateway/internal/config"
	"faceTest/backend/gateway/internal/handler"
	"faceTest/backend/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "gateway/etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建网关服务上下文
	ctx := svc.NewServiceContext(c)

	// 创建HTTP服务器
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	// 注册路由
	handler.RegisterRoutes(server, ctx)

	fmt.Printf("API网关启动中...\n")
	fmt.Printf("监听地址: http://%s:%d\n", c.Host, c.Port)
	fmt.Printf("CORS: 已启用\n")

	server.Start()
}
