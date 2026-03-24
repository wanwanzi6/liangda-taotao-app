package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// 先加载 .env 文件
func init() {
	// 优先尝试当前工作目录
	envPath := ".env"
	if err := godotenv.Load(envPath); err != nil {
		// 如果找不到，尝试 backend/.env
		envPath = "backend/.env"
		if err := godotenv.Load(envPath); err != nil {
			log.Println("⚠️ 未找到 .env 文件，将使用系统环境变量")
		}
	}
	log.Printf("📂 已加载环境变量: %s", envPath)
}

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

// AppDBConfig 从环境变量读取数据库配置
var AppDBConfig = DBConfig{
	User:     getEnv("DB_USER", "root"),
	Password: getEnv("DB_PASSWORD", "lkx411"),
	Host:     getEnv("DB_HOST", "127.0.0.1"),
	Port:     getEnvInt("DB_PORT", 3306),
	DBName:   getEnv("DB_NAME", "liang_da_tao_tao"),
}

// WeChatConfig 微信小程序配置
type WeChatConfig struct {
	AppID     string
	AppSecret string
}

// GetWeChat 获取微信配置
func GetWeChat() WeChatConfig {
	return WeChatConfig{
		AppID:     getEnv("WECHAT_APPID", ""),
		AppSecret: getEnv("WECHAT_SECRET", ""),
	}
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string
}

var JWT = JWTConfig{
	Secret: getEnv("JWT_SECRET", "default-secret-key"),
}

// getEnv 获取环境变量，如果不存在返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取环境变量 int 值
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intVal int
		fmt.Sscanf(value, "%d", &intVal)
		return intVal
	}
	return defaultValue
}
