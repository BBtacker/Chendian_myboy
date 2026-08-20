package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// 滑动窗口限流Lua脚本
// KEYS[1] = 限流key
// ARGV[1] = 窗口大小（秒）
// ARGV[2] = 最大请求数
// ARGV[3] = 当前时间戳（毫秒）
const slidingWindowScript = `
local key = KEYS[1]
local window = tonumber(ARGV[1])
local maxRequests = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local windowStart = now - window * 1000

-- 移除窗口外的记录
redis.call('ZREMRANGEBYSCORE', key, 0, windowStart)

-- 获取当前窗口内的请求数
local count = redis.call('ZCARD', key)

if count < maxRequests then
    -- 添加当前请求
    redis.call('ZADD', key, now, now .. '-' .. math.random(1, 1000000))
    redis.call('EXPIRE', key, window + 1)
    return 1
else
    return 0
end
`

// RateLimiter 基于Redis Lua脚本的滑动窗口限流器
type RateLimiter struct {
	client *redis.Client
}

// NewRateLimiter 创建限流器
func NewRateLimiter(redisClient *redis.Client) *RateLimiter {
	return &RateLimiter{client: redisClient}
}

// Allow 检查是否允许请求
// key: 限流标识（如 "rate_limit:diagnosis:user_123"）
// window: 时间窗口（秒）
// maxRequests: 窗口内最大请求数
func (rl *RateLimiter) Allow(ctx context.Context, key string, window int, maxRequests int) (bool, error) {
	now := time.Now().UnixMilli()
	result, err := rl.client.Eval(ctx, slidingWindowScript, []string{key},
		window, maxRequests, now).Int()
	if err != nil {
		logx.Errorf("限流检查失败: key=%s, err=%v", key, err)
		// 限流器故障时放行，避免影响正常业务
		return true, nil
	}
	return result == 1, nil
}

// AllowWithBackoff 带退避的限流检查
func (rl *RateLimiter) AllowWithBackoff(ctx context.Context, key string, window int, maxRequests int) (bool, time.Duration, error) {
	allowed, err := rl.Allow(ctx, key, window, maxRequests)
	if err != nil {
		return true, 0, err
	}
	if !allowed {
		// 计算需要等待的时间
		windowStart := time.Now().UnixMilli() - int64(window*1000)
		oldest, err := rl.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Min: fmt.Sprintf("%d", windowStart),
			Max: "+inf",
			Count: 1,
		}).Result()
		if err == nil && len(oldest) > 0 {
			// 简单返回窗口大小的1/4作为退避时间
			return false, time.Duration(window/4) * time.Second, nil
		}
		return false, time.Duration(window) * time.Second, nil
	}
	return true, 0, nil
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 诊断接口限流：每用户每分钟最多10次
	DiagnosisWindow    int
	DiagnosisMaxReqs   int
	// 登录接口限流：每IP每分钟最多5次
	LoginWindow        int
	LoginMaxReqs       int
	// 注册接口限流：每IP每小时最多3次
	RegisterWindow     int
	RegisterMaxReqs    int
	// 全局限流：每秒最多1000次
	GlobalWindow       int
	GlobalMaxReqs      int
}

// DefaultRateLimitConfig 默认限流配置
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		DiagnosisWindow:   60,
		DiagnosisMaxReqs:  10,
		LoginWindow:       60,
		LoginMaxReqs:      5,
		RegisterWindow:    3600,
		RegisterMaxReqs:   3,
		GlobalWindow:      1,
		GlobalMaxReqs:     1000,
	}
}
