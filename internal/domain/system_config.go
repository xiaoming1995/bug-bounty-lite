package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// SystemConfig 系统配置实体
type SystemConfig struct {
	ID        uint      `gorm:"primaryKey;comment:配置ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	// 配置类型
	ConfigType string `gorm:"size:50;not null;index;comment:配置类型(vulnerability_type:漏洞类型/severity_level:危害等级/project_category:项目分类等)" json:"config_type"`

	// 配置键
	ConfigKey string `gorm:"size:100;not null;comment:配置键(如:SQL_INJECTION/XSS/CSRF等，用于程序内部识别)" json:"config_key"`

	// 配置值（显示名称）
	ConfigValue string `gorm:"size:255;not null;comment:配置值(显示名称，如:SQL注入/XSS跨站脚本，用于前端显示)" json:"config_value"`

	// 配置描述
	Description string `gorm:"type:text;comment:配置描述" json:"description"`

	// 排序顺序
	SortOrder int `gorm:"default:0;comment:排序顺序(数字越小越靠前)" json:"sort_order"`

	// 状态
	Status string `gorm:"size:20;default:'active';comment:配置状态(active:启用/inactive:禁用)" json:"status"`

	// 扩展数据（JSON格式）
	ExtraData JSON `gorm:"type:json;comment:扩展数据(JSON格式，存储额外信息如图标、颜色等)" json:"extra_data,omitempty"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

// JSON 类型用于存储 JSON 数据
type JSON json.RawMessage

// Value 实现 driver.Valuer 接口
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// Scan 实现 sql.Scanner 接口
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	*j = JSON(bytes)
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return nil
	}
	*j = JSON(data)
	return nil
}

// SystemConfigRepository 系统配置仓库接口
type SystemConfigRepository interface {
	Create(config *SystemConfig) error
	FindByID(id uint) (*SystemConfig, error)
	FindByType(configType string, includeInactive bool) ([]SystemConfig, error)
	Update(config *SystemConfig) error
	Delete(id uint) error
}

// SystemConfigService 系统配置服务接口
type SystemConfigService interface {
	GetConfigsByType(configType string, includeInactive bool) ([]SystemConfig, error)
	GetConfig(id uint) (*SystemConfig, error)
	CreateConfig(config *SystemConfig) error
	UpdateConfig(id uint, config *SystemConfig) error
	DeleteConfig(id uint) error
}

