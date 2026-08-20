package main

import (
	"flag"
	"fmt"

	"faceTest/backend/report/internal/config"
	"faceTest/backend/report/internal/server"
	"faceTest/backend/report/internal/svc"
	reportpb "faceTest/backend/proto/report"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "report/etc/report.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		reportpb.RegisterReportServer(grpcServer, server.NewReportServer(ctx))

		if c.Mode == service.DevMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	s.AddOptions(grpc.MaxRecvMsgSize(50 * 1024 * 1024), grpc.MaxSendMsgSize(50 * 1024 * 1024))

	fmt.Printf("报告生成服务启动中...\n")
	fmt.Printf("监听地址: %s\n", c.ListenOn)

	s.Start()
}
