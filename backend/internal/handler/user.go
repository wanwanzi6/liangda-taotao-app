package handler

import (
	"encoding/json"
	"liangda-taotao/config"
	"liangda-taotao/internal/middleware"
	"liangda-taotao/internal/model"
	"liangda-taotao/internal/repository"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// WeChatLogin 微信登录
// 前端调用 wx.login() 获取 code，发送到后端
// 后端用 code 换 openid，返回 token
func (h *UserHandler) WeChatLogin(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "缺少 code 参数")
		return
	}

	// 1. 用 code 换取 openid 和 session_key
	openID, sessionKey, err := h.getWeChatOpenID(req.Code)
	if err != nil {
		log.Printf("[WeChatLogin] 获取 openid 失败: %v", err)
		ServerError(c, "微信登录失败: "+err.Error())
		return
	}

	// 2. 查询用户是否存在，不存在则创建
	user, err := h.repo.GetByOpenID(openID)
	isNewUser := false
	if err == gorm.ErrRecordNotFound {
		// 新用户注册
		isNewUser = true
		user = &model.User{
			OpenID:      openID,
			SessionKey:  sessionKey,
			LastLoginAt: time.Now(),
		}
		if err := h.repo.Create(user); err != nil {
			ServerError(c, "创建用户失败")
			return
		}
		log.Printf("[WeChatLogin] 新用户注册: user_id=%d, openid=%s", user.ID, openID)
	} else if err != nil {
		ServerError(c, "查询用户失败")
		return
	} else {
		// 老用户更新 session_key 和最后登录时间
		user.SessionKey = sessionKey
		user.LastLoginAt = time.Now()
		if err := h.repo.Update(user); err != nil {
			log.Printf("[WeChatLogin] 更新用户信息失败: %v", err)
		}
	}

	// 3. 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, user.OpenID, user.Nickname)
	if err != nil {
		ServerError(c, "生成 token 失败")
		return
	}

	log.Printf("[WeChatLogin] 用户登录成功: user_id=%d, is_new=%v", user.ID, isNewUser)

	// 4. 返回 token 和用户信息
	Success(c, gin.H{
		"token":      token,
		"is_new":     isNewUser,
		"expires_in": 7 * 24 * 3600, // token 有效期（秒）
		"user": gin.H{
			"id":          user.ID,
			"nickname":    user.Nickname,
			"avatar_url":  user.AvatarURL,
			"is_verify":   user.IsVerify,
			"last_login":  user.LastLoginAt,
		},
	})
}

// getWeChatOpenID 调用微信 API 获取 openid 和 session_key
func (h *UserHandler) getWeChatOpenID(code string) (openID, sessionKey string, err error) {
	wechatConfig := config.GetWeChat()
	appID := wechatConfig.AppID
	appSecret := wechatConfig.AppSecret

	if appID == "" || appSecret == "" {
		return "", "", &WeChatError{Code: -1, Msg: "微信配置未正确设置"}
	}

	// 调用微信登录凭证校验接口
	apiURL := "https://api.weixin.qq.com/sns/jscode2session?" +
		"appid=" + appID +
		"&secret=" + appSecret +
		"&js_code=" + code +
		"&grant_type=authorization_code"

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", "", &WeChatError{Code: -2, Msg: "请求微信接口失败: " + err.Error()}
	}
	defer resp.Body.Close()

	var result struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", &WeChatError{Code: -3, Msg: "解析微信响应失败"}
	}

	if result.ErrCode != 0 {
		return "", "", &WeChatError{Code: result.ErrCode, Msg: getWeChatErrorMsg(result.ErrCode, result.ErrMsg)}
	}

	if result.OpenID == "" {
		return "", "", &WeChatError{Code: -1, Msg: "未获取到 openid"}
	}

	return result.OpenID, result.SessionKey, nil
}

// getWeChatErrorMsg 微信错误码转中文提示
func getWeChatErrorMsg(code int, msg string) string {
	errorMsgs := map[int]string{
		-1:   "系统繁忙，请稍后再试",
		0:    "请求成功",
		40029: "code 无效",
		45011: "API 调用频率限制",
		40226: "高风险用户，需慎重",
	}
	if m, ok := errorMsgs[code]; ok {
		return m
	}
	return msg
}

// WeChatError 微信 API 错误
type WeChatError struct {
	Code int
	Msg  string
}

func (e *WeChatError) Error() string {
	return e.Msg
}

// UpdateUserInfo 更新用户信息
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	// 从上下文获取登录用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	var req struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误")
		return
	}

	user, err := h.repo.GetByID(userID.(uint64))
	if err != nil {
		NotFound(c, "用户不存在")
		return
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := h.repo.Update(user); err != nil {
		ServerError(c, "更新失败")
		return
	}

	SuccessMsg(c, "更新成功", user)
}

// GetUserInfo 获取当前用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	user, err := h.repo.GetByID(userID.(uint64))
	if err != nil {
		NotFound(c, "用户不存在")
		return
	}

	Success(c, gin.H{
		"id":         user.ID,
		"open_id":    user.OpenID,
		"nickname":   user.Nickname,
		"avatar_url": user.AvatarURL,
		"is_verify":  user.IsVerify,
	})
}

// URLEncode 简单 URL 编码
func URLEncode(s string) string {
	return url.QueryEscape(s)
}

// RefreshToken 刷新 Token
// 前端在 token 过期前调用此接口刷新，不需要重新走微信登录
func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	user, err := h.repo.GetByID(userID.(uint64))
	if err != nil {
		NotFound(c, "用户不存在")
		return
	}

	// 生成新的 token
	token, err := middleware.GenerateToken(user.ID, user.OpenID, user.Nickname)
	if err != nil {
		ServerError(c, "刷新 token 失败")
		return
	}

	log.Printf("[RefreshToken] 用户刷新 token: user_id=%d", user.ID)

	Success(c, gin.H{
		"token":      token,
		"expires_in": 7 * 24 * 3600,
	})
}

// Logout 退出登录（可选：后端可以记录日志，前端删除 localStorage）
func (h *UserHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if exists {
		log.Printf("[Logout] 用户退出登录: user_id=%d", userID)
	}
	SuccessMsg(c, "已退出登录", nil)
}
