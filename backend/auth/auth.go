package main

import (
	"flag"
	"fmt"

	"faceTest/backend/auth/internal/config"
	"faceTest/backend/auth/internal/server"
	"faceTest/backend/auth/internal/svc"
	authpb "faceTest/backend/proto/auth"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "auth/etc/auth.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 创建gRPC服务器
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		authpb.RegisterAuthServer(grpcServer, server.NewAuthServer(ctx))

		// 开发模式下注册反射服务
		if c.Mode == service.DevMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 添加服务启动钩子
	s.AddOptions(grpc.MaxRecvMsgSize(10 * 1024 * 1024)) // 10MB

	fmt.Printf("认证服务启动中...\n")
	fmt.Printf("监听地址: %s\n", c.ListenOn)
	fmt.Printf("Etcd注册: %v\n", len(c.Etcd.Hosts) > 0)

	s.Start()
}
