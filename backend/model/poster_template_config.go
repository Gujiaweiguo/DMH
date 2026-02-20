package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type PosterTemplateConfigData map[string]interface{}

func (d *PosterTemplateConfigData) Scan(value interface{}) error {
	if value == nil {
		*d = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unexpected type for PosterTemplateConfigData: %T", value)
	}

	return json.Unmarshal(bytes, d)
}

func (d PosterTemplateConfigData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

type PosterTemplateConfig struct {
	Id           int64                    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string                   `gorm:"column:name;type:varchar(100);not null" json:"name"`
	PreviewImage string                   `gorm:"column:preview_image;type:varchar(255)" json:"previewImage"`
	Config       PosterTemplateConfigData `gorm:"column:config;type:json;serializer:json" json:"config"`
	Status       string                   `gorm:"column:status;type:varchar(20);default:active" json:"status"`
	CreatedAt    time.Time                `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time                `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (PosterTemplateConfig) TableName() string {
	return "poster_template_configs"
}
