package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	milvusCollectionName = "diagnosis_features"
	milvusDimension      = 1536 // EfficientNet-B3 特征维度
)

// MilvusClient Milvus向量数据库客户端
type MilvusClient struct {
	client client.Client
}

// NewMilvusClient 创建Milvus连接
func NewMilvusClient(addr string) (*MilvusClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := client.NewClient(ctx, client.Config{
		Address:  addr,
		Username: "",
		Password: "",
	})
	if err != nil {
		return nil, fmt.Errorf("Milvus连接失败: %w", err)
	}

	mc := &MilvusClient{client: c}

	// 初始化集合
	if err := mc.initCollection(ctx); err != nil {
		logx.Errorf("初始化Milvus集合失败: %v", err)
	}

	logx.Info("Milvus连接成功")
	return mc, nil
}

// initCollection 初始化向量集合
func (mc *MilvusClient) initCollection(ctx context.Context) error {
	has, err := mc.client.HasCollection(ctx, milvusCollectionName)
	if err != nil {
		return err
	}

	if !has {
		// 创建集合schema
		schema := &entity.Schema{
			CollectionName: milvusCollectionName,
			Description:    "腺样体面容诊断特征向量",
			Fields: []*entity.Field{
				{
					Name:       "id",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: true,
					AutoID:     true,
				},
				{
					Name:     "result_id",
					DataType: entity.FieldTypeInt64,
				},
				{
					Name:     "user_id",
					DataType: entity.FieldTypeInt64,
				},
				{
					Name:       "level",
					DataType:   entity.FieldTypeVarChar,
					TypeParams: map[string]string{"max_length": "20"},
				},
				{
					Name:       "feature_vector",
					DataType:   entity.FieldTypeFloatVector,
					TypeParams: map[string]string{"dim": strconv.Itoa(milvusDimension)},
				},
			},
		}

		// 创建集合
		err = mc.client.CreateCollection(ctx, schema, 2) // 2 shards
		if err != nil {
			return fmt.Errorf("创建集合失败: %w", err)
		}

		// 创建IVF索引
		idx, err := entity.NewIndexIvfFlat(entity.IP, 128)
		if err != nil {
			return fmt.Errorf("创建索引失败: %w", err)
		}
		err = mc.client.CreateIndex(ctx, milvusCollectionName, "feature_vector", idx, false)
		if err != nil {
			return fmt.Errorf("创建索引失败: %w", err)
		}

		// 加载集合到内存
		err = mc.client.LoadCollection(ctx, milvusCollectionName, false)
		if err != nil {
			return fmt.Errorf("加载集合失败: %w", err)
		}

		logx.Info("Milvus集合创建成功")
	} else {
		// 确保集合已加载
		err = mc.client.LoadCollection(ctx, milvusCollectionName, false)
		if err != nil {
			logx.Infof("集合已加载或加载失败(可忽略): %v", err)
		}
	}

	return nil
}

// FeatureRecord 特征记录
type FeatureRecord struct {
	ResultID  int64
	UserID    int64
	Level     string
	Vector    []float32
}

// InsertFeature 插入特征向量
func (mc *MilvusClient) InsertFeature(ctx context.Context, record FeatureRecord) (int64, error) {
	// 准备列数据（id 为 AutoID，无需插入）
	resultIDColumn := entity.NewColumnInt64("result_id", []int64{record.ResultID})
	userIDColumn := entity.NewColumnInt64("user_id", []int64{record.UserID})
	levelColumn := entity.NewColumnVarChar("level", []string{record.Level})
	vectorColumn := entity.NewColumnFloatVector("feature_vector", milvusDimension, [][]float32{record.Vector})

	_, err := mc.client.Insert(ctx, milvusCollectionName, "",
		resultIDColumn, userIDColumn, levelColumn, vectorColumn)
	if err != nil {
		return 0, fmt.Errorf("插入特征向量失败: %w", err)
	}

	// Flush确保数据持久化
	err = mc.client.Flush(ctx, milvusCollectionName, false)
	if err != nil {
		logx.Errorf("Flush失败: %v", err)
	}

	logx.Infof("特征向量已插入: resultID=%d, level=%s", record.ResultID, record.Level)
	return 0, nil
}

