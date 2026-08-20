package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// RedisCache Redis缓存封装
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		logx.Errorf("Redis连接失败: %v", err)
	} else {
		logx.Info("Redis连接成功")
	}

	return &RedisCache{client: client}
}

// GetClient 获取原始Redis客户端
func (rc *RedisCache) GetClient() *redis.Client {
	return rc.client
}

// Set 设置缓存
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}
	return rc.client.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

// Exists 检查key是否存在
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := rc.client.Exists(ctx, key).Result()
	return n > 0, err
}

// HSet Hash存储（用于存储特征向量索引）
func (rc *RedisCache) HSet(ctx context.Context, key, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rc.client.HSet(ctx, key, field, data).Err()
}

// HGet Hash获取
func (rc *RedisCache) HGet(ctx context.Context, key, field string, dest interface{}) error {
	data, err := rc.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// HGetAll 获取Hash所有字段
func (rc *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return rc.client.HGetAll(ctx, key).Result()
}

// HDel 删除Hash字段
func (rc *RedisCache) HDel(ctx context.Context, key string, fields ...string) error {
	return rc.client.HDel(ctx, key, fields...).Err()
}

// SetNX 分布式锁
func (rc *RedisCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return rc.client.SetNX(ctx, key, value, expiration).Result()
}

// Incr 自增
func (rc *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return rc.client.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func (rc *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rc.client.Expire(ctx, key, expiration).Err()
}
