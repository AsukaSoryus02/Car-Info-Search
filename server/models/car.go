package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jasonzheng/carrag/utils"
)

// Car 车辆信息模型
type Car struct {
	ID                 string    `json:"id"`                           // 车辆唯一标识符
	Brand              string    `json:"brand"`                        // 品牌
	Model              string    `json:"model"`                        // 车型
	FuelConsumption    float64   `json:"fuelConsumption,omitempty"`    // 油耗(L/100km)
	FuelType           string    `json:"fuelType,omitempty"`           // 燃油类型
	Mileage            float64   `json:"mileage,omitempty"`            // 行驶里程(km)
	AnnualMileage      float64   `json:"annualMileage,omitempty"`      // 年均行驶里程(km)
	StorageEnvironment string    `json:"storageEnvironment,omitempty"` // 存放环境
	UsageScenario      []string  `json:"usageScenario,omitempty"`      // 使用场景
	Remarks            string    `json:"remarks,omitempty"`            // 备注
	CreatedAt          time.Time `json:"createdAt"`                    // 创建时间
	UpdatedAt          time.Time `json:"updatedAt,omitempty"`          // 更新时间
}

// CarRepository 车辆信息仓库接口
type CarRepository interface {
	FindAll() ([]Car, error)                 // 获取所有车辆信息
	FindByID(id string) (*Car, error)        // 根据ID获取车辆信息
	Create(car *Car) error                   // 创建车辆信息
	Update(car *Car) error                   // 更新车辆信息
	Delete(id string) error                  // 删除车辆信息
	FindByBrand(brand string) ([]Car, error) // 根据品牌查找车辆信息
}

// CarService 车辆信息服务
type CarService struct {
	Repo   CarRepository     // 车辆信息仓库
	Logger *utils.Logger     // 日志记录器
	Cache  *utils.RedisCache // Redis缓存
}

// NewCarService 创建车辆信息服务
func NewCarService(repo CarRepository, logger *utils.Logger, cache *utils.RedisCache) *CarService {
	return &CarService{
		Repo:   repo,
		Logger: logger,
		Cache:  cache,
	}
}

// GenerateID 生成唯一的车辆标识符
func GenerateID() string {
	// 使用UUID生成基础唯一ID
	uuid := uuid.New().String()
	// 取UUID前8位作为ID前缀
	shortID := uuid[:8]
	// 添加时间戳（base36格式）确保即使UUID重复也能保持唯一性
	timestamp := time.Now().UnixNano()
	timestampStr := utils.FormatInt36(timestamp)[:6]
	// 组合ID格式：UUID前缀-时间戳
	return shortID + "-" + timestampStr
}

// GetAllCars 获取所有车辆信息
func (s *CarService) GetAllCars() ([]Car, error) {
	s.Logger.Info("获取所有车辆信息")
	return s.Repo.FindAll()
}

// GetCarByID 根据ID获取车辆信息
func (s *CarService) GetCarByID(id string) (*Car, error) {
	s.Logger.Info("获取车辆信息，ID: %s", id)

	// 定义一个空的Car指针用于存储结果
	var car Car

	// 如果缓存可用，尝试从缓存获取
	if s.Cache != nil {
		// 创建上下文
		ctx := s.Cache.Client.Context()

		// 定义回退函数，从数据库获取数据
		fallback := func() (interface{}, error) {
			result, err := s.Repo.FindByID(id)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		// 使用缓存获取，如果缓存不存在则使用回退函数获取并更新缓存
		err := s.Cache.GetWithFallback(ctx, "car:"+id, &car, fallback, 24*time.Hour)
		if err != nil {
			s.Logger.Error("获取车辆信息失败: %v", err)
			return nil, err
		}

		return &car, nil
	}

	// 缓存不可用，直接从数据库获取
	return s.Repo.FindByID(id)
}

// CreateCar 创建车辆信息
func (s *CarService) CreateCar(car *Car) error {
	s.Logger.Info("创建车辆信息: %s %s", car.Brand, car.Model)

	// 设置ID和时间
	car.ID = GenerateID()
	car.CreatedAt = time.Now()

	// 保存到数据库
	err := s.Repo.Create(car)
	if err != nil {
		s.Logger.Error("创建车辆信息失败: %v", err)
		return err
	}

	// 如果缓存可用，保存到缓存
	if s.Cache != nil {
		ctx := s.Cache.Client.Context()
		if err := s.Cache.Set(ctx, "car:"+car.ID, car, 24*time.Hour); err != nil {
			s.Logger.Warning("缓存车辆信息失败: %v", err)
			// 缓存失败不影响正常返回
		}
	}

	return nil
}

// UpdateCar 更新车辆信息
func (s *CarService) UpdateCar(car *Car) error {
	s.Logger.Info("更新车辆信息，ID: %s", car.ID)

	// 设置更新时间
	car.UpdatedAt = time.Now()

	// 更新数据库
	err := s.Repo.Update(car)
	if err != nil {
		s.Logger.Error("更新车辆信息失败: %v", err)
		return err
	}

	// 如果缓存可用，更新缓存
	if s.Cache != nil {
		ctx := s.Cache.Client.Context()
		if err := s.Cache.Set(ctx, "car:"+car.ID, car, 24*time.Hour); err != nil {
			s.Logger.Warning("更新缓存失败: %v", err)
			// 缓存失败不影响正常返回
		}
	}

	return nil
}

// DeleteCar 删除车辆信息
func (s *CarService) DeleteCar(id string) error {
	s.Logger.Info("删除车辆信息，ID: %s", id)

	// 从数据库删除
	err := s.Repo.Delete(id)
	if err != nil {
		s.Logger.Error("删除车辆信息失败: %v", err)
		return err
	}

	// 如果缓存可用，从缓存删除
	if s.Cache != nil {
		ctx := s.Cache.Client.Context()
		if err := s.Cache.Delete(ctx, "car:"+id); err != nil {
			s.Logger.Warning("从缓存删除失败: %v", err)
			// 缓存失败不影响正常返回
		}
	}

	return nil
}

// FindCarsByBrand 根据品牌查找车辆信息
func (s *CarService) FindCarsByBrand(brand string) ([]Car, error) {
	s.Logger.Info("根据品牌查找车辆信息: %s", brand)
	return s.Repo.FindByBrand(brand)
}
