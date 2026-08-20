package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"faceTest/backend/common/model"
	"faceTest/backend/diagnosis/internal/svc"
	diagnosispb "faceTest/backend/proto/diagnosis"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// GetDiagnosisResultLogic 获取诊断结果
type GetDiagnosisResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDiagnosisResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiagnosisResultLogic {
	return &GetDiagnosisResultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDiagnosisResultLogic) GetDiagnosisResult(req *diagnosispb.GetResultReq) (*diagnosispb.DiagnosisResultResp, error) {
	var result model.DiagnosisResult
	err := l.svcCtx.DB.Where("id = ? AND user_id = ?", req.ResultId, req.UserId).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosispb.DiagnosisResultResp{}, nil
		}
		return nil, err
	}

	return convertResultToResp(&result), nil
}

// GetTaskStatusLogic 获取任务状态
type GetTaskStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskStatusLogic {
	return &GetTaskStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskStatusLogic) GetTaskStatus(req *diagnosispb.GetTaskStatusReq) (*diagnosispb.GetTaskStatusResp, error) {
	var task model.DiagnosisTask
	err := l.svcCtx.DB.Where("id = ?", req.TaskId).First(&task).Error
	if err != nil {
		return nil, err
	}

	return &diagnosispb.GetTaskStatusResp{
		TaskId:       int64(task.ID),
		TaskNo:       task.TaskNo,
		Status:       int32(task.Status),
		RetryCount:   int32(task.RetryCount),
		ErrorMessage: task.ErrorMessage,
	}, nil
}

