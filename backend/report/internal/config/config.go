package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 报告生成服务配置
type Config struct {
	zrpc.RpcServerConf

	MySQL struct {
		DataSource string
	}

	// 注意: 不能用 Redis 命名，RpcServerConf 内嵌字段已占用
	Cache struct {
		Host     string
		Password string
		DB       int
	}
}
