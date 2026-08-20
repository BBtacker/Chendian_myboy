package logic

import (
	"context"
	"fmt"
	"time"

	"faceTest/backend/common/model"
	"faceTest/backend/report/internal/svc"
	reportpb "faceTest/backend/proto/report"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

// ExportExcelLogic Excel导出逻辑
type ExportExcelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExportExcelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportExcelLogic {
	return &ExportExcelLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// ExportExcel 导出Excel
func (l *ExportExcelLogic) ExportExcel(req *reportpb.ExportReq) (*reportpb.ExportResp, error) {
	// 查询数据
	results, err := l.queryResults(req)
	if err != nil {
		return &reportpb.ExportResp{Code: 0, Msg: "查询数据失败"}, nil
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "检测记录"
	f.SetSheetName(f.GetSheetName(0), sheetName)

	// 表头
	headers := []string{"ID", "用户ID", "图片路径", "是否腺样体面容", "等级", "可视化描述", "检测时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}
	// 设置表头样式
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#E0E0E0"}},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheetName, "A1", "G1", style)

	// 数据行
	for i, r := range results {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), r.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), r.UserID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), r.ImagePath)
		isGland := "否"
		if r.IsGlandFace {
			isGland = "是"
		}
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), isGland)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), r.Level)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), r.VisualizationDescription)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), r.TestTime.Format("2006-01-02 15:04:05"))
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 8)
	f.SetColWidth(sheetName, "B", "B", 10)
	f.SetColWidth(sheetName, "C", "C", 40)
	f.SetColWidth(sheetName, "D", "D", 15)
	f.SetColWidth(sheetName, "E", "E", 10)
	f.SetColWidth(sheetName, "F", "F", 60)
	f.SetColWidth(sheetName, "G", "G", 22)

	// 保存到buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		l.Errorf("生成Excel失败: %v", err)
		return &reportpb.ExportResp{Code: 0, Msg: "生成Excel失败"}, nil
	}

	filename := fmt.Sprintf("检测记录_%s.xlsx", time.Now().Format("20060102_150405"))

	l.Infof("Excel导出成功: %d条记录", len(results))

	return &reportpb.ExportResp{
		Code:        1,
		Msg:         "导出成功",
		FileData:    buf.Bytes(),
		Filename:    filename,
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}, nil
}

// ExportPDFLogic PDF导出逻辑
type ExportPDFLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExportPDFLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportPDFLogic {
	return &ExportPDFLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// ExportPDF 导出PDF（生成HTML格式，浏览器可另存为PDF）
func (l *ExportPDFLogic) ExportPDF(req *reportpb.ExportReq) (*reportpb.ExportResp, error) {
	results, err := l.queryResults(req)
	if err != nil {
		return &reportpb.ExportResp{Code: 0, Msg: "查询数据失败"}, nil
	}

	// 生成HTML报告（兼容前端PDF导出逻辑）
	html := l.generateHTMLReport(results)

	filename := fmt.Sprintf("检测记录_%s.pdf.html", time.Now().Format("20060102_150405"))

	l.Infof("PDF(HTML)导出成功: %d条记录", len(results))

	return &reportpb.ExportResp{
		Code:        1,
		Msg:         "导出成功",
		FileData:    []byte(html),
		Filename:    filename,
		ContentType: "text/html; charset=utf-8",
	}, nil
}

// ExportSinglePDFLogic 单条PDF导出
type ExportSinglePDFLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExportSinglePDFLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportSinglePDFLogic {
	return &ExportSinglePDFLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ExportSinglePDFLogic) ExportSinglePDF(req *reportpb.SingleExportReq) (*reportpb.ExportResp, error) {
	var result model.DiagnosisResult
	if err := l.svcCtx.DB.Where("id = ? AND user_id = ?", req.ResultId, req.UserId).First(&result).Error; err != nil {
		return &reportpb.ExportResp{Code: 0, Msg: "记录不存在"}, nil
	}

	html := l.generateSingleHTMLReport(&result)
	filename := fmt.Sprintf("检测报告_%d_%s.pdf.html", result.ID, time.Now().Format("20060102"))

	return &reportpb.ExportResp{
		Code:        1,
		Msg:         "导出成功",
		FileData:    []byte(html),
		Filename:    filename,
		ContentType: "text/html; charset=utf-8",
	}, nil
}

// queryResults 查询结果数据
func (l *ExportExcelLogic) queryResults(req *reportpb.ExportReq) ([]model.DiagnosisResult, error) {
	query := l.svcCtx.DB.Model(&model.DiagnosisResult{}).Where("user_id = ?", req.UserId)

	if req.HasFilter {
		if req.IsGlandFace {
			query = query.Where("is_gland_face = ?", true)
		}
		if req.Level != "" {
			query = query.Where("level = ?", req.Level)
		}
	}
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.StartDate); err == nil {
			query = query.Where("test_time >= ?", t)
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.EndDate); err == nil {
			query = query.Where("test_time <= ?", t)
		}
	}

	var results []model.DiagnosisResult
	query.Order("test_time DESC").Find(&results)
	return results, nil
}

