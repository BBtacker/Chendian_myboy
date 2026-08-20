package svc

import (
	"faceTest/backend/common/model"
	"faceTest/backend/common/pkg"
	"faceTest/backend/diagnosis/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ServiceContext 诊断服务上下文
type ServiceContext struct {
	Config           config.Config
	DB               *gorm.DB
	RedisClient      *redis.Client
	RedisCache       *pkg.RedisCache
	RateLimiter      *pkg.RateLimiter
	RabbitMQ         *pkg.RabbitMQ
	MilvusClient     *pkg.MilvusClient
	DeepSeekClient   *pkg.DeepSeekClient
	FeatureExtractor *pkg.FeatureExtractor // EfficientNet-B3 特征提取（Enabled 时才初始化）
}

// NewServiceContext 创建诊断服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// MySQL
	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logx.Must(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 自动迁移诊断相关表（新增列如 skipped_images 时会自动创建，已有列不变）
	if err := db.AutoMigrate(&model.DiagnosisTask{}, &model.DiagnosisResult{}); err != nil {
		logx.Errorf("诊断表 AutoMigrate 失败: %v", err)
	}

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.RedisStore.Host,
		Password: c.RedisStore.Password,
		DB:       c.RedisStore.DB,
	})
	redisCache := pkg.NewRedisCache(c.RedisStore.Host, c.RedisStore.Password, c.RedisStore.DB)

	// RabbitMQ
	rabbitMQ, err := pkg.NewRabbitMQ(c.RabbitMQ.URL)
	if err != nil {
		logx.Errorf("RabbitMQ连接失败: %v", err)
	} else {
		// 声明带死信队列的消息队列
		err = rabbitMQ.DeclareDLXQueue(c.RabbitMQ.Queue, c.RabbitMQ.DLXExchange, c.RabbitMQ.DLXQueue)
		if err != nil {
			logx.Errorf("声明RabbitMQ队列失败: %v", err)
		}
		rabbitMQ.SetQoS(c.RabbitMQ.PrefetchCount)
	}

	// Milvus
	var milvusClient *pkg.MilvusClient
	if c.Milvus.Address != "" {
		milvusClient, err = pkg.NewMilvusClient(c.Milvus.Address)
		if err != nil {
			logx.Errorf("Milvus连接失败: %v", err)
		}
	}

	// DeepSeek
	var deepseekClient *pkg.DeepSeekClient
	if c.DeepSeek.APIKey != "" {
		deepseekClient = pkg.NewDeepSeekClient(c.DeepSeek.APIKey, c.DeepSeek.Model)
	}

	// 特征提取服务（EfficientNet-B3）
	var featureExtractor *pkg.FeatureExtractor
	if c.FeatureExtractor.Enabled {
		featureExtractor = pkg.NewFeatureExtractor(pkg.FeatureExtractorConfig{
			URL:     c.FeatureExtractor.URL,
			Timeout: c.FeatureExtractor.Timeout,
		})
		logx.Infof("真实特征提取已启用: %s", c.FeatureExtractor.URL)
	} else {
		logx.Info("特征提取未启用（使用 Mock），配置 FeatureExtractor.Enabled=true 启用真实 EfficientNet-B3")
	}

	return &ServiceContext{
		Config:           c,
		DB:               db,
		RedisClient:      redisClient,
		RedisCache:       redisCache,
		RateLimiter:      pkg.NewRateLimiter(redisClient),
		RabbitMQ:         rabbitMQ,
		MilvusClient:     milvusClient,
		DeepSeekClient:   deepseekClient,
		FeatureExtractor: featureExtractor,
	}
}
