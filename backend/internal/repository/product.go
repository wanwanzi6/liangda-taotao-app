package repository

import (
	"liangda-taotao/internal/model"
	"time"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// 发布商品
func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

// 软删除商品
func (r *ProductRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Product{}, id).Error
}

// 更新商品信息
func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

// 查询商品（带预加载）
func (r *ProductRepository) GetByID(id uint64) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").First(&product, id).Error
	return &product, err
}

// 分页查询商品（带预加载分类和用户信息，避免 N+1）
func (r *ProductRepository) GetList(page, pageSize int) ([]model.Product, error) {
	var products []model.Product
	offset := (page - 1) * pageSize

	err := r.db.Preload("Category").Preload("User").
		Order("updated_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	return products, err
}

// 按分类分页查询（带预加载）
func (r *ProductRepository) GetListByCategory(categoryID uint64, page, pageSize int) ([]model.Product, error) {
	var products []model.Product
	offset := (page - 1) * pageSize
	err := r.db.Preload("Category").Preload("User").
		Where("category_id = ?", categoryID).
		Order("updated_at desc").
		Offset(offset).Limit(pageSize).
		Find(&products).Error
	return products, err
}

// 一键擦亮 - 重点业务逻辑
func (r *ProductRepository) Refresh(id uint64) error {
	return r.db.Model(&model.Product{}).
		Where("id = ?", id).
		Update("updated_at", time.Now()).Error
}

// 按用户ID查询商品（我的发布，带预加载）
func (r *ProductRepository) GetByUserID(userID uint64) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Preload("Category").
		Where("user_id = ?", userID).
		Order("updated_at desc").
		Find(&products).Error
	return products, err
}
