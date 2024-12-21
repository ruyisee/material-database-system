package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Component 定义组件模型
// φ22*5	304	81383220
// φ45*8	304	81349147
// φ127*28	16Mn	81386223

type Component struct {
	ID uint `gorm:"primarykey" json:"id"`
	// 外径
	OuterDiameter string `gorm:"not null;comment:外径" json:"outerDiameter"`
	// 壁厚
	WallThickness string `gorm:"not null;comment:壁厚" json:"wallThickness"`
	// 材质
	Material  string    `gorm:"not null;comment:材质" json:"material"`          // 材质(如 304)
	Code      string    `gorm:"uniqueIndex;not null;comment:物料码" json:"code"` // 物料码
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`             // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`             // 更新时间
}

// 材质编码
type Material struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Code      string    `gorm:"uniqueIndex;not null;comment:材质编码" json:"code"`
	Remark    string    `gorm:"comment:备注" json:"remark"`         // 备注
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
}

// DB 全局数据库连接
var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("material-database.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&Material{})
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&Component{})
	if err != nil {
		return err
	}
	return nil
}

// ComponentService 组件服务结构体
type ComponentService struct {
	db *gorm.DB
}

// NewMaterialService 创建物料服务实例
func NewComponentService() *ComponentService {
	return &ComponentService{db: DB}
}

// Create 创建新物料
func (s *ComponentService) AddComponent(m *Component) error {
	return s.db.Create(m).Error
}

// GetAll 获取所有物料
func (s *ComponentService) GetAllComponent() ([]Component, error) {
	var components []Component
	err := s.db.Find(&components).Error
	return components, err
}

// GetByCode 根据物料码查询
func (s *ComponentService) GetComponentByCode(code string) (*Component, error) {
	var component Component
	err := s.db.Where("code = ?", code).First(&component).Error
	if err != nil {
		return nil, err
	}
	return &component, nil
}

// UpdateComponent 更新组件信息
func (s *ComponentService) UpdateComponent(m *Component) error {
	return s.db.Save(m).Error
}

// DeleteComponent 删除组件
func (s *ComponentService) DeleteComponent(code string) error {
	return s.db.Where("code = ?", code).Delete(&Component{}).Error
}

// GetBySpecMaterial 根据规格和材质查询
func (s *ComponentService) GetBySpecMaterial(spec string, material string) ([]Component, error) {
	var components []Component
	err := s.db.Where("spec = ? AND material = ?", spec, material).Find(&components).Error
	return components, err
}

// GetByOuterDiameterWallThicknessMaterial 根据外径和壁厚查询
func (s *ComponentService) SearchComponents(outerDiameter string, wallThickness string, material string, code string) ([]Component, error) {
	var components []Component
	var selectSql string = ""
	// 如果物料码不为空，则根据物料码查询
	if code != "" {
		component, err := s.GetComponentByCode(code)
		if err != nil {
			Log.Error("--查询物料失败: ", err)
			return nil, err
		}
		return []Component{*component}, nil
	}

	if outerDiameter != "" {
		selectSql += "outer_diameter = ?"
	}
	if wallThickness != "" {
		selectSql += " AND wall_thickness = ?"
	}
	if material != "" {
		selectSql += " AND material = ?"
	}
	if selectSql != "" {
		err := s.db.Where(selectSql, outerDiameter, wallThickness, material).Find(&components).Error
		return components, err
	} else {
		Log.Info("查询所有物料")
		err := s.db.Find(&components).Error
		return components, err
	}
}

// 添加材质
func (s *ComponentService) AddMaterial(m *Material) error {
	return s.db.Create(m).Error
}

// 获取所有材质
func (s *ComponentService) GetAllMaterial() ([]Material, error) {
	var materials []Material
	err := s.db.Find(&materials).Error
	return materials, err
}
