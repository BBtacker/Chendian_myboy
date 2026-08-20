package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// FeatureExtractorConfig 特征提取服务配置
type FeatureExtractorConfig struct {
	URL     string // 特征提取服务地址，如 http://127.0.0.1:8085
	Timeout int    // 超时（秒），默认 10
}

// FeatureExtractor 特征提取服务客户端
// 调用 Python FastAPI 服务（EfficientNet-B3 ONNX 推理）
type FeatureExtractor struct {
	baseURL string
	client  *http.Client
}

// ExtractResult 特征提取结果
type ExtractResult struct {
	Code           int                `json:"code"`
	Msg            string             `json:"msg"`
	FeatureVector  []float32          `json:"feature_vector"`
	Dimension      int                `json:"dimension"`
	Description    string             `json:"description"`
	FaceDetected   bool               `json:"face_detected"`
	Measurements   map[string]float64 `json:"measurements"`
	ElapsedMs      int                `json:"elapsed_ms"`
}

// NewFeatureExtractor 创建特征提取客户端
func NewFeatureExtractor(cfg FeatureExtractorConfig) *FeatureExtractor {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10
	}
	return &FeatureExtractor{
		baseURL: cfg.URL,
		client:  &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

// Extract 提取图片特征向量
// imagePath 为本地图片完整路径
func (fe *FeatureExtractor) Extract(ctx context.Context, imagePath string) (*ExtractResult, error) {
	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("打开图片失败: %w", err)
	}
	defer file.Close()

	// 构造 multipart 请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 用 CreatePart 显式设置 Content-Type，否则 Python /extract 会因 application/octet-stream 拒绝
	ct := imageContentType(imagePath)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="image"; filename="%s"`, filepath.Base(imagePath)))
	if ct != "" {
		partHeader.Set("Content-Type", ct)
	}
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return nil, fmt.Errorf("创建表单失败: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("读取图片失败: %w", err)
	}
	writer.Close()

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fe.baseURL+"/extract", body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := fe.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用特征服务失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("特征服务返回错误: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	var result ExtractResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 1 {
		return nil, fmt.Errorf("特征提取失败: %s", result.Msg)
	}

	return &result, nil
}

// imageContentType 根据文件扩展名返回图片 MIME 类型
func imageContentType(p string) string {
	switch strings.ToLower(filepath.Ext(p)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	default:
		return ""
	}
}

// Health 健康检查
func (fe *FeatureExtractor) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fe.baseURL+"/health", nil)
	if err != nil {
		return err
	}
	resp, err := fe.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("特征服务健康检查失败: status=%d", resp.StatusCode)
	}
	return nil
}

// logx 保留引用（包内其他文件使用）
var _ = logx.Infof
