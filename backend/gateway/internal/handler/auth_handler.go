package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"faceTest/backend/common/middleware"
	"faceTest/backend/common/pkg"
	"faceTest/backend/gateway/internal/svc"
	authpb "faceTest/backend/proto/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// LoginHandler 登录
func LoginHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("参数解析失败"))
			return
		}

		// 调用认证服务
		rpcCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		resp, err := ctx.AuthClient.Login(rpcCtx, &authpb.LoginReq{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			logx.Errorf("登录RPC调用失败: %v", err)
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("服务调用失败"))
			return
		}

		if resp.Code == 1 {
			// 登录成功，返回token和数据
			httpx.OkJsonCtx(r.Context(), w, pkg.Success(resp.Token))
		} else {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error(resp.Msg))
		}
	}
}

// RegisterHandler 注册
func RegisterHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Name     string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("参数解析失败"))
			return
		}

		rpcCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		resp, err := ctx.AuthClient.Register(rpcCtx, &authpb.RegisterReq{
			Username: req.Username,
			Password: req.Password,
			Name:     req.Name,
		})
		if err != nil {
			logx.Errorf("注册RPC调用失败: %v", err)
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("服务调用失败"))
			return
		}

		if resp.Code == 1 {
			httpx.OkJsonCtx(r.Context(), w, pkg.SuccessMsg(resp.Msg))
		} else {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error(resp.Msg))
		}
	}
}

// LogoutHandler 登出
func LogoutHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.OkJsonCtx(r.Context(), w, pkg.SuccessMsg("登出成功"))
	}
}

// GetUserHandler 获取当前用户信息
func GetUserHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		rpcCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		resp, err := ctx.AuthClient.GetUser(rpcCtx, &authpb.GetUserReq{UserId: int64(userID)})
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("获取用户信息失败"))
			return
		}

		httpx.OkJsonCtx(r.Context(), w, pkg.Success(resp))
	}
}

// UpdateUserHandler 更新用户信息
func UpdateUserHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		var req struct {
			Name     string `json:"name,omitempty"`
			Password string `json:"password,omitempty"`
			Email    string `json:"email,omitempty"`
			Phone    string `json:"phone,omitempty"`
			Gender   int32  `json:"gender,omitempty"`
			Birthday string `json:"birthday,omitempty"`
			Address  string `json:"address,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("参数解析失败"))
			return
		}

		rpcCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		resp, err := ctx.AuthClient.UpdateUser(rpcCtx, &authpb.UpdateUserReq{
			UserId: int64(userID),
			Name:     req.Name,
			Password: req.Password,
			Email:    req.Email,
			Phone:    req.Phone,
			Gender:   req.Gender,
			Birthday: req.Birthday,
			Address:  req.Address,
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

// UploadAvatarHandler 上传头像
func UploadAvatarHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		// 限制文件大小
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB

		file, header, err := r.FormFile("avatar")
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("文件读取失败"))
			return
		}
		defer file.Close()

		// 读取文件内容
		fileData, err := io.ReadAll(file)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("文件读取失败"))
			return
		}

		rpcCtx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		resp, err := ctx.AuthClient.UploadAvatar(rpcCtx, &authpb.UploadAvatarReq{
			UserId: int64(userID),
			FileData: fileData,
			Filename: header.Filename,
		})
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, pkg.Error("头像上传失败"))
			return
		}

		httpx.OkJsonCtx(r.Context(), w, pkg.Success(resp.AvatarUrl))
	}
}

// HealthHandler 健康检查
func HealthHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.OkJsonCtx(r.Context(), w, pkg.Success(map[string]string{
			"status":  "ok",
			"service": "gateway",
			"time":    time.Now().Format("2006-01-02 15:04:05"),
		}))
	}
}

// StaticFileHandler 静态文件服务
func StaticFileHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 简单的静态文件服务
		filePath := r.URL.Path
		if len(filePath) > 9 { // strip "/uploads/"
			filePath = filePath[9:]
			fullPath := ctx.FileStorage.GetFullPath(filePath)
			http.ServeFile(w, r, fullPath)
			return
		}
		http.NotFound(w, r)
	}
}

// helper: parse int from query param
func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// helper: parse bool from query param
func parseBool(s string) bool {
	if s == "" {
		return false
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}

// avoid unused import
var _ = fmt.Sprintf
