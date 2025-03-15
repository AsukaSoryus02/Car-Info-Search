package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache Redis缓存管理器
type RedisCache struct {
	Client *redis.Client // Redis客户端
	Logger *Logger       // 日志记录器
	Prefix string        // 键前缀，用于区分不同应用的缓存
}

// NewRedisCache 创建新的Redis缓存管理器
func NewRedisCache(addr string, password string, db int, prefix string, logger *Logger) (*RedisCache, error) {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         addr,            // Redis服务器地址
		Password:     password,        // Redis密码
		DB:           db,              // 使用的数据库编号
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读取超时
		WriteTimeout: 3 * time.Second, // 写入超时
		PoolSize:     10,              // 连接池大小
		MinIdleConns: 2,               // 最小空闲连接数
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	logger.Info("Redis连接成功: %s", pong)

	return &RedisCache{
		Client: client,
		Logger: logger,
		Prefix: prefix,
	}, nil
}

// Close 关闭Redis连接
func (r *RedisCache) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}

// formatKey 格式化缓存键名
func (r *RedisCache) formatKey(key string) string {
	return fmt.Sprintf("%s:%s", r.Prefix, key)
}

// Set 设置缓存
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		r.Logger.Error("序列化缓存数据失败: %v", err)
		return fmt.Errorf("序列化缓存数据失败: %w", err)
	}

	// 设置缓存
	formattedKey := r.formatKey(key)
	err = r.Client.Set(ctx, formattedKey, data, expiration).Err()
	if err != nil {
		r.Logger.Error("设置缓存失败: %v", err)
		return fmt.Errorf("设置缓存失败: %w", err)
	}

	r.Logger.Debug("成功设置缓存: %s, 过期时间: %v", formattedKey, expiration)
	return nil
}

// Get 获取缓存
func (r *RedisCache) Get(ctx context.Context, key string, target interface{}) error {
	// 获取缓存
	formattedKey := r.formatKey(key)
	data, err := r.Client.Get(ctx, formattedKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			// 键不存在
			r.Logger.Debug("缓存不存在: %s", formattedKey)
			return fmt.Errorf("缓存不存在: %w", err)
		}
		r.Logger.Error("获取缓存失败: %v", err)
		return fmt.Errorf("获取缓存失败: %w", err)
	}

	// 反序列化
	if err := json.Unmarshal(data, target); err != nil {
		r.Logger.Error("解析缓存数据失败: %v", err)
		return fmt.Errorf("解析缓存数据失败: %w", err)
	}

	r.Logger.Debug("成功获取缓存: %s", formattedKey)
	return nil
}

// Delete 删除缓存
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	formattedKey := r.formatKey(key)
	err := r.Client.Del(ctx, formattedKey).Err()
	if err != nil {
		r.Logger.Error("删除缓存失败: %v", err)
		return fmt.Errorf("删除缓存失败: %w", err)
	}

	r.Logger.Debug("成功删除缓存: %s", formattedKey)
	return nil
}

// Exists 检查缓存是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	formattedKey := r.formatKey(key)
	val, err := r.Client.Exists(ctx, formattedKey).Result()
	if err != nil {
		r.Logger.Error("检查缓存是否存在失败: %v", err)
		return false, fmt.Errorf("检查缓存是否存在失败: %w", err)
	}

	exists := val > 0
	r.Logger.Debug("检查缓存是否存在: %s, 结果: %v", formattedKey, exists)
	return exists, nil
}

// SetWithRetry 设置缓存，失败时重试
func (r *RedisCache) SetWithRetry(ctx context.Context, key string, value interface{}, expiration time.Duration, retries int, retryDelay time.Duration) error {
	var lastErr error

	for i := 0; i <= retries; i++ {
		err := r.Set(ctx, key, value, expiration)
		if err == nil {
			return nil
		}

		lastErr = err
		if i < retries {
			r.Logger.Warning("设置缓存失败，将在 %v 后重试: %v", retryDelay, err)
			time.Sleep(retryDelay)
		}
	}

	return lastErr
}

// GetWithFallback 获取缓存，失败时使用回退函数获取数据并更新缓存
func (r *RedisCache) GetWithFallback(ctx context.Context, key string, target interface{}, fallback func() (interface{}, error), expiration time.Duration) error {
	// 尝试从缓存获取
	err := r.Get(ctx, key, target)
	if err == nil {
		// 缓存命中
		return nil
	}

	// 缓存未命中，使用回退函数获取数据
	data, err := fallback()
	if err != nil {
		r.Logger.Error("回退函数执行失败: %v", err)
		return fmt.Errorf("回退函数执行失败: %w", err)
	}

	// 更新目标对象
	dataBytes, err := json.Marshal(data)
	if err != nil {
		r.Logger.Error("序列化回退数据失败: %v", err)
		return fmt.Errorf("序列化回退数据失败: %w", err)
	}

	if err := json.Unmarshal(dataBytes, target); err != nil {
		r.Logger.Error("解析回退数据失败: %v", err)
		return fmt.Errorf("解析回退数据失败: %w", err)
	}

	// 异步更新缓存
	go func() {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.Set(ctxTimeout, key, data, expiration); err != nil {
			r.Logger.Error("异步更新缓存失败: %v", err)
		}
	}()

	return nil
}
