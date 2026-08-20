package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"faceTest/backend/common/middleware"
	"faceTest/backend/common/pkg"
	"faceTest/backend/gateway/internal/svc"
	diagnosispb "faceTest/backend/proto/diagnosis"
	reportpb "faceTest/backend/proto/report"
	uploadpb "faceTest/backend/proto/upload"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

// AnalyzeFaceHandler 面容分析（上传图片 + 提交诊断任务）
func AnalyzeFaceHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		// 限制文件大小 10MB
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

		file, header, err := r.FormFile("image")
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("请上传图片文件"))
			return
		}
		defer file.Close()

		// 读取文件内容
		fileData, err := io.ReadAll(file)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("文件读取失败"))
			return
		}

		// Step 1: 调用上传服务存储图片
		uploadCtx, uploadCancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer uploadCancel()

		uploadResp, err := ctx.UploadClient.UploadImage(uploadCtx, &uploadpb.UploadImageReq{
			FileData:    fileData,
			Filename:    header.Filename,
			ContentType: header.Header.Get("Content-Type"),
		})
		if err != nil || uploadResp.Code != 1 {
			logx.Errorf("图片上传失败: %v", err)
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("图片上传失败"))
			return
		}

		// Step 2: 调用诊断服务提交诊断任务（异步）
		diagCtx, diagCancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer diagCancel()

		diagResp, err := ctx.DiagnosisClient.SubmitDiagnosis(diagCtx, &diagnosispb.SubmitDiagnosisReq{
			UserId: int64(userID),
			ImagePath: uploadResp.ImagePath,
			ImageUrl:  uploadResp.ImageUrl,
		})
		if err != nil {
			logx.Errorf("提交诊断任务失败: %v", err)
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("诊断提交失败"))
			return
		}

		if diagResp.Code != 1 {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error(diagResp.Msg))
			return
		}

		// 同步等待诊断结果（轮询任务状态，最多等待60秒）
		resultCtx, resultCancel := context.WithTimeout(r.Context(), 90*time.Second)
		defer resultCancel()

		var result *diagnosispb.DiagnosisResultResp
		for i := 0; i < 30; i++ {
			time.Sleep(2 * time.Second)

			// 查询任务状态
			statusResp, err := ctx.DiagnosisClient.GetTaskStatus(resultCtx, &diagnosispb.GetTaskStatusReq{
				TaskId: diagResp.TaskId,
			})
			if err != nil {
				continue
			}

			if statusResp.Status == 2 { // 已完成
				// 查询诊断结果
				listResp, err := ctx.DiagnosisClient.ListDiagnosisResults(resultCtx, &diagnosispb.ListResultsReq{
					UserId: int64(userID),
					Page:     1,
					PageSize: 1,
				})
				if err == nil && len(listResp.Records) > 0 {
					result = listResp.Records[0]
				}
				break
			}

			if statusResp.Status == 3 { // 失败
				httpx.OkJsonCtx(r.Context(), w, pkg.Error("诊断失败: "+statusResp.ErrorMessage))
				return
			}
		}

		if result == nil {
			// 超时但任务仍在处理中，返回任务编号供前端轮询
			httpx.OkJsonCtx(r.Context(), w, pkg.Success(map[string]interface{}{
				"task_no":    diagResp.TaskNo,
				"status":     "processing",
				"message":    "诊断任务正在处理中，请稍后查看结果",
				"image_path": uploadResp.ImagePath,
				"image_url":  uploadResp.ImageUrl,
			}))
		return
	}

	httpx.OkJsonCtx(r.Context(), w, pkg.Success(toDiagnosisResultVO(result)))
}
}

// ListResultsHandler 获取检测结果列表
func ListResultsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		rpcCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// 解析筛选条件
		req := &diagnosispb.ListResultsReq{
			UserId: int64(userID),
			Page:     int32(parseInt(r.URL.Query().Get("page"), 1)),
			PageSize: int32(parseInt(r.URL.Query().Get("pageSize"), 10)),
			Level:    r.URL.Query().Get("level"),
		}

		startDate := r.URL.Query().Get("startDate")
		if startDate != "" {
			req.StartDate = startDate
		}
		endDate := r.URL.Query().Get("endDate")
		if endDate != "" {
			req.EndDate = endDate
		}
		isGlandFaceStr := r.URL.Query().Get("isGlandFace")
		if isGlandFaceStr != "" {
			req.HasFilter = true
			req.IsGlandFace = parseBool(isGlandFaceStr)
		}
		if req.Level != "" {
			req.HasFilter = true
		}

		resp, err := ctx.DiagnosisClient.ListDiagnosisResults(rpcCtx, req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("查询失败"))
			return
		}

	httpx.OkJsonCtx(r.Context(), w, pkg.Success(pkg.PageResult{
		Total:   resp.Total,
		Records: toDiagnosisResultVOs(resp.Records),
	}))
	}
}

