package handler

import (
	"encoding/json"
	"liangda-taotao/config"
	"liangda-taotao/internal/middleware"
	"liangda-taotao/internal/model"
	"liangda-taotao/internal/repository"
	"net/http"
	"net/url"

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

	// 1. 用 code 换取 openid
	openID, err := h.getWeChatOpenID(req.Code)
	if err != nil {
		ServerError(c, "微信登录失败: "+err.Error())
		return
	}

	// 2. 查询用户是否存在，不存在则创建
	user, err := h.repo.GetByOpenID(openID)
	if err == gorm.ErrRecordNotFound {
		// 新用户注册
		user = &model.User{OpenID: openID}
		if err := h.repo.Create(user); err != nil {
			ServerError(c, "创建用户失败")
			return
		}
	} else if err != nil {
		ServerError(c, "查询用户失败")
		return
	}

	// 3. 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, user.OpenID, user.Nickname)
	if err != nil {
		ServerError(c, "生成 token 失败")
		return
	}

	// 4. 返回 token 和用户信息
	Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"nickname":   user.Nickname,
			"avatar_url": user.AvatarURL,
			"is_verify":  user.IsVerify,
		},
	})
}

// getWeChatOpenID 调用微信 API 获取 openid
func (h *UserHandler) getWeChatOpenID(code string) (string, error) {
	appID := config.WeChat.AppID
	appSecret := config.WeChat.AppSecret

	// 调用微信登录凭证校验接口
	apiURL := "https://api.weixin.qq.com/sns/jscode2session?" +
		"appid=" + appID +
		"&secret=" + appSecret +
		"&js_code=" + code +
		"&grant_type=authorization_code"

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", &WeChatError{Code: result.ErrCode, Msg: result.ErrMsg}
	}

	if result.OpenID == "" {
		return "", &WeChatError{Code: -1, Msg: "未获取到 openid"}
	}

	return result.OpenID, nil
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
