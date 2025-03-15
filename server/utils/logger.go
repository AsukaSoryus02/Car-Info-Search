package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	// DEBUG 调试级别
	DEBUG LogLevel = iota
	// INFO 信息级别
	INFO
	// WARNING 警告级别
	WARNING
	// ERROR 错误级别
	ERROR
	// FATAL 致命错误级别
	FATAL
)

var levelNames = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

// Logger 自定义日志记录器
type Logger struct {
	logFile    *os.File
	debugLog   *log.Logger
	infoLog    *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
	fatalLog   *log.Logger
	minLevel   LogLevel
}

// NewLogger 创建一个新的日志记录器
func NewLogger(logDir string, minLevel LogLevel) (*Logger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 创建日志文件，使用当前日期作为文件名
	logFileName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}

	// 创建多输出目标，同时输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 创建不同级别的日志记录器
	debugLog := log.New(multiWriter, "DEBUG: ", log.Ldate|log.Ltime)
	infoLog := log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime)
	warningLog := log.New(multiWriter, "WARNING: ", log.Ldate|log.Ltime)
	errorLog := log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime)
	fatalLog := log.New(multiWriter, "FATAL: ", log.Ldate|log.Ltime)

	return &Logger{
		logFile:    logFile,
		debugLog:   debugLog,
		infoLog:    infoLog,
		warningLog: warningLog,
		errorLog:   errorLog,
		fatalLog:   fatalLog,
		minLevel:   minLevel,
	}, nil
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// getCallerInfo 获取调用者信息
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(3) // 跳过两层调用栈
	if !ok {
		return ""
	}
	// 只获取文件名，不包含路径
	shortFile := filepath.Base(file)
	return fmt.Sprintf("%s:%d", shortFile, line)
}

// formatMessage 格式化日志消息，添加调用者信息
func formatMessage(format string, args ...interface{}) string {
	callerInfo := getCallerInfo()
	message := fmt.Sprintf(format, args...)
	if callerInfo != "" {
		return fmt.Sprintf("%s - %s", callerInfo, message)
	}
	return message
}

// Debug 记录调试级别日志
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.minLevel <= DEBUG {
		l.debugLog.Printf(formatMessage(format, args...))
	}
}

// Info 记录信息级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	if l.minLevel <= INFO {
		l.infoLog.Printf(formatMessage(format, args...))
	}
}

// Warning 记录警告级别日志
func (l *Logger) Warning(format string, args ...interface{}) {
	if l.minLevel <= WARNING {
		l.warningLog.Printf(formatMessage(format, args...))
	}
}

// Error 记录错误级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	if l.minLevel <= ERROR {
		l.errorLog.Printf(formatMessage(format, args...))
	}
}

// Fatal 记录致命错误级别日志并退出程序
func (l *Logger) Fatal(format string, args ...interface{}) {
	if l.minLevel <= FATAL {
		l.fatalLog.Printf(formatMessage(format, args...))
		os.Exit(1)
	}
}
