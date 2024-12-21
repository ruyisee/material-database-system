package main

import (
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// InitLogger 初始化日志配置
func InitLogger() {
	// 创建 log 目录
	logPath := "logs"
	if err := os.MkdirAll(logPath, 0777); err != nil {
		panic(err)
	}

	// 设置日志文件
	filename := path.Join(logPath, "log.log")

	// 配置 rotatelogs
	writer, err := rotatelogs.New(
		filename+".%Y%m%d%H%M",                    // 日志文件名格式
		rotatelogs.WithLinkName(filename),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(30*24*time.Hour),    // 最多保存30天
		rotatelogs.WithRotationSize(10*1024*1024), // 10MB分割一次
		rotatelogs.WithRotationCount(5),           // 保留5个文件
	)
	if err != nil {
		panic(err)
	}

	// 配置 logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	logrus.SetOutput(writer)
	logrus.SetLevel(logrus.InfoLevel)
}

var Log = logrus.New()
