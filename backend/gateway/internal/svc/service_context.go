package svc

import (
	"faceTest/backend/common/middleware"
	"faceTest/backend/common/pkg"
	"faceTest/backend/gateway/internal/config"
	authpb "faceTest/backend/proto/auth"
	diagnosispb "faceTest/backend/proto/diagnosis"
	reportpb "faceTest/backend/proto/report"
	uploadpb "faceTest/backend/proto/upload"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ServiceContext 网关服务上下文
type ServiceContext struct {
	Config          config.Config
	AuthMiddleware  *middleware.AuthMiddleware
	AuthClient      authpb.AuthClient
	UploadClient    uploadpb.UploadClient
	DiagnosisClient diagnosispb.DiagnosisClient
	ReportClient    reportpb.ReportClient
	RedisClient     *redis.Client
	RedisCache      *pkg.RedisCache
	RateLimiter     *pkg.RateLimiter
	DeepSeekClient  *pkg.DeepSeekClient
	FileStorage     *pkg.FileStorage
	DB              *gorm.DB
}

// NewServiceContext 创建网关服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// JWT 校验密钥从配置注入（须与 auth 服务 Secret 一致；未配置时保留默认占位）
	if c.JWT.Secret != "" {
		pkg.SetJWTSecret(c.JWT.Secret)
	}

	// 创建RPC客户端（通过etcd服务发现）
	authClient := authpb.NewAuthClient(zrpc.MustNewClient(c.AuthRpc).Conn())
	uploadClient := uploadpb.NewUploadClient(zrpc.MustNewClient(c.UploadRpc).Conn())
	diagnosisClient := diagnosispb.NewDiagnosisClient(zrpc.MustNewClient(c.DiagnosisRpc).Conn())
	reportClient := reportpb.NewReportClient(zrpc.MustNewClient(c.ReportRpc).Conn())

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	// DeepSeek客户端（用于AI对话）
	var deepseekClient *pkg.DeepSeekClient
	if c.DeepSeek.APIKey != "" {
		deepseekClient = pkg.NewDeepSeekClient(c.DeepSeek.APIKey, c.DeepSeek.Model)
	}

	// MySQL（用于对话功能持久化）
	var db *gorm.DB
	if c.MySQL.DataSource != "" {
		var err error
		db, err = gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err != nil {
			logx.Errorf("网关MySQL连接失败: %v", err)
		} else {
			logx.Info("网关MySQL连接成功")
		}
	}

	logx.Info("网关服务上下文初始化完成")

	return &ServiceContext{
		Config:          c,
		AuthMiddleware:  middleware.NewAuthMiddleware(),
		AuthClient:      authClient,
		UploadClient:    uploadClient,
		DiagnosisClient: diagnosisClient,
		ReportClient:    reportClient,
		RedisClient:     redisClient,
		RedisCache:      pkg.NewRedisCache(c.Redis.Host, c.Redis.Password, c.Redis.DB),
		RateLimiter:     pkg.NewRateLimiter(redisClient),
		DeepSeekClient:  deepseekClient,
		FileStorage:     pkg.NewFileStorage(c.Upload.BasePath, c.Upload.BaseURL),
		DB:              db,
	}
}
