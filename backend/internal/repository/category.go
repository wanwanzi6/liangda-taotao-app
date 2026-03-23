package repository

import (
	"liangda-taotao/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// 创建分类
func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// 获取所有分类
func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

// 根据 ID 查找分类
func (r *CategoryRepository) GetByID(id uint64) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	return &category, err
}
