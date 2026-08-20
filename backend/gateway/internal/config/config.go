package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config API网关配置
type Config struct {
	rest.RestConf

	AuthRpc      zrpc.RpcClientConf
	UploadRpc    zrpc.RpcClientConf
	DiagnosisRpc zrpc.RpcClientConf
	ReportRpc    zrpc.RpcClientConf

	MySQL struct {
		DataSource string
	}

	Upload struct {
		BasePath string
		BaseURL  string
	}

	DeepSeek struct {
		APIKey string
		Model  string `json:",default=deepseek-chat"`
	}

	JWT struct {
		Secret string // 与 auth 服务保持一致，用于校验请求 Token
	}

	Redis struct {
		Host     string
		Password string
		DB       int
	}
}
