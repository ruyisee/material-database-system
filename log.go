package main

import (
	"os"
	"path/filepath"
	"runtime"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func InitLogger() error {
	// 获取应用程序的日志目录
	logDir := getLogDir()

	// 创建日志目录
	if err := os.MkdirAll(logDir, 0777); err != nil {
		return err
	}

	// 设置日志文件路径
	logPath := filepath.Join(logDir, "app.log")

	// 配置 rotatelogs
	writer, err := rotatelogs.New(
		logPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithRotationSize(10*1024*1024),
		rotatelogs.WithRotationCount(5),
	)
	if err != nil {
		return err
	}

	// 配置 logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	logrus.SetOutput(writer)
	logrus.SetLevel(logrus.InfoLevel)

	// 记录初始化成功日志
	logrus.WithField("logPath", logPath).Info("日志系统初始化成功")

	return nil
}

// getLogDir 根据操作系统返回合适的日志目录
func getLogDir() string {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%/material-database-system/logs
		baseDir = os.Getenv("APPDATA")
	case "darwin":
		// macOS: ~/Library/Logs/material-database-system
		homeDir, _ := os.UserHomeDir()
		baseDir = filepath.Join(homeDir, "Library", "Logs")
	default:
		// Linux: ~/.local/share/material-database-system/logs
		homeDir, _ := os.UserHomeDir()
		baseDir = filepath.Join(homeDir, ".local", "share")
	}

	return filepath.Join(baseDir, "material-database-system", "logs")
}
