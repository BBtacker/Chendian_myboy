package server

import (
	"context"

	"faceTest/backend/upload/internal/logic"
	"faceTest/backend/upload/internal/svc"
	uploadpb "faceTest/backend/proto/upload"
)

// UploadServer 图像上传服务gRPC服务器
type UploadServer struct {
	svcCtx *svc.ServiceContext
	uploadpb.UnimplementedUploadServer
}

func NewUploadServer(svcCtx *svc.ServiceContext) *UploadServer {
	return &UploadServer{svcCtx: svcCtx}
}

func (s *UploadServer) UploadImage(ctx context.Context, req *uploadpb.UploadImageReq) (*uploadpb.UploadImageResp, error) {
	l := logic.NewUploadImageLogic(ctx, s.svcCtx)
	return l.UploadImage(req)
}

func (s *UploadServer) BatchUploadImages(ctx context.Context, req *uploadpb.BatchUploadReq) (*uploadpb.BatchUploadResp, error) {
	l := logic.NewBatchUploadImagesLogic(ctx, s.svcCtx)
	return l.BatchUploadImages(req)
}
