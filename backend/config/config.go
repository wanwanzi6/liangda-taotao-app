package config

import "fmt"

// DBConfig 存储数据库连接信息
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

// GetDSN 生成 GORM 需要的连接字符串 (Data Source Name)
func GetDSN(cfg DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

// 先手动填入配置，后续升级为读取 .yaml 文件
var ApppDBconfig = DBConfig{
	User:     "root",
	Password: "lkx411",
	Host:     "127.0.0.1",
	Port:     3306,
	DBName:   "liang_da_tao_tao",
}

// WeChatConfig 微信小程序配置
type WeChatConfig struct {
	AppID     string
	AppSecret string
}

var WeChat = WeChatConfig{
	AppID:     "wx6015d4f584306826",
	AppSecret: "YOUR_APP_SECRET", // TODO: 替换为实际的 AppSecret
}
