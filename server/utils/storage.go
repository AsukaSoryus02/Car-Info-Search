package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Storage 文件存储管理器
type Storage struct {
	DataDir  string       // 数据目录路径
	FileLock sync.RWMutex // 读写锁，用于并发控制
	Logger   *Logger      // 日志记录器
}

// NewStorage 创建新的存储管理器
func NewStorage(dataDir string, logger *Logger) (*Storage, error) {
	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	return &Storage{
		DataDir:  dataDir,
		FileLock: sync.RWMutex{},
		Logger:   logger,
	}, nil
}

// SaveJSON 将数据保存为JSON文件
func (s *Storage) SaveJSON(filename string, data interface{}) error {
	// 获取写锁
	s.FileLock.Lock()
	defer s.FileLock.Unlock()

	// 构建完整文件路径
	filePath := filepath.Join(s.DataDir, filename)

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		s.Logger.Error("创建目录失败: %v", err)
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 将数据转换为格式化的JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		s.Logger.Error("序列化数据失败: %v", err)
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	// 如果文件已存在，创建备份
	if s.FileExists(filename) {
		backupFile := filePath + ".bak"
		if err := os.Rename(filePath, backupFile); err != nil {
			s.Logger.Warning("创建备份文件失败: %v，继续保存新文件", err)
		} else {
			s.Logger.Debug("已创建备份文件: %s", backupFile)
		}
	}

	// 先写入临时文件
	tempFile := filePath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		s.Logger.Error("写入临时文件失败: %v", err)
		return fmt.Errorf("写入临时文件失败: %w", err)
	}

	// 重命名临时文件为目标文件（原子操作）
	if err := os.Rename(tempFile, filePath); err != nil {
		// 重命名失败时尝试删除临时文件
		os.Remove(tempFile)
		s.Logger.Error("重命名文件失败: %v", err)
		return fmt.Errorf("重命名文件失败: %w", err)
	}

	s.Logger.Info("成功保存数据到文件: %s", filename)
	return nil
}

// LoadJSON 从JSON文件加载数据
func (s *Storage) LoadJSON(filename string, target interface{}) error {
	// 获取读锁
	s.FileLock.RLock()
	defer s.FileLock.RUnlock()

	// 构建完整文件路径
	filePath := filepath.Join(s.DataDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.Logger.Warning("文件不存在: %s", filename)
		return nil // 文件不存在不视为错误，由调用者处理空数据情况
	}

	// 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		s.Logger.Error("读取文件失败: %v", err)
		return fmt.Errorf("读取文件失败: %w", err)
	}

	// 文件为空时直接返回
	if len(fileData) == 0 {
		s.Logger.Warning("文件为空: %s", filename)
		return nil
	}

	// 尝试从备份文件恢复
	if err := json.Unmarshal(fileData, target); err != nil {
		s.Logger.Error("解析JSON数据失败: %v，尝试从备份恢复", err)

		// 检查是否存在备份文件
		backupFile := filePath + ".bak"
		if _, backupErr := os.Stat(backupFile); !os.IsNotExist(backupErr) {
			// 读取备份文件
			backupData, backupErr := os.ReadFile(backupFile)
			if backupErr == nil && len(backupData) > 0 {
				// 尝试解析备份数据
				if backupErr := json.Unmarshal(backupData, target); backupErr == nil {
					s.Logger.Info("成功从备份文件恢复数据: %s", backupFile)
					return nil
				}
			}
		}

		return fmt.Errorf("解析JSON数据失败: %w", err)
	}

	s.Logger.Info("成功从文件加载数据: %s", filename)
	return nil
}

// FileExists 检查文件是否存在
func (s *Storage) FileExists(filename string) bool {
	filePath := filepath.Join(s.DataDir, filename)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// DeleteFile 删除文件
func (s *Storage) DeleteFile(filename string) error {
	// 获取写锁
	s.FileLock.Lock()
	defer s.FileLock.Unlock()

	filePath := filepath.Join(s.DataDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.Logger.Warning("尝试删除不存在的文件: %s", filename)
		return nil // 文件不存在不视为错误
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	s.Logger.Info("成功删除文件: %s", filename)
	return nil
}