// SearchResult 搜索结果
type SearchResult struct {
	ResultID  int64
	UserID    int64
	Level     string
	Score     float32
}

// SearchSimilar 搜索相似向量（RAG检索）
func (mc *MilvusClient) SearchSimilar(ctx context.Context, queryVector []float32, topK int) ([]SearchResult, error) {
	sp, err := entity.NewIndexIvfFlatSearchParam(128)
	if err != nil {
		return nil, fmt.Errorf("创建搜索参数失败: %w", err)
	}

	searchResult, err := mc.client.Search(ctx, milvusCollectionName, nil,
		"", // expr (过滤条件)
		[]string{"result_id", "user_id", "level"}, // output fields
		[]entity.Vector{entity.FloatVector(queryVector)},
		"feature_vector", // vector field name
		entity.IP,        // metric type (Inner Product)
		topK,             // topK
		sp,
	)
	if err != nil {
		return nil, fmt.Errorf("向量搜索失败: %w", err)
	}

	var results []SearchResult
	for _, sr := range searchResult {
		for i := 0; i < sr.ResultCount; i++ {
			var result SearchResult
			result.Score = sr.Scores[i]

			// 获取字段值
			resultID, err := sr.IDs.GetAsInt64(i)
			if err == nil {
				result.ResultID = resultID
			}

			// 获取output fields
			fields := sr.Fields
			for _, field := range fields {
				switch field.Name() {
				case "result_id":
					val, err := field.GetAsInt64(i)
					if err == nil {
						result.ResultID = val
					}
				case "user_id":
					val, err := field.GetAsInt64(i)
					if err == nil {
						result.UserID = val
					}
				case "level":
					val, err := field.GetAsString(i)
					if err == nil {
						result.Level = val
					}
				}
			}

			results = append(results, result)
		}
	}

	logx.Infof("RAG检索完成: topK=%d, 结果数=%d", topK, len(results))
	return results, nil
}

// FormatReferenceCases 格式化参考病例为prompt文本
func FormatReferenceCases(results []SearchResult) string {
	if len(results) == 0 {
		return "暂无相似历史病例（系统初期运行，知识库尚未积累足够病例）"
	}

	text := ""
	for i, r := range results {
		similarity := (1 - r.Score) * 100 // 转换为相似度百分比
		if similarity < 0 {
			similarity = 0
		}
		text += fmt.Sprintf("病例%d（相似度 %.1f%%）：诊断为%s腺样体面容\n", i+1, similarity, r.Level)
		text += fmt.Sprintf("  病例ID: %d\n\n", r.ResultID)
	}
	return text
}

// Close 关闭连接
func (mc *MilvusClient) Close() error {
	if mc.client != nil {
		return mc.client.Close()
	}
	return nil
}

// MockFeatureVector 生成模拟特征向量（实际使用时替换为EfficientNet提取）
func MockFeatureVector() []float32 {
	vector := make([]float32, milvusDimension)
	// 使用伪随机但确定性的方式生成向量
	for i := 0; i < milvusDimension; i++ {
		vector[i] = float32(i%100) / 100.0
	}
	return vector
}

// MockFeatureDescription 生成模拟特征描述（实际使用时替换为EfficientNet提取）
func MockFeatureDescription() string {
	desc := map[string]float64{
		"上唇上翘程度":   0.87,
		"鼻唇沟变浅程度":  0.63,
		"下颌后缩程度":   0.72,
		"牙列不齐程度":   0.45,
		"张口呼吸特征":   0.91,
		"面部比例失调程度": 0.68,
	}
	data, _ := json.Marshal(desc)
	return string(data)
}
