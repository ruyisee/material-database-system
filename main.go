package main

import (
	"context"
	"embed"

	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化日志
	if err := InitLogger(); err != nil {
		logrus.WithError(err).Fatal("初始化日志失败")
	}
	// 初始化数据库
	if err := InitDB(); err != nil {
		logrus.WithError(err).Fatal("初始化数据库失败")
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	w_err := wails.Run(&options.App{
		Title:  "物料管理助手",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			logrus.Info("应用程序启动")
			app.startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			logrus.Info("应用程序关闭")
		},
		Bind: []interface{}{
			app,
		},
	})

	if w_err != nil {
		logrus.WithError(w_err).Error("应用启动失败")
		println("Error:", w_err.Error())
	}
}
