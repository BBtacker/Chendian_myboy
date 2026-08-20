package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"time"

	"faceTest/backend/common/model"
	"faceTest/backend/common/pkg"
	"faceTest/backend/diagnosis/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// DiagnosisConsumer 诊断消息消费者
type DiagnosisConsumer struct {
	svcCtx *svc.ServiceContext
}

// NewDiagnosisConsumer 创建诊断消费者
func NewDiagnosisConsumer(svcCtx *svc.ServiceContext) *DiagnosisConsumer {
	return &DiagnosisConsumer{svcCtx: svcCtx}
}

// Start 启动消费者（包含Outbox中继和RabbitMQ消费）
func (dc *DiagnosisConsumer) Start() {
	// 1. 启动Outbox中继（定时扫描待发送消息）
	go dc.startOutboxRelay()

	// 2. 启动RabbitMQ消费者
	go dc.startRabbitMQConsumer()

	// 3. 启动死信队列消费者（处理失败消息）
	go dc.startDeadLetterConsumer()

	logx.Info("诊断消费者已启动（Outbox中继 + RabbitMQ消费 + 死信队列）")
}

// startOutboxRelay Outbox中继：定时扫描待发送消息并发布到RabbitMQ
func (dc *DiagnosisConsumer) startOutboxRelay() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		dc.relayOutboxMessages()
	}
}

// relayOutboxMessages 中继Outbox消息
func (dc *DiagnosisConsumer) relayOutboxMessages() {
	// RabbitMQ 未连接时直接跳过，避免空指针 panic 导致整个诊断服务崩溃；
	// 待发送消息仍留在 outbox_message 表中，下次轮询会重试。
	if dc.svcCtx.RabbitMQ == nil {
		return
	}

	var messages []model.OutboxMessage
	// 查询待发送或需要重试的消息
	dc.svcCtx.DB.Where("status = ? OR (status = ? AND next_retry_time <= ?)",
		model.OutboxStatusPending, model.OutboxStatusSent, time.Now()).
		Limit(100).
		Find(&messages)

	for _, msg := range messages {
		// 检查是否超过最大重试次数
		if msg.RetryCount >= msg.MaxRetry {
			dc.svcCtx.DB.Model(&msg).Update("status", model.OutboxStatusFailed)
			logx.Errorf("Outbox消息超过最大重试次数: id=%d, aggregateID=%s", msg.ID, msg.AggregateID)
			continue
		}

		// 解析消息内容
		var diagMsg pkg.DiagnosisMessage
		if err := json.Unmarshal([]byte(msg.Payload), &diagMsg); err != nil {
			logx.Errorf("解析Outbox消息失败: id=%d, err=%v", msg.ID, err)
			continue
		}

		// 发布到RabbitMQ
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := dc.svcCtx.RabbitMQ.PublishMessage(ctx, "", dc.svcCtx.Config.RabbitMQ.Queue, diagMsg)
		cancel()

		if err != nil {
			// 发布失败，更新重试次数和下次重试时间
			nextRetry := time.Now().Add(time.Duration(1<<uint(msg.RetryCount)) * time.Second)
			dc.svcCtx.DB.Model(&msg).Updates(map[string]interface{}{
				"retry_count":     msg.RetryCount + 1,
				"status":          model.OutboxStatusSent,
				"next_retry_time": nextRetry,
			})
			logx.Infof("Outbox消息发布失败，将重试: id=%d, retry=%d, nextRetry=%v", msg.ID, msg.RetryCount+1, nextRetry)
		} else {
			// 发布成功，标记为已确认
			dc.svcCtx.DB.Model(&msg).Update("status", model.OutboxStatusConfirmed)
			logx.Infof("Outbox消息已发布: id=%d, taskNo=%s", msg.ID, diagMsg.TaskNo)
		}
	}
}

// startRabbitMQConsumer 启动RabbitMQ消费者
func (dc *DiagnosisConsumer) startRabbitMQConsumer() {
	if dc.svcCtx.RabbitMQ == nil {
		logx.Error("RabbitMQ未连接，消费者无法启动")
		return
	}

	queueName := dc.svcCtx.Config.RabbitMQ.Queue
	err := dc.svcCtx.RabbitMQ.ConsumeMessages(queueName, dc.handleDiagnosisMessage)
	if err != nil {
		logx.Errorf("启动RabbitMQ消费者失败: %v", err)
	}
}

