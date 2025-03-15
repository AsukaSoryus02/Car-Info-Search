package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonzheng/carrag/models"
	"github.com/jasonzheng/carrag/utils"
)

// CarController 车辆控制器
type CarController struct {
	CarService *models.CarService // 车辆服务
	Logger     *utils.Logger      // 日志记录器
}

// NewCarController 创建新的车辆控制器
func NewCarController(carService *models.CarService, logger *utils.Logger) *CarController {
	return &CarController{
		CarService: carService,
		Logger:     logger,
	}
}

// RegisterRoutes 注册路由
func (c *CarController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/cars", c.GetCars)
	router.GET("/cars/:id", c.GetCarByID)
	router.POST("/cars", c.CreateCar)
	router.PUT("/cars/:id", c.UpdateCar)
	router.DELETE("/cars/:id", c.DeleteCar)
	router.GET("/cars/brand/:brand", c.GetCarsByBrand)
}

// GetCars 获取所有车辆信息
func (c *CarController) GetCars(ctx *gin.Context) {
	// 使用服务层获取所有车辆信息
	cars, err := c.CarService.GetAllCars()
	if err != nil {
		c.Logger.Error("获取车辆信息失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取车辆信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, cars)
}

// GetCarByID 根据ID获取车辆信息
func (c *CarController) GetCarByID(ctx *gin.Context) {
	id := ctx.Param("id")

	// 使用服务层获取车辆信息
	car, err := c.CarService.GetCarByID(id)
	if err != nil {
		c.Logger.Warning("获取车辆信息失败: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "车辆信息不存在"})
		return
	}

	ctx.JSON(http.StatusOK, car)
}

// CreateCar 创建车辆信息
func (c *CarController) CreateCar(ctx *gin.Context) {
	var car models.Car

	// 解析请求体
	if err := ctx.ShouldBindJSON(&car); err != nil {
		c.Logger.Warning("解析请求体失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 生成唯一ID和设置创建时间会在服务层处理
	// 不需要在控制器中设置

	// 使用服务层创建车辆信息
	if err := c.CarService.CreateCar(&car); err != nil {
		c.Logger.Error("创建车辆信息失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建车辆信息失败"})
		return
	}

	ctx.JSON(http.StatusCreated, car)
}

// UpdateCar 更新车辆信息
func (c *CarController) UpdateCar(ctx *gin.Context) {
	id := ctx.Param("id")
	var car models.Car

	// 解析请求体
	if err := ctx.ShouldBindJSON(&car); err != nil {
		c.Logger.Warning("解析请求体失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 确保ID一致
	car.ID = id

	// 更新时间会在服务层设置

	// 使用服务层更新车辆信息
	if err := c.CarService.UpdateCar(&car); err != nil {
		c.Logger.Error("更新车辆信息失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新车辆信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, car)
}

// DeleteCar 删除车辆信息
func (c *CarController) DeleteCar(ctx *gin.Context) {
	id := ctx.Param("id")

	// 使用服务层删除车辆信息
	if err := c.CarService.DeleteCar(id); err != nil {
		c.Logger.Error("删除车辆信息失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除车辆信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "车辆信息已删除"})
}

// GetCarsByBrand 根据品牌获取车辆信息
func (c *CarController) GetCarsByBrand(ctx *gin.Context) {
	brand := ctx.Param("brand")

	// 使用服务层获取车辆信息
	cars, err := c.CarService.FindCarsByBrand(brand)
	if err != nil {
		c.Logger.Error("获取车辆信息失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取车辆信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, cars)
}
