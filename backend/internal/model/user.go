package model

import "time"

type User struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time `gorm:"precision:3"`
	UpdatedAt   time.Time `gorm:"precision:3"`
	OpenID      string    `gorm:"type:varchar(128);unique;not null;comment:微信OpenID"`
	SessionKey  string    `gorm:"type:varchar(256);comment:微信SessionKey(用于解密用户数据)"`
	Nickname    string    `gorm:"type:varchar(64)"`
	AvatarURL   string    `gorm:"type:varchar(255)"`
	IsVerify    bool      `gorm:"type:tinyint(1);default:0;comment:校内认证状态"`
	LastLoginAt time.Time `gorm:"precision:3;comment:最后登录时间"`
}
