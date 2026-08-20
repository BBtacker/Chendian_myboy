package svc

import (
	"faceTest/backend/common/pkg"
	"faceTest/backend/report/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ServiceContext 报告服务上下文
type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	RedisClient *redis.Client
	RedisCache  *pkg.RedisCache
}

// NewServiceContext 创建报告服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logx.Must(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(50)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Cache.Host,
		Password: c.Cache.Password,
		DB:       c.Cache.DB,
	})

	return &ServiceContext{
		Config:      c,
		DB:          db,
		RedisClient: redisClient,
		RedisCache:  pkg.NewRedisCache(c.Cache.Host, c.Cache.Password, c.Cache.DB),
	}
}
