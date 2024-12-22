package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"encoding/csv"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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
	OuterDiameter string `gorm:"not null;uniqueIndex:Idx_outer_diameter_wall_thickness_material;comment:外径" json:"outerDiameter"`
	// 壁厚
	WallThickness string `gorm:"not null;uniqueIndex:Idx_outer_diameter_wall_thickness_material;comment:壁厚" json:"wallThickness"`
	// 材质
	Material  string    `gorm:"not null;uniqueIndex:Idx_outer_diameter_wall_thickness_material;comment:材质" json:"material"` // 材质(如 304)
	Code      string    `gorm:"uniqueIndex;not null;comment:物料码" json:"code"`                                               // 物料码
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                                                           // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`                                                           // 更新时间
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
	return s.db.Model(&Component{}).Where("id = ?", m.ID).Updates(m).Error
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
			logrus.WithError(err).Warn("查询物料失败")
			return nil, err
		}
		return []Component{*component}, nil
	}

	if outerDiameter != "" {
		selectSql += "outer_diameter = ?"
	}
	if wallThickness != "" {
		if selectSql != "" {
			selectSql += " AND wall_thickness = ?"
		} else {
			selectSql += "wall_thickness = ?"
		}
	}
	if material != "" {
		if selectSql != "" {
			selectSql += " AND material = ?"
		} else {
			selectSql += "material = ?"
		}
	}
	if selectSql != "" {
		err := s.db.Where(selectSql, outerDiameter, wallThickness, material).Find(&components).Error
		return components, err
	} else {
		logrus.Info("查询所有物料")
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

// 检测文件编码
func detectFileEncoding(file string) (string, error) {
	// 读取文件前几个字节来判断编码
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 读取前4个字节
	bom := make([]byte, 4)
	n, err := f.Read(bom)
	if err != nil {
		return "", err
	}

	// 重置文件指针
	f.Seek(0, 0)

	// 检查BOM标记
	if n >= 3 && bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
		return "UTF-8", nil
	}
	if n >= 2 && bom[0] == 0xFF && bom[1] == 0xFE {
		return "UTF-16LE", nil
	}
	if n >= 2 && bom[0] == 0xFE && bom[1] == 0xFF {
		return "UTF-16BE", nil
	}

	// 如果没有BOM，读取更多内容进行分析
	content := make([]byte, 1024)
	n, err = f.Read(content)
	if err != nil {
		return "", err
	}
	content = content[:n]

	// 检查是否符合UTF-8编码规则
	if isUTF8(content) {
		return "UTF-8", nil
	}

	// 假设是GBK
	return "GBK", nil
}

// 判断是否是有效的UTF-8编码
func isUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { // 非ASCII字符
				for (buf[i] & 0x80) != 0 {
					buf[i] <<= 1
					nBytes++
				}
				nBytes--         // 减去首字节本身
				if nBytes == 0 { // 非法编码
					return false
				}
			}
		} else { // 多字节字符的非首字节
			if buf[i]&0xC0 != 0x80 {
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

// 根据编码读取CSV文件
func readCSVWithEncoding(file string) ([][]string, error) {
	// 检测文件编码
	encoding, err := detectFileEncoding(file)
	if err != nil {
		return nil, fmt.Errorf("检测文件编码失败: %v", err)
	}
	logrus.Infof("检测到文件编码: %s", encoding)

	// 打开文件
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	var reader *csv.Reader
	if encoding == "GBK" {
		decoder := simplifiedchinese.GBK.NewDecoder()
		reader = csv.NewReader(transform.NewReader(f, decoder))
	} else {
		reader = csv.NewReader(f)
	}

	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%s解码失败: %v", encoding, err)
	}

	return records, nil
}

// 导入物料
func (s *ComponentService) ImportComponents(file string) string {
	// 读取CSV文件
	records, err := readCSVWithEncoding(file)
	if err != nil {
		return "读取文件失败: " + err.Error()
	}

	if len(records) < 2 {
		return "CSV文件为空或只有表头"
	}

	// 解析表头，获取字段索引
	headers := records[0]
	fieldIndex := make(map[string]int)
	requiredFields := []string{"外径", "壁厚", "材质", "物料码"}

	logrus.Info("解析表头: ", headers)
	// 构建字段索引映射
	for i, header := range headers {
		header = strings.TrimSpace(header)
		fieldIndex[header] = i
	}

	// 验证必需字段是否存在
	var missingFields []string
	for _, field := range requiredFields {
		if _, exists := fieldIndex[field]; !exists {
			missingFields = append(missingFields, field)
		}
	}
	if len(missingFields) > 0 {
		return fmt.Sprintf("CSV文件缺少必需字段: %v", missingFields)
	}

	var meterialTypeSet = make(map[string]bool)
	var errorLinesErrorMap = make(map[int]string)
	var totalLines = len(records) - 1 // 减去表头

	// 从第二行开始处理数据
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < len(headers) {
			errorLinesErrorMap[i] = fmt.Sprintf("行 %d 格式错误: 列数不正确", i)
			continue
		}

		outerDiameter := strings.TrimSpace(record[fieldIndex["外径"]])
		wallThickness := strings.TrimSpace(record[fieldIndex["壁厚"]])
		material := strings.TrimSpace(record[fieldIndex["材质"]])
		code := strings.TrimSpace(record[fieldIndex["物料码"]])

		if outerDiameter == "" || wallThickness == "" || material == "" || code == "" {
			errorLinesErrorMap[i] = fmt.Sprintf("行 %d 存在空字段", i)
			continue
		}

		_err := s.AddComponent(&Component{
			OuterDiameter: outerDiameter,
			WallThickness: wallThickness,
			Material:      material,
			Code:          code,
		})
		if _err != nil {
			errorLinesErrorMap[i] = _err.Error()
		}
		meterialTypeSet[material] = true
	}

	// 添加材质
	for material := range meterialTypeSet {
		s.AddMaterial(&Material{
			Code:   material,
			Remark: "",
		})
	}

	// 构建错误信息
	var errorMsg string
	if len(errorLinesErrorMap) > 0 {
		errorMsg = "\n失败行详情:\n"
		for line, err := range errorLinesErrorMap {
			errorMsg += fmt.Sprintf("第 %d 行: %s\n", line, err)
		}
	}

	return fmt.Sprintf("导入完成\n总行数: %d\n成功: %d\n失败: %d%s",
		totalLines,
		totalLines-len(errorLinesErrorMap),
		len(errorLinesErrorMap),
		errorMsg)
}
