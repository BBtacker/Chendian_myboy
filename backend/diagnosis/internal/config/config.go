package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 诊断服务配置
type Config struct {
	zrpc.RpcServerConf

	MySQL struct {
		DataSource string
	}

	// 注意: 不能用 Redis 命名，RpcServerConf 内嵌字段已占用
	RedisStore struct {
		Host     string
		Password string
		DB       int
	}

	RabbitMQ struct {
		URL            string
		Exchange       string `json:",default=diagnosis.exchange"`
		Queue          string `json:",default=diagnosis.queue"`
		DLXExchange    string `json:",default=diagnosis.dlx.exchange"`
		DLXQueue       string `json:",default=diagnosis.dlx.queue"`
		PrefetchCount  int    `json:",default=1"`
	}

	Milvus struct {
		Address string
	}

	DeepSeek struct {
		APIKey string
		Model  string `json:",default=deepseek-chat"`
	}

	// 特征提取服务（EfficientNet-B3 ONNX, Python 微服务）
	FeatureExtractor struct {
		URL     string `json:",default=http://127.0.0.1:8085"`
		Timeout int    `json:",default=60"`    // 秒
		Enabled bool   `json:",default=false"` // true 时使用真实特征提取，false 回退 Mock
	}

	// 上传图片存储根目录（相对诊断服务工作目录）
	// 用于读取 task.ImagePath 对应的图片文件
	UploadBasePath string `json:",default=./uploads"`

	Cache struct {
		TTL int `json:",default=3600"` // 缓存过期时间（秒）
	}

	Diagnosis struct {
		MaxRetry       int `json:",default=5"`
		RetryBaseDelay int `json:",default=1"` // 重试基础延迟（秒）
	}
}
