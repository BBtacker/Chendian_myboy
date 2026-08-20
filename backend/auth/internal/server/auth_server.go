package server

import (
	"context"

	"faceTest/backend/auth/internal/logic"
	"faceTest/backend/auth/internal/svc"
	authpb "faceTest/backend/proto/auth"
)

// AuthServer 认证服务gRPC服务器
type AuthServer struct {
	svcCtx *svc.ServiceContext
	authpb.UnimplementedAuthServer
}

// NewAuthServer 创建认证服务器
func NewAuthServer(svcCtx *svc.ServiceContext) *AuthServer {
	return &AuthServer{svcCtx: svcCtx}
}

// Login 用户登录
func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(req)
}

// Register 用户注册
func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterReq) (*authpb.RegisterResp, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(req)
}

// GetUser 获取用户信息
func (s *AuthServer) GetUser(ctx context.Context, req *authpb.GetUserReq) (*authpb.UserResp, error) {
	l := logic.NewGetUserLogic(ctx, s.svcCtx)
	return l.GetUser(req)
}

// UpdateUser 更新用户信息
func (s *AuthServer) UpdateUser(ctx context.Context, req *authpb.UpdateUserReq) (*authpb.UpdateUserResp, error) {
	l := logic.NewUpdateUserLogic(ctx, s.svcCtx)
	return l.UpdateUser(req)
}

// ValidateToken 验证Token
func (s *AuthServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenReq) (*authpb.ValidateTokenResp, error) {
	l := logic.NewValidateTokenLogic(ctx, s.svcCtx)
	return l.ValidateToken(req)
}

// UploadAvatar 上传头像
func (s *AuthServer) UploadAvatar(ctx context.Context, req *authpb.UploadAvatarReq) (*authpb.UploadAvatarResp, error) {
	l := logic.NewUploadAvatarLogic(ctx, s.svcCtx)
	return l.UploadAvatar(req)
}
