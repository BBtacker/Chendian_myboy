package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"faceTest/backend/common/model"
	"faceTest/backend/common/pkg"
	"faceTest/backend/diagnosis/internal/svc"
	diagnosispb "faceTest/backend/proto/diagnosis"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// SubmitDiagnosisLogic 提交诊断任务逻辑
type SubmitDiagnosisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitDiagnosisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitDiagnosisLogic {
	return &SubmitDiagnosisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SubmitDiagnosis 提交诊断任务（使用Outbox模式保证原子性）
func (l *SubmitDiagnosisLogic) SubmitDiagnosis(req *diagnosispb.SubmitDiagnosisReq) (*diagnosispb.SubmitDiagnosisResp, error) {
	// 生成唯一任务编号（幂等键）
	taskNo := fmt.Sprintf("DIAG-%s-%d", uuid.New().String()[:8], time.Now().Unix())

	// 限流检查：每用户每分钟最多10次诊断
	rateKey := fmt.Sprintf("rate_limit:diagnosis:user_%d", req.UserId)
	allowed, err := l.svcCtx.RateLimiter.Allow(l.ctx, rateKey, 60, 10)
	if err == nil && !allowed {
		return &diagnosispb.SubmitDiagnosisResp{
			Code: 0,
			Msg:  "诊断请求过于频繁，请稍后再试",
		}, nil
	}

	// 使用事务 + Outbox模式：保证任务创建和消息发送的原子性
	var task model.DiagnosisTask
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 转换类型（proto: int64/int32 -> model: uint64/*int8）
		gender := int8(req.Gender)

		// 1. 创建诊断任务
		task = model.DiagnosisTask{
			TaskNo:     taskNo,
			UserID:     uint64(req.UserId),
			ImagePath:  req.ImagePath,
			ImageURL:   req.ImageUrl,
			Status:     model.TaskStatusPending,
			Age:        &req.Age,
			Gender:     &gender,
			MaxRetry:   l.svcCtx.Config.Diagnosis.MaxRetry,
		}

		if err := tx.Create(&task).Error; err != nil {
			return fmt.Errorf("创建诊断任务失败: %w", err)
		}

		// 2. 创建Outbox消息（在同一事务中）
		msgPayload, _ := json.Marshal(pkg.DiagnosisMessage{
			TaskID:        task.ID,
			TaskNo:        taskNo,
			UserID:        uint64(req.UserId),
			ImagePath:     req.ImagePath,
			ImageURL:      req.ImageUrl,
			Age:           req.Age,
			Gender:        gender,
			IdempotentKey: taskNo, // 唯一幂等键
			Timestamp:     time.Now().Unix(),
		})

		outboxMsg := model.OutboxMessage{
			AggregateID:   fmt.Sprintf("%d", task.ID),
			AggregateType: "diagnosis_task",
			MessageType:   "diagnosis.submit",
			Payload:       string(msgPayload),
			Status:        model.OutboxStatusPending,
			MaxRetry:      l.svcCtx.Config.Diagnosis.MaxRetry,
		}

		if err := tx.Create(&outboxMsg).Error; err != nil {
			return fmt.Errorf("创建Outbox消息失败: %w", err)
		}

		l.Infof("诊断任务已创建: taskID=%d, taskNo=%s", task.ID, taskNo)
		return nil
	})

	if err != nil {
		l.Errorf("提交诊断任务失败: %v", err)
		return &diagnosispb.SubmitDiagnosisResp{
			Code: 0,
			Msg:  "提交诊断任务失败: " + err.Error(),
		}, nil
	}

	return &diagnosispb.SubmitDiagnosisResp{
		Code:   1,
		Msg:    "诊断任务已提交，正在处理中",
		TaskId: int64(task.ID), // 返回真实任务ID，供网关轮询 GetTaskStatus
		TaskNo: taskNo,
	}, nil
}
