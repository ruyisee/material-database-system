package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// 根据规格和材质查询物料
func (a *App) GoSearchComponents(outerDiameter string, wallThickness string, material string, code string) ([]Component, error) {
	logrus.Info("查询物料: ", outerDiameter, wallThickness, material, code)
	materialService := NewComponentService()
	components, err := materialService.SearchComponents(outerDiameter, wallThickness, material, code)
	if err != nil {
		logrus.WithError(err).Warn("查询物料失败")
		return nil, err
	}
	return components, nil
}

// 添加物料
func (a *App) GoAddComponent(outerDiameter string, wallThickness string, material string, code string) error {
	logrus.Info("添加物料: ", outerDiameter, wallThickness, material, code)
	materialService := NewComponentService()
	err := materialService.AddComponent(&Component{
		OuterDiameter: outerDiameter,
		WallThickness: wallThickness,
		Material:      material,
		Code:          code,
	})
	if err != nil {
		logrus.WithError(err).Warn("添加物料失败")
		return err
	}
	return nil
}

// 添加材质编码
func (a *App) GoAddMaterial(code string, remark string) error {
	logrus.Info("添加材质: ", code, remark)
	materialService := NewComponentService()
	err := materialService.AddMaterial(&Material{
		Code:   code,
		Remark: remark,
	})
	if err != nil {
		logrus.WithError(err).Warn("添加材质失败")
		return err
	}
	return nil
}

// 获取所有材质
func (a *App) GoGetAllMaterial() ([]Material, error) {
	logrus.Info("查询所有材质")
	materialService := NewComponentService()
	materials, err := materialService.GetAllMaterial()
	if err != nil {
		logrus.WithError(err).Warn("查询材质失败")
		return nil, err
	}
	return materials, nil
}
