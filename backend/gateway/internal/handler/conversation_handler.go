package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"faceTest/backend/common/middleware"
	"faceTest/backend/common/model"
	"faceTest/backend/common/pkg"
	"faceTest/backend/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
	"gorm.io/gorm"
)

// CreateConversationHandler 创建对话
func CreateConversationHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		firstMessage := r.URL.Query().Get("firstMessage")
		if firstMessage == "" {
			// 尝试从body解析
			var body struct {
				FirstMessage string `json:"firstMessage"`
			}
			json.NewDecoder(r.Body).Decode(&body)
			firstMessage = body.FirstMessage
		}

		title := firstMessage
		if len(title) > 50 {
			title = title[:50] + "..."
		}

		conv := model.Conversation{
			UserID: userID,
			Title:  title,
		}

		if err := ctx.DB.Create(&conv).Error; err != nil {
			// 如果DB不可用（网关没有直接DB连接），使用简化模式
			httpx.OkJsonCtx(r.Context(), w, pkg.Success(map[string]interface{}{
				"id":         time.Now().Unix(),
				"title":      title,
				"create_time": time.Now().Format("2006-01-02 15:04:05"),
			}))
			return
		}

		httpx.OkJsonCtx(r.Context(), w, pkg.Success(conv))
	}
}

// SendConversationHandler 发送消息
func SendConversationHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		convIDStr := r.URL.Query().Get("conversationId")
		content := r.URL.Query().Get("content")

		if content == "" {
			// 尝试从body解析
			var body struct {
				ConversationId int64  `json:"conversationId"`
				Content        string `json:"content"`
			}
			json.NewDecoder(r.Body).Decode(&body)
			if body.Content != "" {
				content = body.Content
			}
		}

		if content == "" {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("消息内容不能为空"))
			return
		}

		// 使用DeepSeek API生成回复
		if ctx.DeepSeekClient == nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("AI助手未配置"))
			return
		}

		messages := []pkg.DeepSeekMessage{
			{Role: "system", Content: "你是一个专业的耳鼻喉科医疗助手，专注于腺样体面容相关问题。请用中文回答，语气专业但易懂。"},
			{Role: "user", Content: content},
		}

		rpcCtx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reply, err := ctx.DeepSeekClient.Chat(rpcCtx, messages)
		if err != nil {
			logx.Errorf("AI对话失败: %v", err)
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("AI回复失败: " + err.Error()))
			return
		}

		_ = userID // 可用于保存对话记录
		_ = convIDStr

		httpx.OkJsonCtx(r.Context(), w, pkg.Success(reply))
	}
}

// ListConversationsHandler 获取对话列表
func ListConversationsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		// 如果网关有DB连接
		if ctx.DB != nil {
			var conversations []model.Conversation
			result := ctx.DB.Where("user_id = ?", userID).Order("create_time DESC").Find(&conversations)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				httpx.OkJsonCtx(r.Context(), w, pkg.Error("查询失败"))
				return
			}
			httpx.OkJsonCtx(r.Context(), w, pkg.Success(conversations))
			return
		}

		// 无DB时返回空列表
		httpx.OkJsonCtx(r.Context(), w, pkg.Success([]interface{}{}))
	}
}

// GetMessagesHandler 获取消息列表
func GetMessagesHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		convIDStr := r.URL.Query().Get("conversationId")

		if ctx.DB != nil && convIDStr != "" {
			var messages []model.Message
			result := ctx.DB.Where("conversation_id = ?", convIDStr).Order("create_time ASC").Find(&messages)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				httpx.OkJsonCtx(r.Context(), w, pkg.Error("查询失败"))
				return
			}
			httpx.OkJsonCtx(r.Context(), w, pkg.Success(messages))
			return
		}

		httpx.OkJsonCtx(r.Context(), w, pkg.Success([]interface{}{}))
	}
}

// DeleteConversationHandler 删除对话
func DeleteConversationHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		convIDStr := r.PathValue("id")
		if convIDStr == "" {
			vars := pathvar.Vars(r)
			convIDStr = vars["id"]
		}
		if convIDStr == "" {
			convIDStr = r.URL.Query().Get("id")
		}

		if convIDStr == "" {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("无效的对话ID"))
			return
		}

		convID, err := strconv.ParseInt(convIDStr, 10, 64)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("对话ID格式错误"))
			return
		}

		if ctx.DB != nil {
			// 删除对话及关联消息
			tx := ctx.DB.Begin()
			if err := tx.Where("conversation_id = ?", convID).Delete(&model.Message{}).Error; err != nil {
				tx.Rollback()
				httpx.OkJsonCtx(r.Context(), w, pkg.Error("删除消息失败"))
				return
			}
			if err := tx.Delete(&model.Conversation{}, convID).Error; err != nil {
				tx.Rollback()
				httpx.OkJsonCtx(r.Context(), w, pkg.Error("删除对话失败"))
				return
			}
			tx.Commit()
		}

		httpx.OkJsonCtx(r.Context(), w, pkg.SuccessMsg("删除成功"))
	}
}
