package svc

import (
	"faceTest/backend/auth/internal/config"
	"faceTest/backend/common/model"
	"faceTest/backend/common/pkg"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ServiceContext 认证服务上下文
type ServiceContext struct {
	Config       config.Config
	DB           *gorm.DB
	RedisClient  *redis.Client
	RedisCache   *pkg.RedisCache
	FileStorage  *pkg.FileStorage
	RateLimiter  *pkg.RateLimiter
}

// NewServiceContext 创建认证服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化MySQL
	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logx.Must(err)
	}

	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		logx.Must(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 自动迁移表结构，补齐 init.sql 与 GORM 模型不一致的列（如 user.deleted_at 软删除列）
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logx.Errorf("用户表结构迁移失败（不影响启动）: %v", err)
	}

	// 初始化Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Cache.Host,
		Password: c.Cache.Password,
		DB:       c.Cache.DB,
	})

	// 设置JWT密钥
	if c.JWT.Secret != "" {
		pkg.SetJWTSecret(c.JWT.Secret)
	}

	return &ServiceContext{
		Config:      c,
		DB:          db,
		RedisClient: redisClient,
		RedisCache:  pkg.NewRedisCache(c.Cache.Host, c.Cache.Password, c.Cache.DB),
		FileStorage: pkg.NewFileStorage(c.Upload.BasePath, c.Upload.BaseURL),
		RateLimiter: pkg.NewRateLimiter(redisClient),
	}
}
