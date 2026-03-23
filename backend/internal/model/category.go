package model

type Category struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(32);unique;not null"`
}
