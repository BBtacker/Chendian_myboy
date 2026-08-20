package logic

import (
	"context"

	"faceTest/backend/common/pkg"
	"faceTest/backend/upload/internal/svc"
	uploadpb "faceTest/backend/proto/upload"

	"github.com/zeromicro/go-zero/core/logx"
)

// UploadImageLogic 图片上传逻辑
type UploadImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImageLogic {
	return &UploadImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UploadImage 上传单张图片
func (l *UploadImageLogic) UploadImage(req *uploadpb.UploadImageReq) (*uploadpb.UploadImageResp, error) {
	// 验证文件大小
	maxSize := l.svcCtx.Config.Upload.MaxFileSize
	if maxSize == 0 {
		maxSize = 10 << 20 // 默认10MB
	}
	if int64(len(req.FileData)) > maxSize {
		return &uploadpb.UploadImageResp{
			Code: 0,
			Msg:  "文件大小超过限制",
		}, nil
	}

	// 存储文件
	path, url, err := l.svcCtx.FileStorage.UploadFromBytes(req.FileData, req.Filename)
	if err != nil {
		l.Errorf("图片上传失败: %v", err)
		return &uploadpb.UploadImageResp{
			Code: 0,
			Msg:  "图片上传失败: " + err.Error(),
		}, nil
	}

	l.Infof("图片上传成功: path=%s, size=%d", path, len(req.FileData))

	return &uploadpb.UploadImageResp{
		Code:      1,
		Msg:       "上传成功",
		ImagePath: path,
		ImageUrl:  url,
		FileSize:  int64(len(req.FileData)),
	}, nil
}

// BatchUploadImagesLogic 批量上传逻辑
type BatchUploadImagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchUploadImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUploadImagesLogic {
	return &BatchUploadImagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchUploadImages 批量上传图片
func (l *BatchUploadImagesLogic) BatchUploadImages(req *uploadpb.BatchUploadReq) (*uploadpb.BatchUploadResp, error) {
	var results []*uploadpb.UploadImageResp

	for _, file := range req.Files {
		uploadLogic := NewUploadImageLogic(l.ctx, l.svcCtx)
		resp, err := uploadLogic.UploadImage(file)
		if err != nil {
			resp = &uploadpb.UploadImageResp{
				Code: 0,
				Msg:  "上传失败: " + err.Error(),
			}
		}
		results = append(results, resp)
	}

	_ = pkg.NewFileStorage // avoid unused import

	return &uploadpb.BatchUploadResp{
		Code:    1,
		Msg:     "批量上传完成",
		Results: results,
	}, nil
}
