package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// DeepSeekClient DeepSeek API客户端
type DeepSeekClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// DeepSeekMessage 消息
type DeepSeekMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// DeepSeekRequest 请求体
type DeepSeekRequest struct {
	Model       string            `json:"model"`
	Messages    []DeepSeekMessage `json:"messages"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	Stream      bool              `json:"stream,omitempty"`
}

// DeepSeekResponse 响应体
type DeepSeekResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// DiagnosisResult 诊断结果（DeepSeek返回解析）
type DiagnosisAIResult struct {
	IsGlandFace             bool    `json:"isGlandFace"`
	Level                   string  `json:"level"`
	VisualizationDescription string  `json:"visualizationDescription"`
	ReferenceCases          []string `json:"referenceCases,omitempty"`
}

// NewDeepSeekClient 创建DeepSeek客户端
func NewDeepSeekClient(apiKey, model string) *DeepSeekClient {
	return &DeepSeekClient{
		apiKey:  apiKey,
		baseURL: "https://api.deepseek.com/v1/chat/completions",
		model:   model,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// ChatCompletion 聊天补全
func (dc *DeepSeekClient) ChatCompletion(ctx context.Context, messages []DeepSeekMessage) (*DeepSeekResponse, error) {
	reqBody := DeepSeekRequest{
		Model:       dc.model,
		Messages:    messages,
		MaxTokens:   4096,
		Temperature: 0.3, // 低温度保证一致性
		Stream:      false,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("请求序列化失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", dc.baseURL, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+dc.apiKey)

	resp, err := dc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logx.Errorf("DeepSeek API错误: status=%d, body=%s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("API返回错误: status=%d", resp.StatusCode)
	}

	var result DeepSeekResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("响应解析失败: %w", err)
	}

	logx.Infof("DeepSeek API调用成功: model=%s, tokens=%d", dc.model, result.Usage.TotalTokens)
	return &result, nil
}

// Diagnose 诊断（RAG增强 + 客观几何测量）
// measurements: MediaPipe Face Mesh 提取的真实面部几何测量值（可为空）
func (dc *DeepSeekClient) Diagnose(ctx context.Context, featureDescription string, referenceCases string, measurements map[string]float64, extraContext string) (*DiagnosisAIResult, error) {
	systemPrompt := `你是一名专业的耳鼻喉科医学影像分析专家，专注于儿童腺样体面容的早期筛查。

【腺样体面容诊断标准】
- 上唇短翘，上翻明显
- 下颌后缩，颏部发育不足
- 鼻唇沟变浅或消失
- 牙列拥挤不齐，前牙开合
- 张口呼吸，唇肌松弛
- 面部比例失调（下面高增加）

【证据采信优先级（务必遵守）】
1. 首要证据：下方【客观面部几何测量值】（MediaPipe Face Mesh 提取的真实解剖学测量）。这是可解释的客观硬指标，必须优先采信；其中：
   - 张口度/口宽比偏大 => 张口呼吸倾向
   - 面下1/3比例偏大、面宽/面高比偏小 => 长面型（腺样体面容特征）
   - 鼻唇角偏小、面凸角异常 => 下颌后缩/面型改变
2. 辅助参考：【视觉特征摘要（EfficientNet-B3 派生）】为模型统计特征，仅作补充，不可单独作为结论依据。
3. 若提供【多图综合说明】，请重点评估多张照片之间的一致性；当各图关键测量或特征存在明显分歧时，必须明确说明分歧点与最终采信依据，不要简单平均掩盖矛盾。

【输出要求】
请综合分析后，严格按照以下JSON格式返回结果，不要包含任何其他文字：
{
  "isGlandFace": true/false,
  "level": "轻度"/"中度"/"重度"（如果isGlandFace为false则level为"非腺样体面容"）,
  "visualizationDescription": "专业的面部特征医学描述，内容要专业且丰富，需结合客观测量值与视觉特征"
}`

	// 组织客观几何测量值文本
	measurementText := "（本次未获取到可量化的面部几何测量值，请基于视觉特征判断）"
	if len(measurements) > 0 {
		labels := map[string]string{
			"mouth_aperture_ratio":       "张口度/口宽比",
			"lower_face_ratio":           "面下1/3比例",
			"face_width_height_ratio":    "面宽/面高比",
			"nasolabial_angle_deg":       "鼻唇角(度)",
			"facial_convexity_angle_deg": "面凸角(度，下颌后缩proxy)",
		}
		var sb strings.Builder
		sb.WriteString("面部几何测量值（由 MediaPipe Face Mesh 提取，客观辅助证据）:\n")
		for k, v := range measurements {
			label := labels[k]
			if label == "" {
				label = k
			}
			sb.WriteString(fmt.Sprintf("- %s: %.4f\n", label, v))
		}
		measurementText = sb.String()
	}

	userPrompt := fmt.Sprintf(`请基于以下信息，判断该儿童是否具有腺样体面容特征：

【相似历史病例参考】
%s

【客观面部几何测量值】
%s

【视觉特征摘要（EfficientNet-B3 派生，辅助参考）】
%s%s
请结合以上所有信息，以客观几何测量值为首要依据进行综合判断，并严格按照指定JSON格式返回诊断结果。`, referenceCases, measurementText, featureDescription, extraContextBlock(extraContext))

	messages := []DeepSeekMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	resp, err := dc.ChatCompletion(ctx, messages)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("API返回为空")
	}

	content := resp.Choices[0].Message.Content

	// 解析JSON结果
	var result DiagnosisAIResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// 尝试从markdown代码块中提取JSON
		jsonStr := extractJSON(content)
		if jsonStr == "" {
			return nil, fmt.Errorf("无法解析诊断结果: %s", content)
		}
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			return nil, fmt.Errorf("JSON解析失败: %w, content=%s", err, content)
		}
	}

	return &result, nil
}

// Chat 对话补全（用于AI助手）
func (dc *DeepSeekClient) Chat(ctx context.Context, messages []DeepSeekMessage) (string, error) {
	resp, err := dc.ChatCompletion(ctx, messages)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("API返回为空")
	}
	return resp.Choices[0].Message.Content, nil
}

// extraContextBlock 将多图综合说明包裹为提示词段落；为空时返回空串。
func extraContextBlock(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return "\n【多图综合说明】\n" + s + "\n"
}

// extractJSON 从可能包含markdown代码块的文本中提取JSON
func extractJSON(s string) string {
	start := -1
	end := -1

	// 尝试找 ```json ... ``` 格式
	for i := 0; i < len(s)-7; i++ {
		if s[i:i+7] == "```json" {
			start = i + 7
			break
		}
	}
	if start == -1 {
		// 尝试找 { 开始
		for i := 0; i < len(s); i++ {
			if s[i] == '{' {
				start = i
				break
			}
		}
	}
	if start == -1 {
		return ""
	}

	// 找到最后一个 }
	for i := len(s) - 1; i >= start; i-- {
		if s[i] == '}' {
			end = i + 1
			break
		}
	}
	if end == -1 {
		return ""
	}

	return s[start:end]
}