// handleDiagnosisMessage 处理诊断消息
func (dc *DiagnosisConsumer) handleDiagnosisMessage(body []byte) error {
	var msg pkg.DiagnosisMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("消息反序列化失败: %w", err)
	}

	logx.Infof("开始处理诊断消息: taskNo=%s, taskID=%d", msg.TaskNo, msg.TaskID)

	// 1. 幂等性检查：检查任务是否已经处理过
	var task model.DiagnosisTask
	if err := dc.svcCtx.DB.Where("task_no = ?", msg.IdempotentKey).First(&task).Error; err != nil {
		return fmt.Errorf("查询任务失败: %w", err)
	}

	if task.Status == model.TaskStatusCompleted {
		logx.Infof("任务已处理完成，跳过（幂等）: taskNo=%s", msg.TaskNo)
		return nil
	}

	if task.Status == model.TaskStatusProcessing && task.RetryCount > 0 {
		// 可能是重复消息，检查是否已有结果
		var count int64
		dc.svcCtx.DB.Model(&model.DiagnosisResult{}).Where("task_id = ?", task.ID).Count(&count)
		if count > 0 {
			logx.Infof("任务已有结果，跳过（幂等）: taskNo=%s", msg.TaskNo)
			return nil
		}
	}

	// 2. 更新任务状态为处理中（使用乐观锁）
	result := dc.svcCtx.DB.Model(&model.DiagnosisTask{}).
		Where("id = ? AND version = ?", task.ID, task.Version).
		Updates(map[string]interface{}{
			"status":  model.TaskStatusProcessing,
			"version": task.Version + 1,
		})

	if result.RowsAffected == 0 {
		// 乐观锁冲突，说明有其他消费者正在处理
		logx.Infof("任务正在被其他消费者处理，跳过: taskNo=%s", msg.TaskNo)
		return nil
	}

	// 3. 执行诊断流水线
	err := dc.executeDiagnosisPipeline(task)
	if err != nil {
		logx.Errorf("诊断流水线失败: taskNo=%s, err=%v", msg.TaskNo, err)

		// 更新任务状态为失败
		dc.svcCtx.DB.Model(&model.DiagnosisTask{}).
			Where("id = ?", task.ID).
			Updates(map[string]interface{}{
				"status":         model.TaskStatusFailed,
				"error_message":  err.Error(),
				"retry_count":    task.RetryCount + 1,
			})

		// 如果未超过最大重试次数，重新入队（指数退避）
		if task.RetryCount+1 < task.MaxRetry {
			return fmt.Errorf("诊断失败，将重试: %w", err) // 返回错误触发Nack进入死信队列
		}
		return nil
	}

	logx.Infof("诊断任务处理完成: taskNo=%s", msg.TaskNo)
	return nil
}

// perImageFeature 单张图片的特征提取结果
type perImageFeature struct {
	path        string
	url         string
	vector      []float32
	description string
	measurements map[string]float64
	faceDetected bool
}

// featureAggregate 多张图片特征聚合结果
type featureAggregate struct {
	used          []perImageFeature // 检测到人脸、参与综合判断的图片
	skippedPaths  []string          // 未检测到人脸、被自动剔除的图片路径
	avgVector     []float32         // 平均后的特征向量（L2 归一化）
	avgDescription string           // 基于平均向量的 B3 统计描述
	avgMeasurements map[string]float64 // 跨图平均的几何测量
	perImageNotes []string          // 每张图片的处理说明（用于一致性说明）
}

