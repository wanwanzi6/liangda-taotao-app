package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement"`
	UserID      uint64          `gorm:"index:idx_user_id;not null"`
	CategoryID  uint64          `gorm:"not null"`
	Title       string          `gorm:"type:varchar(128);not null"`
	ImageURL    string          `gorm:"type:varchar(255)"`
	Description string          `gorm:"type:text"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Type        int             `gorm:"default:1;comment:1:实物, 2:租赁"`
	Status      int             `gorm:"default:1;comment:1:待售, 2:交易中, 3:已售"`
	UpdatedAt   time.Time       `gorm:"precision:3;index:idx_updated_at;autoUpdateTime:milli"`
	DeletedAt   gorm.DeletedAt  `gorm:"index"`

	// 关联关系
	User     User     `gorm:"foreignKey:UserID;preload:false"`
	Category Category `gorm:"foreignKey:CategoryID;preload:false"`
}
