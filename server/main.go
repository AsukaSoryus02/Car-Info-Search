package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jasonzheng/carrag/config"
	"github.com/jasonzheng/carrag/controllers"
	"github.com/jasonzheng/carrag/middleware"
	"github.com/jasonzheng/carrag/models"
	"github.com/jasonzheng/carrag/repositories"
	"github.com/jasonzheng/carrag/utils"
)

func main() {
	// 加载配置
	appConfig := config.NewDefaultConfig()

	// 初始化日志记录器
	logger, err := utils.NewLogger(appConfig.LogDir, appConfig.LogLevel)
	if err != nil {
		fmt.Printf("初始化日志记录器失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// 初始化存储管理器
	storage, err := utils.NewStorage(appConfig.DataDir, logger)
	if err != nil {
		logger.Fatal("初始化存储管理器失败: %v", err)
	}

	// 初始化文件仓库
	fileRepo := repositories.NewFileCarRepository(storage, logger, "cars.json")

	// 初始化Redis缓存
	var redisCache *utils.RedisCache
	redisCache, err = utils.NewRedisCache(
		appConfig.RedisAddr,
		appConfig.RedisPassword,
		appConfig.RedisDB,
		appConfig.RedisPrefix,
		logger,
	)
	if err != nil {
		logger.Warning("Redis连接失败: %v，系统将降级为仅使用文件存储", err)
		redisCache = nil
	} else {
		defer redisCache.Close()
	}

	// 初始化车辆服务
	carService := models.NewCarService(fileRepo, logger, redisCache)

	// 初始化车辆控制器
	carController := controllers.NewCarController(carService, logger)

	// 初始化Gin路由
	gin.SetMode(appConfig.GinMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(logger))

	// 设置受信任的代理
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// 配置CORS
	r.Use(cors.New(config.GetCorsConfig()))

	// API路由
	api := r.Group("/api")
	carController.RegisterRoutes(api)

	// 启动服务器
	logger.Info("服务器启动在 http://localhost:%d", appConfig.ServerPort)
	if err := r.Run(fmt.Sprintf(":%d", appConfig.ServerPort)); err != nil {
		logger.Fatal("服务器启动失败: %v", err)
	}
}