// executeDiagnosisPipeline 执行诊断流水线（支持单图与多图合并）
// 多图：逐张提取 EfficientNet-B3 特征 + MediaPipe 几何测量，人脸门控剔除无脸图，
//       跨图平均特征向量与各几何测量，Milvus 检索平均向量，一次 DeepSeek 综合判断。
// 单图：image_path 为普通路径，退化为原逻辑。
func (dc *DiagnosisConsumer) executeDiagnosisPipeline(task model.DiagnosisTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	// 解析图片列表（多图：image_path 存 JSON 数组；单图：普通路径）
	imgPaths, imgUrls := parseImageList(task.ImagePath, task.ImageURL)
	if len(imgPaths) == 0 {
		return fmt.Errorf("任务没有关联图片")
	}

	// Step 1: 逐张特征提取 + 人脸门控 + 跨图平均
	logx.Infof("[Step 1] 特征提取中(共%d张): taskNo=%s", len(imgPaths), task.TaskNo)
	agg, err := dc.extractFeaturesMulti(ctx, imgPaths, imgUrls)
	if err != nil {
		return fmt.Errorf("特征提取失败: %w", err)
	}
	if len(agg.used) == 0 {
		return fmt.Errorf("所有照片均未检测到人脸，无法生成诊断，请上传包含清晰面部的照片")
	}

	featureVector := agg.avgVector
	featureDescription := agg.avgDescription
	measurements := agg.avgMeasurements

	// Step 2: Milvus RAG检索相似病例（基于平均特征向量）
	logx.Infof("[Step 2] RAG检索中: taskNo=%s", task.TaskNo)
	var referenceCases string
	if dc.svcCtx.MilvusClient != nil {
		results, err := dc.svcCtx.MilvusClient.SearchSimilar(ctx, featureVector, 5)
		if err != nil {
			logx.Errorf("Milvus检索失败: %v", err)
			referenceCases = "Milvus检索失败，使用默认参考"
		} else {
			referenceCases = pkg.FormatReferenceCases(results)
		}
	} else {
		referenceCases = "Milvus未连接，暂无相似病例参考"
	}
	time.Sleep(200 * time.Millisecond) // 模拟检索耗时

	// Step 4: DeepSeek API生成诊断报告（综合多图 + 几何测量为主证据）
	logx.Infof("[Step 3] DeepSeek推理中: taskNo=%s", task.TaskNo)
	var aiResult *pkg.DiagnosisAIResult
	if dc.svcCtx.DeepSeekClient != nil {
		multiSummary := agg.consistencyNote(len(imgPaths))
		var err error
		aiResult, err = dc.svcCtx.DeepSeekClient.Diagnose(ctx, featureDescription, referenceCases, measurements, multiSummary)
		if err != nil {
			return fmt.Errorf("DeepSeek推理失败: %w", err)
		}
	} else {
		// DeepSeek未配置，使用模拟结果
		aiResult = &pkg.DiagnosisAIResult{
			IsGlandFace:               true,
			Level:                     "中度",
			VisualizationDescription:  "患儿面部呈典型腺样体面容特征：上唇明显上翘外翻，下颌后缩，鼻唇沟变浅，牙列轻度拥挤，存在张口呼吸特征。建议进一步进行腺样体检查确认。",
		}
	}

	// 归一化 level，防止 DeepSeek 返回变体（如"中重度"/"中度腺样体面容"/"疑似轻度"）
	// 污染统计分组与 Milvus 检索。统一收敛为四个标准值。
	aiResult.Level = normalizeLevel(aiResult.Level)

	// Step 5: 保存诊断结果到MySQL（多图路径以 JSON 数组存储）
	logx.Infof("[Step 4] 保存结果: taskNo=%s", task.TaskNo)
	now := time.Now()
	diagResult := model.DiagnosisResult{
		TaskID:                   task.ID,
		UserID:                   task.UserID,
		ImagePath:                toJSONArray(imgPaths),
		ImageURL:                 toJSONArray(imgUrls),
		IsGlandFace:              aiResult.IsGlandFace,
		Level:                    aiResult.Level,
		VisualizationDescription: aiResult.VisualizationDescription,
		ReferenceCases:           referenceCases,
		SkippedImages:            toJSONArray(agg.skippedPaths),
		TestTime:                 now,
	}

	if err := dc.svcCtx.DB.Create(&diagResult).Error; err != nil {
		return fmt.Errorf("保存诊断结果失败: %w", err)
	}

	// Step 6: 特征向量存入Milvus（RAG知识库积累）
	if dc.svcCtx.MilvusClient != nil {
		logx.Infof("[Step 5] 存入Milvus: taskNo=%s", task.TaskNo)
		_, err := dc.svcCtx.MilvusClient.InsertFeature(ctx, pkg.FeatureRecord{
			ResultID: int64(diagResult.ID),
			UserID:   int64(task.UserID),
			Level:    aiResult.Level,
			Vector:   featureVector,
		})
		if err != nil {
			logx.Errorf("Milvus插入失败（不影响主流程）: %v", err)
		}
	}

	// Step 7: 更新任务状态为完成
	dc.svcCtx.DB.Model(&model.DiagnosisTask{}).
		Where("id = ?", task.ID).
		Updates(map[string]interface{}{
			"status":         model.TaskStatusCompleted,
			"complete_time":  now,
		})

	logx.Infof("诊断流水线完成: taskNo=%s, 图片=%d, 有效=%d, 剔除=%d, level=%s",
		task.TaskNo, len(imgPaths), len(agg.used), len(agg.skippedPaths), aiResult.Level)

	return nil
}

