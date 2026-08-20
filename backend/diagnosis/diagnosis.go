package main

import (
	"flag"
	"fmt"

	"faceTest/backend/diagnosis/internal/config"
	"faceTest/backend/diagnosis/internal/consumer"
	"faceTest/backend/diagnosis/internal/server"
	"faceTest/backend/diagnosis/internal/svc"
	diagnosispb "faceTest/backend/proto/diagnosis"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "diagnosis/etc/diagnosis.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		diagnosispb.RegisterDiagnosisServer(grpcServer, server.NewDiagnosisServer(ctx))

		if c.Mode == service.DevMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 启动诊断消息消费者（Outbox中继 + RabbitMQ + 死信队列）
	diagConsumer := consumer.NewDiagnosisConsumer(ctx)
	diagConsumer.Start()
	defer diagConsumer.Stop()

	s.AddOptions(grpc.MaxRecvMsgSize(10 * 1024 * 1024))

	fmt.Printf("诊断服务启动中...\n")
	fmt.Printf("监听地址: %s\n", c.ListenOn)
	fmt.Printf("RabbitMQ队列: %s\n", c.RabbitMQ.Queue)
	fmt.Printf("DeepSeek模型: %s\n", c.DeepSeek.Model)

	s.Start()
}
