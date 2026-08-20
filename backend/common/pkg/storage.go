package pkg

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

// FileStorage 文件存储（本地存储，可扩展为OSS）
type FileStorage struct {
	basePath    string
	baseURL     string
	maxFileSize int64
	allowedTypes map[string]bool
}

// NewFileStorage 创建文件存储
func NewFileStorage(basePath, baseURL string) *FileStorage {
	// 确保目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		logx.Errorf("创建存储目录失败: %v", err)
	}

	return &FileStorage{
		basePath:    basePath,
		baseURL:     baseURL,
		maxFileSize: 10 << 20, // 10MB
		allowedTypes: map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/webp": true,
			"image/gif":  true,
		},
	}
}

// UploadFile 上传文件
func (fs *FileStorage) UploadFile(file multipart.File, header *multipart.FileHeader) (path, url string, err error) {
	// 验证文件大小
	if header.Size > fs.maxFileSize {
		return "", "", fmt.Errorf("文件大小超过限制（最大%dMB）", fs.maxFileSize>>20)
	}

	// 验证文件类型
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		case ".gif":
			contentType = "image/gif"
		default:
			return "", "", fmt.Errorf("不支持的文件类型: %s", ext)
		}
	}
	if !fs.allowedTypes[contentType] {
		return "", "", fmt.Errorf("不支持的文件类型: %s", contentType)
	}

	// 按日期分目录
	dateDir := time.Now().Format("2006/01")
	ext := filepath.Ext(header.Filename)
	newFileName := uuid.New().String() + ext
	relativePath := filepath.Join(dateDir, newFileName)
	fullPath := filepath.Join(fs.basePath, relativePath)

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(fullPath) // 清理失败的文件
		return "", "", fmt.Errorf("写入文件失败: %w", err)
	}

	// 构建路径和URL
	path = relativePath
	url = fmt.Sprintf("%s/%s", strings.TrimSuffix(fs.baseURL, "/"), relativePath)

	logx.Infof("文件上传成功: path=%s, size=%d", path, header.Size)
	return path, url, nil
}

// UploadFromBytes 从字节上传文件
func (fs *FileStorage) UploadFromBytes(data []byte, filename string) (path, url string, err error) {
	// 验证大小
	if int64(len(data)) > fs.maxFileSize {
		return "", "", fmt.Errorf("文件大小超过限制（最大%dMB）", fs.maxFileSize>>20)
	}

	// 按日期分目录
	dateDir := time.Now().Format("2006/01")
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg"
	}
	newFileName := uuid.New().String() + ext
	relativePath := filepath.Join(dateDir, newFileName)
	fullPath := filepath.Join(fs.basePath, relativePath)

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", "", fmt.Errorf("写入文件失败: %w", err)
	}

	path = relativePath
	url = fmt.Sprintf("%s/%s", strings.TrimSuffix(fs.baseURL, "/"), relativePath)

	logx.Infof("文件上传成功(bytes): path=%s, size=%d", path, len(data))
	return path, url, nil
}

// DeleteFile 删除文件
func (fs *FileStorage) DeleteFile(path string) error {
	fullPath := filepath.Join(fs.basePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// FileExists 检查文件是否存在
func (fs *FileStorage) FileExists(path string) bool {
	fullPath := filepath.Join(fs.basePath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// GetFullPath 获取文件完整路径
func (fs *FileStorage) GetFullPath(path string) string {
	return filepath.Join(fs.basePath, path)
}
