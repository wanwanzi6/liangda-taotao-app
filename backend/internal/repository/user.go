package repository

import (
	"liangda-taotao/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// 注册用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// 通过 OpenID 查找用户 - 核心
func (r *UserRepository) GetByOpenID(openID string) (*model.User, error) {
	var user model.User
	// 使用 Where 条件查询
	err := r.db.Where("open_id = ?", openID).First(&user).Error
	return &user, err
}

// 通过数据库自增 ID 查找
func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// 修改个人资料 (昵称/头像/认证状态)
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