// ListDiagnosisResultsLogic 获取诊断结果列表
type ListDiagnosisResultsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDiagnosisResultsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDiagnosisResultsLogic {
	return &ListDiagnosisResultsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListDiagnosisResultsLogic) ListDiagnosisResults(req *diagnosispb.ListResultsReq) (*diagnosispb.ListResultsResp, error) {
	page := int(req.Page)
	if page == 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize == 0 {
		pageSize = 10
	}

	query := l.svcCtx.DB.Model(&model.DiagnosisResult{}).Where("user_id = ?", req.UserId)

	// 条件筛选
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

	var total int64
	query.Count(&total)

	var results []model.DiagnosisResult
	query.Order("test_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&results)

	var records []*diagnosispb.DiagnosisResultResp
	for _, r := range results {
		records = append(records, convertResultToResp(&r))
	}

	return &diagnosispb.ListResultsResp{
		Total:   total,
		Records: records,
	}, nil
}

// DeleteDiagnosisResultLogic 删除诊断结果
type DeleteDiagnosisResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDiagnosisResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDiagnosisResultLogic {
	return &DeleteDiagnosisResultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteDiagnosisResultLogic) DeleteDiagnosisResult(req *diagnosispb.DeleteResultReq) (*diagnosispb.DeleteResultResp, error) {
	result := l.svcCtx.DB.Where("id = ? AND user_id = ?", req.ResultId, req.UserId).Delete(&model.DiagnosisResult{})
	if result.Error != nil {
		return &diagnosispb.DeleteResultResp{Code: 0, Msg: "删除失败"}, nil
	}
	if result.RowsAffected == 0 {
		return &diagnosispb.DeleteResultResp{Code: 0, Msg: "记录不存在或无权限"}, nil
	}
	return &diagnosispb.DeleteResultResp{Code: 1, Msg: "删除成功"}, nil
}

// BatchDeleteResultsLogic 批量删除
type BatchDeleteResultsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteResultsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteResultsLogic {
	return &BatchDeleteResultsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchDeleteResultsLogic) BatchDeleteResults(req *diagnosispb.BatchDeleteReq) (*diagnosispb.BatchDeleteResp, error) {
	if len(req.Ids) == 0 {
		return &diagnosispb.BatchDeleteResp{Code: 0, Msg: "请选择要删除的记录"}, nil
	}

	result := l.svcCtx.DB.Where("id IN ? AND user_id = ?", req.Ids, req.UserId).Delete(&model.DiagnosisResult{})
	if result.Error != nil {
		return &diagnosispb.BatchDeleteResp{Code: 0, Msg: "批量删除失败"}, nil
	}

	return &diagnosispb.BatchDeleteResp{
		Code:         1,
		Msg:          fmt.Sprintf("成功删除 %d 条记录", result.RowsAffected),
		DeletedCount: int32(result.RowsAffected),
	}, nil
}

// UpdateDiagnosisResultLogic 更新诊断结果
type UpdateDiagnosisResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDiagnosisResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDiagnosisResultLogic {
	return &UpdateDiagnosisResultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDiagnosisResultLogic) UpdateDiagnosisResult(req *diagnosispb.UpdateResultReq) (*diagnosispb.UpdateResultResp, error) {
	// 使用乐观锁更新
	var result model.DiagnosisResult
	if err := l.svcCtx.DB.Where("id = ? AND user_id = ?", req.ResultId, req.UserId).First(&result).Error; err != nil {
		return &diagnosispb.UpdateResultResp{Code: 0, Msg: "记录不存在或无权限"}, nil
	}

	updates := map[string]interface{}{
		"version": result.Version + 1,
	}
	if req.Level != "" {
		updates["level"] = req.Level
	}
	if req.VisualizationDescription != "" {
		updates["visualization_description"] = req.VisualizationDescription
	}

	dbResult := l.svcCtx.DB.Model(&model.DiagnosisResult{}).
		Where("id = ? AND version = ?", req.ResultId, result.Version).
		Updates(updates)

	if dbResult.RowsAffected == 0 {
		// 乐观锁冲突，直接更新
		l.svcCtx.DB.Model(&model.DiagnosisResult{}).
			Where("id = ?", req.ResultId).
			Updates(map[string]interface{}{
				"level":                   req.Level,
				"visualization_description": req.VisualizationDescription,
			})
	}

	return &diagnosispb.UpdateResultResp{Code: 1, Msg: "更新成功"}, nil
}

// GetStatisticsLogic 获取统计数据
type GetStatisticsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStatisticsLogic {
	return &GetStatisticsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetStatisticsLogic) GetStatistics(req *diagnosispb.StatisticsReq) (*diagnosispb.StatisticsResp, error) {
	query := l.svcCtx.DB.Model(&model.DiagnosisResult{}).Where("user_id = ?", req.UserId)

	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			query = query.Where("test_time >= ?", t)
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			query = query.Where("test_time <= ?", t.Add(24*time.Hour))
		}
	}

	var totalTests, glandFaceCount int64

	query.Count(&totalTests)
	query.Where("is_gland_face = ?", true).Count(&glandFaceCount)

	// 趋势数据
	type TrendItem struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var trendData []TrendItem
	l.svcCtx.DB.Model(&model.DiagnosisResult{}).
		Select("DATE(test_time) as date, COUNT(*) as count").
		Where("user_id = ?", req.UserId).
		Group("DATE(test_time)").
		Order("date DESC").
		Limit(30).
		Scan(&trendData)
	trendJSON, _ := json.Marshal(trendData)

	// 等级分布
	type LevelItem struct {
		Level string `json:"level"`
		Count int64  `json:"count"`
	}
	var levelData []LevelItem
	l.svcCtx.DB.Model(&model.DiagnosisResult{}).
		Select("level, COUNT(*) as count").
		Where("user_id = ?", req.UserId).
		Group("level").
		Scan(&levelData)
	levelJSON, _ := json.Marshal(levelData)

	return &diagnosispb.StatisticsResp{
		TotalTests:        totalTests,
		GlandFaceCount:    glandFaceCount,
		NonGlandFaceCount: totalTests - glandFaceCount,
		TrendData:         string(trendJSON),
		LevelData:         string(levelJSON),
	}, nil
}

// convertResultToResp 转换结果模型为gRPC响应
func convertResultToResp(r *model.DiagnosisResult) *diagnosispb.DiagnosisResultResp {
	return &diagnosispb.DiagnosisResultResp{
		Id:                      int64(r.ID),
		TaskId:                  int64(r.TaskID),
		UserId:                  int64(r.UserID),
		ImagePath:               r.ImagePath,
		ImageUrl:                r.ImageURL,
		IsGlandFace:             r.IsGlandFace,
		Level:                   r.Level,
		VisualizationDescription: r.VisualizationDescription,
		ReferenceCases:          r.ReferenceCases,
		SkippedImages:           r.SkippedImages,
		TestTime:                r.TestTime.Format("2006-01-02 15:04:05"),
	}
}
