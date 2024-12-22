package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

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

// 删除物料
func (a *App) GoDeleteComponent(code string) error {
	logrus.Info("删除物料: ", code)
	materialService := NewComponentService()
	err := materialService.DeleteComponent(code)
	if err != nil {
		logrus.WithError(err).Warn("删除物料失败")
		return err
	}
	return nil
}

// 编辑物料
func (a *App) GoEditComponent(id int, code string, outerDiameter string, wallThickness string, material string) error {
	logrus.Info("编辑物料: ", code, outerDiameter, wallThickness, material)
	materialService := NewComponentService()
	err := materialService.UpdateComponent(&Component{
		ID:            uint(id),
		Code:          code,
		OuterDiameter: outerDiameter,
		WallThickness: wallThickness,
		Material:      material,
	})
	if err != nil {
		logrus.WithError(err).Warn("编辑物料失败")
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

// 导入物料
func (a *App) GoImportComponents(file string) string {
	logrus.Info("导入物料: ", file)
	materialService := NewComponentService()
	result := materialService.ImportComponents(file)
	logrus.Warn("导入物料情况汇总: ", result)
	return result
}

// 处理文件上传
func (a *App) GoUploadFile(fileContent string) (string, error) {
	// 创建临时目录
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "import.csv")

	// 将 Base64 字符串转换为字节
	fileData, err := base64.StdEncoding.DecodeString(fileContent)
	if err != nil {
		logrus.WithError(err).Error("解码文件内容失败")
		return "", err
	}

	// 写入文件
	err = os.WriteFile(tempFile, fileData, 0644)
	if err != nil {
		logrus.WithError(err).Error("写入文件失败")
		return "", err
	}

	return tempFile, nil
}