// normalizeLevel 将 DeepSeek 返回的各种 level 表述收敛为四个标准值，
// 避免模型"嘴瓢"产生的非预期字符串进入 GROUP BY level 统计与 Milvus 检索。
// 标准值：轻度 / 中度 / 重度 / 非腺样体面容。
// 优先级：重 > 中 > 轻 > 非；命中"非"且未命中轻重中时归为无面容；其余兜底为"非腺样体面容"。
func normalizeLevel(s string) string {
	t := strings.TrimSpace(s)
	switch {
	case strings.Contains(t, "重度"):
		return "重度"
	case strings.Contains(t, "中度"), strings.Contains(t, "中重"):
		return "中度"
	case strings.Contains(t, "轻度"):
		return "轻度"
	case strings.Contains(t, "非腺样体"), strings.Contains(t, "非腺体"), strings.Contains(t, "无腺样体"), strings.Contains(t, "正常"):
		return "非腺样体面容"
	default:
		// 空串或无法识别的变体：兜底为无面容，避免污染已知三档
		return "非腺样体面容"
	}
}

// extractFeaturesMulti 逐张提取特征并跨图聚合。
// 人脸门控：仅 faceDetected 为 true 的图片参与平均（无脸图被剔除并计入 skippedPaths）。
// 特征服务不可用时对该张回退 Mock（视为有人脸并参与平均）。
func (dc *DiagnosisConsumer) extractFeaturesMulti(ctx context.Context, imgPaths, imgUrls []string) (*featureAggregate, error) {
	agg := &featureAggregate{}
	useMock := dc.svcCtx.FeatureExtractor == nil

	for i, p := range imgPaths {
		full := resolveImagePath(p, dc.svcCtx.Config.UploadBasePath)
		var res *pkg.ExtractResult
		var err error
		if !useMock {
			res, err = dc.svcCtx.FeatureExtractor.Extract(ctx, full)
		}
		if useMock || err != nil || len(res.FeatureVector) == 0 {
			// 该张回退 Mock（视为有人脸，参与平均，避免单张失败拖垮整体）
			res = &pkg.ExtractResult{
				FeatureVector:  pkg.MockFeatureVector(),
				Description:    pkg.MockFeatureDescription(),
				FaceDetected:   true,
				Measurements:   map[string]float64{},
			}
			logx.Errorf("[Step 1] 图片特征提取失败/回退Mock: path=%s, err=%v", p, err)
		}

		url := ""
		if i < len(imgUrls) {
			url = imgUrls[i]
		}
		pi := perImageFeature{
			path:         p,
			url:          url,
			vector:       res.FeatureVector,
			description:  res.Description,
			measurements: res.Measurements,
			faceDetected: res.FaceDetected,
		}

		if !res.FaceDetected {
			agg.skippedPaths = append(agg.skippedPaths, p)
			agg.perImageNotes = append(agg.perImageNotes,
				fmt.Sprintf("图片[%s]：未检测到人脸，已自动剔除", shortName(p)))
			continue
		}
		agg.used = append(agg.used, pi)
		agg.perImageNotes = append(agg.perImageNotes,
			fmt.Sprintf("图片[%s]：已纳入（人脸检测OK）", shortName(p)))
	}

	if len(agg.used) == 0 {
		return agg, nil
	}

	// 平均特征向量（L2 归一化）
	dim := len(agg.used[0].vector)
	sum := make([]float32, dim)
	for _, u := range agg.used {
		for j := 0; j < dim; j++ {
			sum[j] += u.vector[j]
		}
	}
	var norm float64
	for j := 0; j < dim; j++ {
		norm += float64(sum[j]) * float64(sum[j])
	}
	if norm > 0 {
		s := float32(1.0 / math.Sqrt(norm))
		for j := 0; j < dim; j++ {
			sum[j] *= s
		}
	}
	agg.avgVector = sum
	agg.avgDescription = describeB3Vector(sum)

	// 平均几何测量（仅跨已纳入图片）
	measKeys := []string{
		"mouth_aperture_ratio",
		"lower_face_ratio",
		"face_width_height_ratio",
		"nasolabial_angle_deg",
		"facial_convexity_angle_deg",
	}
	summ := map[string]float64{}
	counts := map[string]int{}
	for _, u := range agg.used {
		for _, k := range measKeys {
			if v, ok := u.measurements[k]; ok {
				summ[k] += v
				counts[k]++
			}
		}
	}
	avg := map[string]float64{}
	for _, k := range measKeys {
		if c := counts[k]; c > 0 {
			avg[k] = summ[k] / float64(c)
		}
	}
	agg.avgMeasurements = avg

	return agg, nil
}

