package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"faceTest/backend/common/pkg"

	"github.com/zeromicro/go-zero/core/logx"
)

// AuthMiddleware JWT认证中间件
type AuthMiddleware struct {
	excludePaths map[string]bool
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware() *AuthMiddleware {
	excludePaths := map[string]bool{
		"/login":               true,
		"/user/register":       true,
		"/register":            true,
		"/health":              true,
		"/api/login":           true,
		"/api/user/register":   true,
		"/api/register":        true,
		"/api/health":          true,
	}

	return &AuthMiddleware{excludePaths: excludePaths}
}

// Handle 认证中间件处理
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// 检查是否为排除路径
		if m.excludePaths[path] {
			next(w, r)
			return
		}

		// OPTIONS请求直接放行
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		// 从Header获取Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.writeUnauthorized(w, "未提供认证Token")
			return
		}

		// 移除Bearer前缀（如果有）
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			m.writeUnauthorized(w, "Token为空")
			return
		}

		// 解析Token
		claims, err := pkg.ParseToken(token)
		if err != nil {
			logx.Infof("Token验证失败: %v, token=%s", err, token)
			m.writeUnauthorized(w, "Token无效或已过期")
			return
		}

		// 将用户信息存入Header，后续handler可获取
		r.Header.Set("X-User-Id", fmt.Sprintf("%d", claims.UserID))
		r.Header.Set("X-Username", claims.Username)
		r.Header.Set("X-User-Name", claims.Name)
		r.Header.Set("X-User-Role", claims.Role)

		next(w, r)
	}
}

// GetUserID 从请求中获取用户ID
func GetUserID(r *http.Request) uint64 {
	idStr := r.Header.Get("X-User-Id")
	if idStr == "" {
		return 0
	}
	var id uint64
	fmt.Sscanf(idStr, "%d", &id)
	return id
}

// GetUsername 从请求中获取用户名
func GetUsername(r *http.Request) string {
	return r.Header.Get("X-Username")
}

// GetUserRole 从请求中获取用户角色
func GetUserRole(r *http.Request) string {
	return r.Header.Get("X-User-Role")
}

func (m *AuthMiddleware) writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 返回200，但code=0，兼容前端处理
	w.Write([]byte(fmt.Sprintf(`{"code":0,"msg":"%s"}`, msg)))
}
