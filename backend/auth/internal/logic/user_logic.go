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

// GetUserLogic 获取用户信息逻辑
type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUser 获取用户信息
func (l *GetUserLogic) GetUser(req *authpb.GetUserReq) (*authpb.UserResp, error) {
	var user model.User
	result := l.svcCtx.DB.Where("id = ?", req.UserId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &authpb.UserResp{}, errors.New("用户不存在")
		}
		l.Errorf("查询用户失败: %v", result.Error)
		return &authpb.UserResp{}, result.Error
	}

	var birthday string
	if user.Birthday != nil {
		birthday = user.Birthday.Format("2006-01-02")
	}

	return &authpb.UserResp{
		Id:         int64(user.ID),
		Username:   user.Username,
		Name:       user.Name,
		Avatar:     user.Avatar,
		Email:      user.Email,
		Phone:      user.Phone,
		Gender:     int32(user.Gender),
		Birthday:   birthday,
		Address:    user.Address,
		Role:       user.Role,
		CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateUserLogic 更新用户信息逻辑
type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUser 更新用户信息（使用乐观锁）
func (l *UpdateUserLogic) UpdateUser(req *authpb.UpdateUserReq) (*authpb.UpdateUserResp, error) {
	// 先查询用户获取当前版本号
	var user model.User
	result := l.svcCtx.DB.Where("id = ?", req.UserId).First(&user)
	if result.Error != nil {
		return &authpb.UpdateUserResp{
			Code: 0,
			Msg:  "用户不存在",
		}, nil
	}

	// 更新字段
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Password != "" {
		// 密码加密
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return &authpb.UpdateUserResp{Code: 0, Msg: "密码加密失败"}, nil
		}
		updates["password"] = hashed
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Gender != 0 {
		updates["gender"] = req.Gender
	}
	if req.Birthday != "" {
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			updates["birthday"] = birthday
		}
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}

	// 使用乐观锁更新
	updates["version"] = user.Version + 1
	result = l.svcCtx.DB.Model(&model.User{}).
		Where("id = ? AND version = ?", req.UserId, user.Version).
		Updates(updates)

	if result.Error != nil {
		l.Errorf("更新用户失败: %v", result.Error)
		return &authpb.UpdateUserResp{
			Code: 0,
			Msg:  "更新失败",
		}, nil
	}

	if result.RowsAffected == 0 {
		// 乐观锁冲突，重试
		l.Infof("乐观锁冲突，用户ID=%d，重试中...", req.UserId)
		updates["version"] = gorm.Expr("version + 1")
		result = l.svcCtx.DB.Model(&model.User{}).
			Where("id = ?", req.UserId).
			Updates(updates)
		if result.Error != nil {
			return &authpb.UpdateUserResp{Code: 0, Msg: "更新失败"}, nil
		}
	}

	return &authpb.UpdateUserResp{
		Code: 1,
		Msg:  "更新成功",
	}, nil
}

// ValidateTokenLogic 验证Token逻辑
type ValidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateTokenLogic {
	return &ValidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ValidateToken 验证JWT Token
func (l *ValidateTokenLogic) ValidateToken(req *authpb.ValidateTokenReq) (*authpb.ValidateTokenResp, error) {
	claims, err := pkg.ParseToken(req.Token)
	if err != nil {
		return &authpb.ValidateTokenResp{
			Valid: false,
		}, nil
	}

	return &authpb.ValidateTokenResp{
		Valid:    true,
		UserId:   int64(claims.UserID),
		Username: claims.Username,
		Name:     claims.Name,
		Role:     claims.Role,
	}, nil
}