// consistencyNote 生成多图一致性说明，供 DeepSeek 综合判断时参考。
func (a *featureAggregate) consistencyNote(total int) string {
	if total <= 1 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("【多图综合说明】本次共提交 %d 张照片，其中 %d 张检测到清晰面部并参与综合判断", total, len(a.used)))
	if len(a.skippedPaths) > 0 {
		sb.WriteString(fmt.Sprintf("，%d 张未检测到人脸已自动排除", len(a.skippedPaths)))
	}
	sb.WriteString("。各图片处理情况：\n")
	for _, n := range a.perImageNotes {
		sb.WriteString("- " + n + "\n")
	}
	sb.WriteString("请结合多张照片的一致性进行综合判断；若各图在关键几何测量或面容特征上存在明显分歧，请明确说明分歧点及最终采信依据。")
	return sb.String()
}

// resolveImagePath 拼接图片完整路径（相对路径 + 上传根目录）
func resolveImagePath(imagePath, base string) string {
	if imagePath == "" || filepath.IsAbs(imagePath) {
		return imagePath
	}
	if base == "" {
		base = "./uploads"
	}
	return filepath.Join(base, imagePath)
}

// parseImageList 解析 image_path / image_url：JSON 数组或单值路径都支持。
func parseImageList(paths, urls string) ([]string, []string) {
	return parseJSONArrayOrSingle(paths), parseJSONArrayOrSingle(urls)
}

func parseJSONArrayOrSingle(s string) []string {
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

// toJSONArray 将字符串切片序列化为 JSON 数组字符串。
func toJSONArray(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	b, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// describeB3Vector 基于平均特征向量生成确定性统计描述（与 Python 端 generate_description 对应）。
func describeB3Vector(v []float32) string {
	if len(v) == 0 {
		return "{}"
	}
	keys := []string{
		"上唇上翘程度",
		"鼻唇沟变浅程度",
		"下颌后缩程度",
		"牙列不齐程度",
		"张口呼吸特征",
		"面部比例失调程度",
	}
	segLen := len(v) / 6
	desc := make(map[string]float64, 6)
	for i := 0; i < 6; i++ {
		start := i * segLen
		end := (i + 1) * segLen
		if i == 5 {
			end = len(v)
		}
		var s float64
		for j := start; j < end; j++ {
			s += math.Abs(float64(v[j]))
		}
		n := float64(end - start)
		if n == 0 {
			n = 1
		}
		val := s / n * 3.0
		if val > 1 {
			val = 1
		}
		if val < 0 {
			val = 0
		}
		desc[keys[i]] = math.Round(val*100) / 100
	}
	b, _ := json.Marshal(desc)
	return string(b)
}

// shortName 取路径最后一段作为展示名。
func shortName(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return "?"
	}
	if idx := strings.LastIndexAny(p, "/\\"); idx >= 0 {
		return p[idx+1:]
	}
	return p
}

// startDeadLetterConsumer 启动死信队列消费者
func (dc *DiagnosisConsumer) startDeadLetterConsumer() {
	if dc.svcCtx.RabbitMQ == nil {
		return
	}

	dlxQueue := dc.svcCtx.Config.RabbitMQ.DLXQueue
	if dlxQueue == "" {
		dlxQueue = "diagnosis.dlx.queue"
	}

	// 创建新的channel用于消费死信队列
	ch, err := dc.svcCtx.RabbitMQ.NewChannel()
	if err != nil {
		logx.Errorf("创建死信队列消费者channel失败: %v", err)
		return
	}

	msgs, err := ch.Consume(dlxQueue, "", false, false, false, false, nil)
	if err != nil {
		logx.Errorf("消费死信队列失败: %v", err)
		return
	}

	go func() {
		for msg := range msgs {
			logx.Errorf("收到死信消息: messageId=%s, body=%s", msg.MessageId, string(msg.Body))

			// 解析消息，更新任务状态为失败
			var diagMsg pkg.DiagnosisMessage
			if err := json.Unmarshal(msg.Body, &diagMsg); err == nil {
				dc.svcCtx.DB.Model(&model.DiagnosisTask{}).
					Where("task_no = ?", diagMsg.IdempotentKey).
					Updates(map[string]interface{}{
						"status":        model.TaskStatusFailed,
						"error_message": "诊断处理失败，已达到最大重试次数",
					})
			}

			_ = msg.Ack(false)
		}
	}()

	logx.Info("死信队列消费者已启动")
}

// Stop 停止消费者
func (dc *DiagnosisConsumer) Stop() {
	if dc.svcCtx.RabbitMQ != nil {
		dc.svcCtx.RabbitMQ.Close()
	}
}

// 避免unused import
var _ = gorm.ErrRecordNotFound
