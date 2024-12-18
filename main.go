package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化数据库
	if err := InitDB(); err != nil {
		log.Fatal(err)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	w_err := wails.Run(&options.App{
		Title:  "material-database-system",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if w_err != nil {
		println("Error:", w_err.Error())
	}
}
