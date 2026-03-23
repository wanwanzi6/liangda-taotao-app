package repository

import (
	"liangda-taotao/internal/model"

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
	// Save 会更新对象的所有字段
	return r.db.Save(product).Error
}

// 查询商品
func (r *ProductRepository) GetByID(id uint64) (*model.Product, error) {
	var product model.Product
	// First 会根据主键查询，如果找不到会返回 gorm.ErrRecordNotFound
	err := r.db.First(&product, id).Error
	return &product, err
}

// 分页查询商品
func (r *ProductRepository) GetList(page, pageSize int) ([]model.Product, error) {
	var products []model.Product
	offset := (page - 1) * pageSize

	// 按 updated_at 倒序排序，实现“最新擦亮/发布”排在最前面
	err := r.db.Order("updated_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	return products, err
}

// 按分类分页查询
func (r *ProductRepository) GetListByCategory(categoryID uint64, page, pageSize int) ([]model.Product, error) {
	var products []model.Product
	offset := (page - 1) * pageSize
	err := r.db.Where("category_id = ?", categoryID).
		Order("updated_at desc").
		Offset(offset).Limit(pageSize).
		Find(&products).Error
	return products, err
}

// 一键擦亮 - 重点业务逻辑
func (r *ProductRepository) Refresh(id uint64) error {
	// 只需要更新 updated_at，GORM 的 autoUpdatedTime 会自动处理
	return r.db.Model(&model.Product{}).
		Where("id = ?", id).
		Update("updated_at", gorm.
			Expr("updated_at")).Error
}
