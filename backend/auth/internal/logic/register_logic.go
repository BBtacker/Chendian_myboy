package logic

import (
	"context"
	"errors"
	"time"

	"faceTest/backend/auth/internal/svc"
	"faceTest/backend/common/model"
	authpb "faceTest/backend/proto/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterLogic 注册逻辑
type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Register 用户注册
func (l *RegisterLogic) Register(req *authpb.RegisterReq) (*authpb.RegisterResp, error) {
	// 检查用户名是否已存在
	var count int64
	l.svcCtx.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return &authpb.RegisterResp{
			Code: 0,
			Msg:  "用户名已存在",
		}, nil
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码加密失败: %v", err)
		return &authpb.RegisterResp{
			Code: 0,
			Msg:  "系统错误",
		}, nil
	}

	// 设置默认姓名
	name := req.Name
	if name == "" {
		name = req.Username
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Name:     name,
		Role:     "doctor",
	}

	if err := l.svcCtx.DB.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &authpb.RegisterResp{
				Code: 0,
				Msg:  "用户名已存在",
			}, nil
		}
		l.Errorf("创建用户失败: %v", err)
		return &authpb.RegisterResp{
			Code: 0,
			Msg:  "注册失败",
		}, nil
	}

	_ = time.Now() // 避免unused import

	l.Infof("用户注册成功: username=%s, userID=%d", user.Username, user.ID)

	return &authpb.RegisterResp{
		Code: 1,
		Msg:  "注册成功",
	}, nil
}
