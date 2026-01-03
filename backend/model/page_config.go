package model

import (
	"database/sql"
	"time"
)

// PageConfig 页面配置模型
type PageConfig struct {
	Id         int64        `db:"id"`
	CampaignId int64        `db:"campaign_id"`
	Components string       `db:"components"` // JSON格式存储组件配置
	Theme      string       `db:"theme"`      // JSON格式存储主题配置
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at"`
	DeletedAt  sql.NullTime `db:"deleted_at"`
}

// TableName 表名
func (m *PageConfig) TableName() string {
	return "page_configs"
}