// DeleteResultHandler 删除检测结果
func DeleteResultHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		idStr := r.PathValue("id")
		if idStr == "" {
			// go-zero pathvar fallback
			vars := pathvar.Vars(r)
			idStr = vars["id"]
		}
		if idStr == "" {
			// 兼容旧路径格式
			idStr = r.URL.Query().Get("id")
		}

		resultID := parseInt(idStr, 0)
		if resultID == 0 {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("无效的ID"))
			return
		}

		rpcCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		resp, err := ctx.DiagnosisClient.DeleteDiagnosisResult(rpcCtx, &diagnosispb.DeleteResultReq{
			ResultId: int64(resultID),
			UserId: int64(userID),
		})
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("删除失败"))
			return
		}

		if resp.Code == 1 {
			httpx.OkJsonCtx(r.Context(), w, pkg.SuccessMsg(resp.Msg))
		} else {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error(resp.Msg))
		}
	}
}

// BatchDeleteResultsHandler 批量删除
func BatchDeleteResultsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		var ids []int64
		if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("参数解析失败"))
			return
		}

		rpcCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		resp, err := ctx.DiagnosisClient.BatchDeleteResults(rpcCtx, &diagnosispb.BatchDeleteReq{
			Ids:    ids,
			UserId: int64(userID),
		})
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("批量删除失败"))
			return
		}

		if resp.Code == 1 {
			httpx.OkJsonCtx(r.Context(), w, pkg.SuccessMsg(resp.Msg))
		} else {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error(resp.Msg))
		}
	}
}

// UpdateResultHandler 更新检测结果
func UpdateResultHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		var req struct {
			ID                      int64  `json:"id"`
			Level                   string `json:"level"`
			VisualizationDescription string `json:"visualizationDescription"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("参数解析失败"))
			return
		}

		rpcCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		resp, err := ctx.DiagnosisClient.UpdateDiagnosisResult(rpcCtx, &diagnosispb.UpdateResultReq{
			ResultId:                req.ID,
			UserId: int64(userID),
			Level:                   req.Level,
			VisualizationDescription: req.VisualizationDescription,
		})
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("更新失败"))
			return
		}

		if resp.Code == 1 {
			httpx.OkJsonCtx(r.Context(), w, pkg.SuccessMsg(resp.Msg))
		} else {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error(resp.Msg))
		}
	}
}

// ExportExcelHandler 导出Excel
func ExportExcelHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		rpcCtx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		req := &reportpb.ExportReq{
			UserId: int64(userID),
		}

		startDate := r.URL.Query().Get("startDate")
		if startDate != "" {
			req.StartDate = startDate
		}
		endDate := r.URL.Query().Get("endDate")
		if endDate != "" {
			req.EndDate = endDate
		}
		level := r.URL.Query().Get("level")
		if level != "" {
			req.HasFilter = true
			req.Level = level
		}
		isGlandFaceStr := r.URL.Query().Get("isGlandFace")
		if isGlandFaceStr != "" {
			req.HasFilter = true
			req.IsGlandFace = parseBool(isGlandFaceStr)
		}

		resp, err := ctx.ReportClient.ExportExcel(rpcCtx, req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("导出失败"))
			return
		}

		if resp.Code != 1 {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error(resp.Msg))
			return
		}

		w.Header().Set("Content-Type", resp.ContentType)
		w.Header().Set("Content-Disposition", "attachment; filename="+resp.Filename)
		w.Write(resp.FileData)
	}
}

// ExportPDFHandler 导出PDF
func ExportPDFHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		rpcCtx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		req := &reportpb.ExportReq{
			UserId: int64(userID),
		}

		startDate := r.URL.Query().Get("startDate")
		if startDate != "" {
			req.StartDate = startDate
		}
		endDate := r.URL.Query().Get("endDate")
		if endDate != "" {
			req.EndDate = endDate
		}
		level := r.URL.Query().Get("level")
		if level != "" {
			req.HasFilter = true
			req.Level = level
		}
		isGlandFaceStr := r.URL.Query().Get("isGlandFace")
		if isGlandFaceStr != "" {
			req.HasFilter = true
			req.IsGlandFace = parseBool(isGlandFaceStr)
		}

		resp, err := ctx.ReportClient.ExportPDF(rpcCtx, req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("导出失败"))
			return
		}

		if resp.Code != 1 {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error(resp.Msg))
			return
		}

		w.Header().Set("Content-Type", resp.ContentType)
		w.Header().Set("Content-Disposition", "attachment; filename="+resp.Filename)
		w.Write(resp.FileData)
	}
}

// StatisticsOverviewHandler 统计概览
func StatisticsOverviewHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		rpcCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		req := &diagnosispb.StatisticsReq{
			UserId: int64(userID),
		}
		if sd := r.URL.Query().Get("startDate"); sd != "" {
			req.StartDate = sd
		}
		if ed := r.URL.Query().Get("endDate"); ed != "" {
			req.EndDate = ed
		}

		resp, err := ctx.DiagnosisClient.GetStatistics(rpcCtx, req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("获取统计失败"))
			return
		}

		// 组装返回数据
		var trendData, levelData interface{}
		json.Unmarshal([]byte(resp.TrendData), &trendData)
		json.Unmarshal([]byte(resp.LevelData), &levelData)

		result := map[string]interface{}{
			"totalTests":        resp.TotalTests,
			"glandFaceCount":    resp.GlandFaceCount,
			"nonGlandFaceCount": resp.NonGlandFaceCount,
			"trendData":         trendData,
			"levelData":         levelData,
		}

		httpx.OkJsonCtx(r.Context(), w, pkg.Success(result))
	}
}

// StatisticsDetailHandler 统计详情
func StatisticsDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 复用overview逻辑
		StatisticsOverviewHandler(ctx)(w, r)
	}
}
