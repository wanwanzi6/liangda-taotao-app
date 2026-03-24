package handler

import (
	"fmt"
	"liangda-taotao/internal/model"
	"liangda-taotao/internal/repository"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
