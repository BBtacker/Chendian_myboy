package svc

import (
	"faceTest/backend/common/pkg"
	"faceTest/backend/upload/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
)

// ServiceContext 图像上传服务上下文
type ServiceContext struct {
	Config      config.Config
	FileStorage *pkg.FileStorage
}

// NewServiceContext 创建上传服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	storage := pkg.NewFileStorage(c.Upload.BasePath, c.Upload.BaseURL)
	logx.Info("图像上传服务上下文初始化完成")
	return &ServiceContext{
		Config:      c,
		FileStorage: storage,
	}
}
