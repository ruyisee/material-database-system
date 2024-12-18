package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Material 定义物料模型
// φ22*5	304	81383220
// φ45*8	304	81349147
// φ127*28	16Mn	81386223

type Material struct {
	ID uint `gorm:"primarykey" json:"id"`
	// 外径
	OuterDiameter string `gorm:"not null;comment:外径" json:"outerDiameter"`
	// 壁厚
	WallThickness string `gorm:"not null;comment:壁厚" json:"wallThickness"`
	// 材质
	Material  string    `gorm:"not null;comment:材质" json:"material"`          // 材质(如 304)
	Code      string    `gorm:"uniqueIndex;not null;comment:物料码" json:"code"` // 物料码
	CreatedAt time.Time `json:"created_at"`                                   // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                                   // 更新时间
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

	return nil
}

// MaterialService 物料服务结构体
type MaterialService struct {
	db *gorm.DB
}

// NewMaterialService 创建物料服务实例
func NewMaterialService() *MaterialService {
	return &MaterialService{db: DB}
}

// Create 创建新物料
func (s *MaterialService) Create(m *Material) error {
	return s.db.Create(m).Error
}

// GetAll 获取所有物料
func (s *MaterialService) GetAll() ([]Material, error) {
	var materials []Material
	err := s.db.Find(&materials).Error
	return materials, err
}

// GetByCode 根据物料码查询
func (s *MaterialService) GetByCode(code string) (*Material, error) {
	var material Material
	err := s.db.Where("code = ?", code).First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

// Update 更新物料信息
func (s *MaterialService) Update(m *Material) error {
	return s.db.Save(m).Error
}

// Delete 删除物料
func (s *MaterialService) Delete(code string) error {
	return s.db.Where("code = ?", code).Delete(&Material{}).Error
}

// Search 搜索物料
func (s *MaterialService) Search(query string) ([]Material, error) {
	var materials []Material
	err := s.db.Where("spec LIKE ? OR material LIKE ? OR code LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Find(&materials).Error
	return materials, err
}

// GetBySpecMaterial 根据规格和材质查询
func (s *MaterialService) GetBySpecMaterial(spec string, material string) ([]Material, error) {
	var materials []Material
	err := s.db.Where("spec = ? AND material = ?", spec, material).Find(&materials).Error
	return materials, err
}

// GetByOuterDiameterWallThicknessMaterial 根据外径和壁厚查询
func (s *MaterialService) GetByOuterDiameterWallThicknessMaterial(outerDiameter string, wallThickness string, material string) ([]Material, error) {
	var materials []Material
	if material != "" && outerDiameter != "" && wallThickness != "" {
		err := s.db.Where("outer_diameter = ? AND wall_thickness = ? AND material = ?", outerDiameter, wallThickness, material).Find(&materials).Error
		if err != nil {
			return nil, err
		}
	} else if outerDiameter != "" && wallThickness != "" {
		err := s.db.Where("outer_diameter = ? AND wall_thickness = ?", outerDiameter, wallThickness).Find(&materials).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := s.db.Find(&materials).Error
		if err != nil {
			return nil, err
		}
	}
	return materials, nil
}
