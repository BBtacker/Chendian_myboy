package logic

import (
	"context"
	"errors"

	"faceTest/backend/auth/internal/svc"
	"faceTest/backend/common/model"
	authpb "faceTest/backend/proto/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// UploadAvatarLogic 头像上传逻辑
type UploadAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UploadAvatar 上传头像
func (l *UploadAvatarLogic) UploadAvatar(req *authpb.UploadAvatarReq) (*authpb.UploadAvatarResp, error) {
	// 验证用户存在
	var user model.User
	result := l.svcCtx.DB.Where("id = ?", req.UserId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &authpb.UploadAvatarResp{
				Code: 0,
				Msg:  "用户不存在",
			}, nil
		}
		return &authpb.UploadAvatarResp{
			Code: 0,
			Msg:  "系统错误",
		}, nil
	}

	// 存储文件
	_, url, err := l.svcCtx.FileStorage.UploadFromBytes(req.FileData, req.Filename)
	if err != nil {
		l.Errorf("头像上传失败: %v", err)
		return &authpb.UploadAvatarResp{
			Code: 0,
			Msg:  "头像上传失败: " + err.Error(),
		}, nil
	}

	// 更新用户头像（使用乐观锁）
	result = l.svcCtx.DB.Model(&model.User{}).
		Where("id = ? AND version = ?", req.UserId, user.Version).
		Updates(map[string]interface{}{
			"avatar":  url,
			"version": user.Version + 1,
		})

	if result.RowsAffected == 0 {
		// 乐观锁冲突，直接更新
		l.svcCtx.DB.Model(&model.User{}).Where("id = ?", req.UserId).Update("avatar", url)
	}

	l.Infof("头像上传成功: userID=%d, url=%s", req.UserId, url)

	return &authpb.UploadAvatarResp{
		Code:      1,
		Msg:       "上传成功",
		AvatarUrl: url,
	}, nil
}
