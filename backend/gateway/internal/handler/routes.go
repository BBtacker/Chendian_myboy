package handler

import (
	"net/http"

	"faceTest/backend/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterRoutes 注册所有HTTP路由
func RegisterRoutes(server *rest.Server, ctx *svc.ServiceContext) {
	// ==================== 不需要认证的路由 ====================
	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: LoginHandler(ctx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: RegisterHandler(ctx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/user/register",
			Handler: RegisterHandler(ctx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/logout",
			Handler: LogoutHandler(ctx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: HealthHandler(ctx),
		},
	})

	// ==================== 需要认证的路由 ====================
	// go-zero v1.7.6: WithMiddleware(m, rs...) 返回包装后的 []Route，作为 AddRoutes 第一个参数
	server.AddRoutes(
		rest.WithMiddleware(ctx.AuthMiddleware.Handle, []rest.Route{
			// 用户管理
			{
				Method:  http.MethodGet,
				Path:    "/user",
				Handler: GetUserHandler(ctx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/user",
				Handler: UpdateUserHandler(ctx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/avatar",
				Handler: UploadAvatarHandler(ctx),
			},
			// 面容分析
			{
				Method:  http.MethodPost,
				Path:    "/doubao/analyzeFace",
				Handler: AnalyzeFaceHandler(ctx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/doubao/analyzeMulti",
				Handler: AnalyzeFaceMultiHandler(ctx),
			},
			// 检测记录
			{
				Method:  http.MethodGet,
				Path:    "/testResult/result",
				Handler: ListResultsHandler(ctx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/testResult/{id}",
				Handler: DeleteResultHandler(ctx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/testResult/batch",
				Handler: BatchDeleteResultsHandler(ctx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/testResult/update",
				Handler: UpdateResultHandler(ctx),
			},
			// 导出
			{
				Method:  http.MethodGet,
				Path:    "/testResult/download",
				Handler: ExportExcelHandler(ctx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/testResult/download/pdf",
				Handler: ExportPDFHandler(ctx),
			},
			// 统计
			{
				Method:  http.MethodGet,
				Path:    "/statistics/overview",
				Handler: StatisticsOverviewHandler(ctx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/statistics/detail",
				Handler: StatisticsDetailHandler(ctx),
			},
			// 对话
			{
				Method:  http.MethodPost,
				Path:    "/conversation/create",
				Handler: CreateConversationHandler(ctx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/conversation/send",
				Handler: SendConversationHandler(ctx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/conversation/list",
				Handler: ListConversationsHandler(ctx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/conversation/messages",
				Handler: GetMessagesHandler(ctx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/conversation/{id}",
				Handler: DeleteConversationHandler(ctx),
			},
		}...),
	)

	// 静态文件服务（上传的文件）
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/uploads/{file}",
		Handler: StaticFileHandler(ctx),
	})
}
