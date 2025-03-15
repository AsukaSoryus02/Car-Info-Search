package config

import (
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jasonzheng/carrag/utils"
)

// AppConfig 应用配置
type AppConfig struct {
	// 服务器配置
	ServerPort int
	GinMode    string

	// 路径配置
	LogDir  string
	DataDir string

	// Redis配置
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	RedisPrefix   string

	// 日志级别
	LogLevel utils.LogLevel
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *AppConfig {
	return &AppConfig{
		ServerPort:    8080,
		GinMode:       gin.ReleaseMode,
		LogDir:        filepath.Join("logs"),
		DataDir:       filepath.Join("data"),
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		RedisPrefix:   "carrag",
		LogLevel:      utils.INFO,
	}
}

// GetCorsConfig 获取CORS配置
func GetCorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}
