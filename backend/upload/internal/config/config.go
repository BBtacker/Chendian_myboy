package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 图像上传服务配置
type Config struct {
	zrpc.RpcServerConf

	Upload struct {
		BasePath    string
		BaseURL     string
		MaxFileSize int64 `json:",default=10485760"` // 10MB
	}
}
