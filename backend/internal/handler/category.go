package handler

import (
	"liangda-taotao/internal/repository"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo *repository.CategoryRepository
}

func NewCategoryHandler(repo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.repo.GetAll()
	if err != nil {
		ServerError(c, "获取分类失败")
		return
	}
	Success(c, categories)
}
