package logic

import (
	"context"
	"errors"
	"time"

	"faceTest/backend/auth/internal/svc"
	"faceTest/backend/common/model"
	"faceTest/backend/common/pkg"
	authpb "faceTest/backend/proto/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginLogic 登录逻辑
type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 用户登录
func (l *LoginLogic) Login(req *authpb.LoginReq) (*authpb.LoginResp, error) {
	// 限流检查
	rateKey := "rate_limit:login:" + req.Username
	allowed, err := l.svcCtx.RateLimiter.Allow(l.ctx, rateKey, 60, 5)
	if err == nil && !allowed {
		return &authpb.LoginResp{
			Code: 0,
			Msg:  "登录尝试过于频繁，请稍后再试",
		}, nil
	}

	// 查询用户
	var user model.User
	result := l.svcCtx.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &authpb.LoginResp{
				Code: 0,
				Msg:  "用户名不存在",
			}, nil
		}
		l.Errorf("查询用户失败: %v", result.Error)
		return &authpb.LoginResp{
			Code: 0,
			Msg:  "系统错误",
		}, nil
	}

	// 验证密码（兼容bcrypt和明文）
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// 尝试明文比较（兼容旧数据）
		if user.Password != req.Password {
			return &authpb.LoginResp{
				Code: 0,
				Msg:  "密码错误",
			}, nil
		}
		// 明文密码匹配，升级为bcrypt
		if hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err == nil {
			l.svcCtx.DB.Model(&user).Update("password", string(hashed))
		}
	}

	// 生成JWT Token
	token, err := pkg.GenerateToken(user.ID, user.Username, user.Name, user.Role)
	if err != nil {
		l.Errorf("生成Token失败: %v", err)
		return &authpb.LoginResp{
			Code: 0,
			Msg:  "生成Token失败",
		}, nil
	}

	// 缓存用户信息到Redis
	userInfo := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"name":     user.Name,
		"role":     user.Role,
	}
	cacheKey := "user:info:" + token
	_ = l.svcCtx.RedisCache.Set(l.ctx, cacheKey, userInfo, time.Duration(l.svcCtx.Config.JWT.Expire)*time.Hour)

	l.Infof("用户登录成功: username=%s, userID=%d", user.Username, user.ID)

	return &authpb.LoginResp{
		Code:     1,
		Msg:      "登录成功",
		Token:    token,
		UserId:   int64(user.ID),
		Username: user.Username,
		Name:     user.Name,
	}, nil
}
