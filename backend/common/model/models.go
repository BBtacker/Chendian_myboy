package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表模型
type User struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Username   string         `gorm:"size:50;uniqueIndex;not null;column:username" json:"username"`
	Password   string         `gorm:"size:255;not null;column:password" json:"-"`
	Name       string         `gorm:"size:50;not null;column:name" json:"name"`
	Avatar     string         `gorm:"size:500;column:avatar" json:"avatar"`
	Email      string         `gorm:"size:100;column:email" json:"email"`
	Phone      string         `gorm:"size:20;column:phone" json:"phone"`
	Gender     int8           `gorm:"column:gender;default:0" json:"gender"`
	Birthday   *time.Time     `gorm:"column:birthday" json:"birthday,omitempty"`
	Address    string         `gorm:"size:255;column:address" json:"address"`
	Role       string         `gorm:"size:20;not null;default:doctor;column:role" json:"role"`
	Version    int            `gorm:"not null;default:0;column:version" json:"-"`
	CreateTime time.Time      `gorm:"not null;autoCreateTime;column:create_time" json:"create_time"`
	UpdateTime time.Time      `gorm:"not null;autoUpdateTime;column:update_time" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "user"
}

// DiagnosisTask 诊断任务表
type DiagnosisTask struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TaskNo       string     `gorm:"size:64;uniqueIndex;not null;column:task_no" json:"task_no"`
	UserID       uint64     `gorm:"index;not null;column:user_id" json:"user_id"`
	ImagePath    string     `gorm:"size:500;not null;column:image_path" json:"image_path"`
	ImageURL     string     `gorm:"size:500;column:image_url" json:"image_url"`
	Status       int8       `gorm:"not null;default:0;index;column:status" json:"status"` // 0-待处理 1-处理中 2-已完成 3-失败
	Age          *int32     `gorm:"column:age" json:"age,omitempty"`
	Gender       *int8      `gorm:"column:gender" json:"gender,omitempty"`
	RetryCount   int        `gorm:"not null;default:0;column:retry_count" json:"retry_count"`
	MaxRetry     int        `gorm:"not null;default:5;column:max_retry" json:"max_retry"`
	ErrorMessage string     `gorm:"type:text;column:error_message" json:"error_message"`
	Version      int        `gorm:"not null;default:0;column:version" json:"-"`
	CreateTime   time.Time  `gorm:"not null;autoCreateTime;index;column:create_time" json:"create_time"`
	UpdateTime   time.Time  `gorm:"not null;autoUpdateTime;column:update_time" json:"update_time"`
	CompleteTime *time.Time `gorm:"column:complete_time" json:"complete_time,omitempty"`
}

func (DiagnosisTask) TableName() string {
	return "diagnosis_task"
}

// 诊断任务状态常量
const (
	TaskStatusPending   int8 = 0
	TaskStatusProcessing int8 = 1
	TaskStatusCompleted  int8 = 2
	TaskStatusFailed     int8 = 3
)

// DiagnosisResult 诊断结果表
type DiagnosisResult struct {
	ID                      uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TaskID                  uint64    `gorm:"index;not null;column:task_id" json:"task_id"`
	UserID                  uint64    `gorm:"index;not null;column:user_id" json:"user_id"`
	ImagePath               string    `gorm:"size:500;not null;column:image_path" json:"image_path"`
	ImageURL                string    `gorm:"size:500;column:image_url" json:"image_url"`
	IsGlandFace             bool      `gorm:"not null;default:false;index;column:is_gland_face" json:"is_gland_face"`
	Level                   string    `gorm:"size:20;index;column:level" json:"level"`
	VisualizationDescription string   `gorm:"type:text;column:visualization_description" json:"visualization_description"`
	FeatureVectorID         string    `gorm:"size:100;column:feature_vector_id" json:"feature_vector_id"`
	ReferenceCases          string    `gorm:"type:text;column:reference_cases" json:"reference_cases"`
	SkippedImages           string    `gorm:"type:text;column:skipped_images" json:"skipped_images"` // 多图诊断中被剔除（未检测到人脸）的图片路径 JSON 数组
	Version                 int       `gorm:"not null;default:0;column:version" json:"-"`
	TestTime                time.Time `gorm:"not null;index;column:test_time" json:"test_time"`
	CreateTime              time.Time `gorm:"not null;autoCreateTime;column:create_time" json:"create_time"`
	UpdateTime              time.Time `gorm:"not null;autoUpdateTime;column:update_time" json:"update_time"`
}

func (DiagnosisResult) TableName() string {
	return "diagnosis_result"
}

// OutboxMessage 本地消息表（Outbox Pattern）
type OutboxMessage struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	AggregateID   string     `gorm:"size:64;index;not null;column:aggregate_id" json:"aggregate_id"`
	AggregateType string     `gorm:"size:50;not null;column:aggregate_type" json:"aggregate_type"`
	MessageType   string     `gorm:"size:50;not null;column:message_type" json:"message_type"`
	Payload       string     `gorm:"type:text;not null;column:payload" json:"payload"`
	Status        int8       `gorm:"not null;default:0;index;column:status" json:"status"` // 0-待发送 1-已发送 2-已确认 3-失败
	RetryCount    int        `gorm:"not null;default:0;column:retry_count" json:"retry_count"`
	MaxRetry      int        `gorm:"not null;default:5;column:max_retry" json:"max_retry"`
	NextRetryTime *time.Time `gorm:"index;column:next_retry_time" json:"next_retry_time,omitempty"`
	CreateTime    time.Time  `gorm:"not null;autoCreateTime;column:create_time" json:"create_time"`
	UpdateTime    time.Time  `gorm:"not null;autoUpdateTime;column:update_time" json:"update_time"`
}

func (OutboxMessage) TableName() string {
	return "outbox_message"
}

// Outbox 消息状态常量
const (
	OutboxStatusPending   int8 = 0
	OutboxStatusSent      int8 = 1
	OutboxStatusConfirmed int8 = 2
	OutboxStatusFailed    int8 = 3
)

// Conversation 对话表
type Conversation struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID     uint64    `gorm:"index;not null;column:user_id" json:"user_id"`
	Title      string    `gorm:"size:100;column:title" json:"title"`
	Version    int       `gorm:"not null;default:0;column:version" json:"-"`
	CreateTime time.Time `gorm:"not null;autoCreateTime;column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;autoUpdateTime;column:update_time" json:"update_time"`
}

func (Conversation) TableName() string {
	return "conversation"
}

// Message 消息表
type Message struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ConversationID uint64    `gorm:"index;not null;column:conversation_id" json:"conversation_id"`
	Role           string    `gorm:"size:20;not null;column:role" json:"role"`
	Content        string    `gorm:"type:text;not null;column:content" json:"content"`
	CreateTime     time.Time `gorm:"not null;autoCreateTime;column:create_time" json:"create_time"`
}

func (Message) TableName() string {
	return "message"
}

// DiagnosisCache 诊断参数缓存表（三级缓存-MySQL层）
type DiagnosisCache struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CacheKey   string     `gorm:"size:200;uniqueIndex;not null;column:cache_key" json:"cache_key"`
	CacheValue string     `gorm:"type:text;not null;column:cache_value" json:"cache_value"`
	ExpireTime *time.Time `gorm:"column:expire_time" json:"expire_time,omitempty"`
	CreateTime time.Time  `gorm:"not null;autoCreateTime;column:create_time" json:"create_time"`
	UpdateTime time.Time  `gorm:"not null;autoUpdateTime;column:update_time" json:"update_time"`
}

func (DiagnosisCache) TableName() string {
	return "diagnosis_cache"
}
