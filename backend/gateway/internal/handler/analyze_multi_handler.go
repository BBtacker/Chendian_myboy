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
	uploadpb "faceTest/backend/proto/upload"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// AnalyzeFaceMultiHandler 多图面容分析（合并诊断）
// 接收多张图片，逐张上传后提交一次合并诊断任务；后端对多张图做特征平均 +
// 人脸门控 + Milvus 检索 + 一次 DeepSeek，最终返回一份综合诊断结果。
func AnalyzeFaceMultiHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		// 限制总体积（最多 9 张 * 10MB）
		r.Body = http.MaxBytesReader(w, r.Body, 100<<20)
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("请求解析失败，请确认上传的是图片"))
			return
		}
		files := r.MultipartForm.File["images"]
		if len(files) == 0 {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("请至少上传一张图片"))
			return
		}
		if len(files) > 9 {
			files = files[:9]
		}

		// 逐张上传到上传服务
		var paths, urls []string
		for _, fh := range files {
			file, err := fh.Open()
			if err != nil {
				continue
			}
			data, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				continue
			}

			upCtx, upCancel := context.WithTimeout(r.Context(), 30*time.Second)
			upResp, err := ctx.UploadClient.UploadImage(upCtx, &uploadpb.UploadImageReq{
				FileData:    data,
				Filename:    fh.Filename,
				ContentType: fh.Header.Get("Content-Type"),
			})
			upCancel()
			if err != nil || upResp.Code != 1 {
				logx.Errorf("图片上传失败: %v", err)
				continue
			}
			paths = append(paths, upResp.ImagePath)
			urls = append(urls, upResp.ImageUrl)
		}
		if len(paths) == 0 {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("图片上传失败，请重试"))
			return
		}

		// 提交合并诊断任务（image_path / image_url 以 JSON 数组形式传递）
		jsonPaths, _ := json.Marshal(paths)
		jsonUrls, _ := json.Marshal(urls)
		diagCtx, diagCancel := context.WithTimeout(r.Context(), 10*time.Second)
		diagResp, err := ctx.DiagnosisClient.SubmitDiagnosis(diagCtx, &diagnosispb.SubmitDiagnosisReq{
			UserId:    int64(userID),
			ImagePath: string(jsonPaths),
			ImageUrl:  string(jsonUrls),
		})
		diagCancel()
		if err != nil {
			logx.Errorf("提交诊断任务失败: %v", err)
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("诊断提交失败"))
			return
		}
		if diagResp.Code != 1 {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error(diagResp.Msg))
			return
		}

		// 轮询任务状态（最多约 90 秒）
		resultCtx, resultCancel := context.WithTimeout(r.Context(), 120*time.Second)
		defer resultCancel()

		var result *diagnosispb.DiagnosisResultResp
		for i := 0; i < 45; i++ {
			time.Sleep(2 * time.Second)

			statusResp, err := ctx.DiagnosisClient.GetTaskStatus(resultCtx, &diagnosispb.GetTaskStatusReq{
				TaskId: diagResp.TaskId,
			})
			if err != nil {
				continue
			}

			if statusResp.Status == 2 { // 已完成
				listResp, err := ctx.DiagnosisClient.ListDiagnosisResults(resultCtx, &diagnosispb.ListResultsReq{
					UserId:   int64(userID),
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
				"task_no": diagResp.TaskNo,
				"status":  "processing",
				"message": "诊断任务正在处理中，请稍后查看结果",
			}))
			return
		}

		httpx.OkJsonCtx(r.Context(), w, pkg.Success(toDiagnosisResultVO(result)))
	}
}