// queryResults for ExportPDFLogic
func (l *ExportPDFLogic) queryResults(req *reportpb.ExportReq) ([]model.DiagnosisResult, error) {
	query := l.svcCtx.DB.Model(&model.DiagnosisResult{}).Where("user_id = ?", req.UserId)

	if req.HasFilter {
		if req.IsGlandFace {
			query = query.Where("is_gland_face = ?", true)
		}
		if req.Level != "" {
			query = query.Where("level = ?", req.Level)
		}
	}
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.StartDate); err == nil {
			query = query.Where("test_time >= ?", t)
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.EndDate); err == nil {
			query = query.Where("test_time <= ?", t)
		}
	}

	var results []model.DiagnosisResult
	query.Order("test_time DESC").Find(&results)
	return results, nil
}

// generateHTMLReport 生成HTML报告
func (l *ExportPDFLogic) generateHTMLReport(results []model.DiagnosisResult) string {
	html := `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>检测记录报告</title>
<style>
body { font-family: Arial, sans-serif; margin: 20px; }
h1 { text-align: center; color: #333; }
table { width: 100%; border-collapse: collapse; margin-top: 20px; }
th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
th { background-color: #f2f2f2; font-weight: bold; }
tr:nth-child(even) { background-color: #f9f9f9; }
.footer { margin-top: 30px; text-align: center; color: #666; font-size: 12px; }
.image-preview { max-width: 100px; max-height: 100px; }
</style>
</head>
<body>
<h1>检测记录报告</h1>
<table>
<thead>
<tr>
<th>ID</th><th>用户ID</th><th>图片</th><th>是否腺样体面容</th><th>等级</th><th>可视化描述</th><th>检测时间</th>
</tr>
</thead>
<tbody>`

	for _, r := range results {
		isGland := "否"
		if r.IsGlandFace {
			isGland = "是"
		}
		imageTag := ""
		if r.ImageURL != "" {
			imageTag = fmt.Sprintf(`<a href="%s" target="_blank"><img src="%s" class="image-preview" alt="检测图片"/></a>`, r.ImageURL, r.ImageURL)
		}
		html += fmt.Sprintf(`<tr>
<td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td>
</tr>`, r.ID, r.UserID, imageTag, isGland, r.Level, r.VisualizationDescription, r.TestTime.Format("2006-01-02 15:04:05"))
	}

	html += `</tbody>
</table>
<div class="footer">报告生成时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</div>
</body>
</html>`
	return html
}

// generateSingleHTMLReport 生成单条HTML报告
func (l *ExportSinglePDFLogic) generateSingleHTMLReport(r *model.DiagnosisResult) string {
	isGland := "否"
	if r.IsGlandFace {
		isGland = "是"
	}
	imageTag := ""
	if r.ImageURL != "" {
		imageTag = fmt.Sprintf(`<img src="%s" style="max-width:300px;max-height:300px;" alt="检测图片"/>`, r.ImageURL)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>腺样体面容检测报告</title>
<style>
body { font-family: Arial, sans-serif; margin: 40px; max-width: 800px; margin: 0 auto; }
.header { text-align: center; margin-bottom: 30px; }
.header h1 { color: #2c3e50; }
.info-card { background: #f8f9fa; border-radius: 8px; padding: 20px; margin: 20px 0; }
.info-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #eee; }
.info-label { font-weight: bold; color: #555; }
.result-badge { display: inline-block; padding: 5px 15px; border-radius: 20px; font-weight: bold; }
.badge-positive { background: #ffe0e0; color: #c0392b; }
.badge-negative { background: #e0ffe0; color: #27ae60; }
.image-section { text-align: center; margin: 20px 0; }
.description { background: #fff; border-left: 4px solid #3498db; padding: 15px; margin: 20px 0; }
.footer { margin-top: 40px; text-align: center; color: #666; font-size: 12px; border-top: 1px solid #eee; padding-top: 20px; }
</style>
</head>
<body>
<div class="header">
<h1>腺样体面容检测报告</h1>
<p>报告编号: #%d</p>
</div>
<div class="info-card">
<div class="info-row"><span class="info-label">检测时间</span><span>%s</span></div>
<div class="info-row"><span class="info-label">检测结论</span><span class="result-badge %s">%s</span></div>
<div class="info-row"><span class="info-label">严重程度</span><span>%s</span></div>
</div>
<div class="image-section">%s</div>
<div class="description">
<h3>医学描述</h3>
<p>%s</p>
</div>
<div class="footer">
<p>本报告由AI腺样体面容智能筛查系统自动生成</p>
<p>生成时间: %s</p>
</div>
</body>
</html>`, r.ID, r.TestTime.Format("2006-01-02 15:04:05"),
		map[bool]string{true: "badge-positive", false: "badge-negative"}[r.IsGlandFace],
		isGland, r.Level, imageTag, r.VisualizationDescription,
		time.Now().Format("2006-01-02 15:04:05"))
}
