package repositories

import (
	"fmt"

	"github.com/jasonzheng/carrag/models"
	"github.com/jasonzheng/carrag/utils"
)

// FileCarRepository 基于文件的车辆信息仓库实现
type FileCarRepository struct {
	Storage  *utils.Storage // 文件存储管理器
	Logger   *utils.Logger  // 日志记录器
	FileName string         // 数据文件名
}

// NewFileCarRepository 创建新的文件车辆信息仓库
func NewFileCarRepository(storage *utils.Storage, logger *utils.Logger, fileName string) *FileCarRepository {
	return &FileCarRepository{
		Storage:  storage,
		Logger:   logger,
		FileName: fileName,
	}
}

// FindAll 获取所有车辆信息
func (r *FileCarRepository) FindAll() ([]models.Car, error) {
	r.Logger.Debug("从文件加载所有车辆信息: %s", r.FileName)
	var cars []models.Car
	err := r.Storage.LoadJSON(r.FileName, &cars)
	if err != nil {
		r.Logger.Error("加载车辆信息失败: %v", err)
		return nil, fmt.Errorf("加载车辆信息失败: %w", err)
	}

	// 如果文件不存在或为空，返回空数组
	if cars == nil {
		cars = []models.Car{}
	}

	r.Logger.Debug("成功加载 %d 条车辆信息", len(cars))
	return cars, nil
}

// FindByID 根据ID获取车辆信息
func (r *FileCarRepository) FindByID(id string) (*models.Car, error) {
	r.Logger.Debug("根据ID查找车辆信息: %s", id)
	cars, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	for _, car := range cars {
		if car.ID == id {
			r.Logger.Debug("找到车辆信息: %s", id)
			return &car, nil
		}
	}

	r.Logger.Warning("未找到车辆信息: %s", id)
	return nil, fmt.Errorf("车辆信息不存在: %s", id)
}

// Create 创建车辆信息
func (r *FileCarRepository) Create(car *models.Car) error {
	r.Logger.Debug("创建车辆信息: %s %s", car.Brand, car.Model)
	cars, err := r.FindAll()
	if err != nil {
		return err
	}

	// 添加新车辆
	cars = append(cars, *car)

	// 保存到文件
	if err := r.Storage.SaveJSON(r.FileName, cars); err != nil {
		r.Logger.Error("保存车辆信息失败: %v", err)
		return fmt.Errorf("保存车辆信息失败: %w", err)
	}

	r.Logger.Debug("成功创建车辆信息: %s", car.ID)
	return nil
}

// Update 更新车辆信息
func (r *FileCarRepository) Update(car *models.Car) error {
	r.Logger.Debug("更新车辆信息: %s", car.ID)
	cars, err := r.FindAll()
	if err != nil {
		return err
	}

	// 查找并更新
	found := false
	for i, c := range cars {
		if c.ID == car.ID {
			cars[i] = *car
			found = true
			break
		}
	}

	if !found {
		r.Logger.Warning("未找到要更新的车辆信息: %s", car.ID)
		return fmt.Errorf("车辆信息不存在: %s", car.ID)
	}

	// 保存到文件
	if err := r.Storage.SaveJSON(r.FileName, cars); err != nil {
		r.Logger.Error("保存车辆信息失败: %v", err)
		return fmt.Errorf("保存车辆信息失败: %w", err)
	}

	r.Logger.Debug("成功更新车辆信息: %s", car.ID)
	return nil
}

// Delete 删除车辆信息
func (r *FileCarRepository) Delete(id string) error {
	r.Logger.Debug("删除车辆信息: %s", id)
	cars, err := r.FindAll()
	if err != nil {
		return err
	}

	// 查找并删除
	found := false
	newCars := make([]models.Car, 0, len(cars))
	for _, car := range cars {
		if car.ID != id {
			newCars = append(newCars, car)
		} else {
			found = true
		}
	}

	if !found {
		r.Logger.Warning("未找到要删除的车辆信息: %s", id)
		return fmt.Errorf("车辆信息不存在: %s", id)
	}

	// 保存到文件
	if err := r.Storage.SaveJSON(r.FileName, newCars); err != nil {
		r.Logger.Error("保存车辆信息失败: %v", err)
		return fmt.Errorf("保存车辆信息失败: %w", err)
	}

	r.Logger.Debug("成功删除车辆信息: %s", id)
	return nil
}

// FindByBrand 根据品牌查找车辆信息
func (r *FileCarRepository) FindByBrand(brand string) ([]models.Car, error) {
	r.Logger.Debug("根据品牌查找车辆信息: %s", brand)
	cars, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	// 筛选符合条件的车辆
	result := make([]models.Car, 0)
	for _, car := range cars {
		if car.Brand == brand {
			result = append(result, car)
		}
	}

	r.Logger.Debug("找到 %d 条符合品牌 %s 的车辆信息", len(result), brand)
	return result, nil
}
