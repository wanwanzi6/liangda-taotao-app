package handler

import (
	"fmt"
	"liangda-taotao/internal/model"
	"liangda-taotao/internal/repository"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) GetList(c *gin.Context) {
	// 接收分页参数，默认第 1 页，每页 10 条
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "10"))
	// 接收分类 ID 参数
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)

	var products []model.Product
	var err error

	if categoryID > 0 {
		products, err = h.repo.GetListByCategory(categoryID, page, pageSize)
	} else {
		products, err = h.repo.GetList(page, pageSize)
	}

	if err != nil {
		ServerError(c, "获取商品列表失败")
		return
	}

	Success(c, products)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var p model.Product
	// 将前端传来的 JSON 绑定到结构体
	if err := c.ShouldBindJSON(&p); err != nil {
		BadRequest(c, "参数格式错误")
		return
	}

	// 基础校验：价格不能为负
	if p.Price.IsNegative() {
		BadRequest(c, "价格不能为负数")
		return
	}

	// 商品名称不能为空
	if p.Title == "" {
		BadRequest(c, "商品名称不能为空")
		return
	}

	// 从 JWT token 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}
	p.UserID = userID.(uint64)
	p.Status = 1 // 初始状态为待售

	if err := h.repo.Create(&p); err != nil {
		ServerError(c, "发布失败")
		return
	}

	SuccessMsg(c, "发布成功", p)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	// 获取 URL 中的 id 参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的商品ID")
		return
	}

	if err := h.repo.Delete(uint64(id)); err != nil {
		ServerError(c, "删除失败")
		return
	}

	SuccessMsg(c, "商品已成功下架", nil)
}

func (h *ProductHandler) GetDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的商品ID格式")
		return
	}

	product, err := h.repo.GetByID(id)
	if err != nil {
		NotFound(c, "该商品不存在或已下架")
		return
	}

	Success(c, product)
}

func (h *ProductHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		BadRequest(c, "未获取到文件")
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	savePath := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		ServerError(c, "保存失败")
		return
	}

	Success(c, gin.H{"url": "/uploads/" + filename})
}

// GetUserProducts 获取当前用户的商品列表（我的发布）
func (h *ProductHandler) GetUserProducts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	products, err := h.repo.GetByUserID(userID.(uint64))
	if err != nil {
		ServerError(c, "获取我的发布失败")
		return
	}

	Success(c, products)
}

// Refresh 一键擦亮
func (h *ProductHandler) Refresh(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的商品ID")
		return
	}

	// 验证商品所有权
	product, err := h.repo.GetByID(id)
	if err != nil {
		NotFound(c, "商品不存在")
		return
	}
	if product.UserID != userID.(uint64) {
		BadRequest(c, "无权操作此商品")
		return
	}

	if err := h.repo.Refresh(id); err != nil {
		ServerError(c, "擦亮失败")
		return
	}

	SuccessMsg(c, "擦亮成功，商品已置顶", nil)
}

// Update 更新商品信息
func (h *ProductHandler) Update(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, "无效的商品ID")
		return
	}

	// 验证商品所有权
	product, err := h.repo.GetByID(id)
	if err != nil {
		NotFound(c, "商品不存在")
		return
	}
	if product.UserID != userID.(uint64) {
		BadRequest(c, "无权操作此商品")
		return
	}

	var req struct {
		Title       string          `json:"title"`
		Description string          `json:"description"`
		Price       string          `json:"price"`
		ImageURL    string          `json:"image_url"`
		CategoryID  uint64          `json:"category_id"`
		Type        int             `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数格式错误")
		return
	}

	// 更新字段
	if req.Title != "" {
		product.Title = req.Title
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}
	if req.Type > 0 {
		product.Type = req.Type
	}
	if req.Price != "" {
		// 使用 shopspring/decimal 解析价格
		price, err := decimal.NewFromString(req.Price)
		if err == nil && !price.IsNegative() {
			product.Price = price
		}
	}

	if err := h.repo.Update(product); err != nil {
		ServerError(c, "更新失败")
		return
	}

	SuccessMsg(c, "更新成功", product)
}
