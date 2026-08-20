package main

import (
	"flag"
	"fmt"

	"faceTest/backend/upload/internal/config"
	"faceTest/backend/upload/internal/server"
	"faceTest/backend/upload/internal/svc"
	uploadpb "faceTest/backend/proto/upload"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "upload/etc/upload.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		uploadpb.RegisterUploadServer(grpcServer, server.NewUploadServer(ctx))

		if c.Mode == service.DevMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	s.AddOptions(grpc.MaxRecvMsgSize(20 * 1024 * 1024)) // 20MB for image uploads

	fmt.Printf("图像上传服务启动中...\n")
	fmt.Printf("监听地址: %s\n", c.ListenOn)

	s.Start()
}
