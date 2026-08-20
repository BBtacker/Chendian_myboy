package handler

import (
	"encoding/json"
	"strings"

	diagnosispb "faceTest/backend/proto/diagnosis"
)

// DiagnosisResultVO 诊断结果视图对象。
// 原前端（Vue3）按 camelCase 约定读取字段（isGlandFace / visualizationDescription /
// imagePath / testTime 等），而 proto 生成的 DiagnosisResultResp 序列化后是
// snake_case（is_gland_face / visualization_description / ...）。这里在网关边界把
// 字段名映射回 camelCase，使原前端无需改动即可正确渲染报告内容。
type DiagnosisResultVO struct {
	ID                       int64   `json:"id"`
	TaskID                   int64   `json:"taskId"`
	UserID                   int64   `json:"userId"`
	ImagePath                string  `json:"imagePath"`
	ImageURL                 string  `json:"imageUrl"`
	IsGlandFace              bool    `json:"isGlandFace"`
	Level                    string  `json:"level"`
	VisualizationDescription string  `json:"visualizationDescription"`
	ReferenceCases           string  `json:"referenceCases"`
	ImagePaths               []string `json:"imagePaths"`
	ImageUrls                []string `json:"imageUrls"`
	SkippedImages            []string `json:"skippedImages"`
	TestTime                 string  `json:"testTime"`
}

// toDiagnosisResultVO 将单个 proto 响应映射为前端友好的 VO。
func toDiagnosisResultVO(r *diagnosispb.DiagnosisResultResp) DiagnosisResultVO {
	if r == nil {
		return DiagnosisResultVO{}
	}
	// 兼容历史：单值 imagePath/imageUrl 对多图场景取数组首张，
	// 避免历史列表页 TestResult.vue 与 Excel/PDF 导出拿到 JSON 数组串而裂图。
	imagePaths := parseImageField(r.ImagePath)
	imageUrls := parseImageField(r.ImageUrl)
	firstPath := ""
	firstURL := ""
	if len(imagePaths) > 0 {
		firstPath = imagePaths[0]
	}
	if len(imageUrls) > 0 {
		firstURL = imageUrls[0]
	}
	return DiagnosisResultVO{
		ID:                       r.Id,
		TaskID:                   r.TaskId,
		UserID:                   r.UserId,
		ImagePath:                firstPath,
		ImageURL:                 firstURL,
		IsGlandFace:              r.IsGlandFace,
		Level:                    r.Level,
		VisualizationDescription: r.VisualizationDescription,
		ReferenceCases:           r.ReferenceCases,
		ImagePaths:               imagePaths,
		ImageUrls:                imageUrls,
		SkippedImages:            parseImageField(r.SkippedImages),
		TestTime:                 r.TestTime,
	}
}

// parseImageField 解析图片路径字段：兼容 JSON 数组（多图）与单值路径（单图）。
func parseImageField(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if strings.HasPrefix(s, "[") {
		var arr []string
		if err := json.Unmarshal([]byte(s), &arr); err == nil {
			return arr
		}
	}
	return []string{s}
}

// toDiagnosisResultVOs 批量映射。
func toDiagnosisResultVOs(records []*diagnosispb.DiagnosisResultResp) []DiagnosisResultVO {
	vos := make([]DiagnosisResultVO, 0, len(records))
	for _, r := range records {
		vos = append(vos, toDiagnosisResultVO(r))
	}
	return vos
}
