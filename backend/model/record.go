package model

import "time"

// VerificationRecord 核销记录表
type VerificationRecord struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID            int64      `gorm:"column:order_id;not null;index" json:"orderId"`
	VerificationStatus string     `gorm:"column:verification_status;type:varchar(20);not null;default:pending;index" json:"verificationStatus"` // pending/verified/cancelled
	VerifiedAt         *time.Time `gorm:"column:verified_at" json:"verifiedAt"`
	VerifiedBy         *int64     `gorm:"column:verified_by;index" json:"verifiedBy"`
	VerificationCode   string     `gorm:"column:verification_code;type:varchar(128)" json:"verificationCode"`
	VerificationMethod string     `gorm:"column:verification_method;type:varchar(20);not null;default:manual" json:"verificationMethod"` // manual/auto/qrcode
	Remark             string     `gorm:"column:remark;type:varchar(500)" json:"remark"`
	CreatedAt          time.Time  `gorm:"column:created_at;not null;autoCreateTime;index" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (VerificationRecord) TableName() string {
	return "verification_records"
}

// PosterRecord 海报生成记录表
type PosterRecord struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RecordType     string    `gorm:"column:record_type;type:varchar(20);not null;default:personal;index" json:"recordType"` // personal/brand
	CampaignID     int64     `gorm:"column:campaign_id;not null;index" json:"campaignId"`
	DistributorID  int64     `gorm:"column:distributor_id;default:0;index" json:"distributorId"`
	TemplateName   string    `gorm:"column:template_name;type:varchar(100);not null" json:"templateName"`
	PosterUrl      string    `gorm:"column:poster_url;type:varchar(500);not null" json:"posterUrl"`
	ThumbnailUrl   string    `gorm:"column:thumbnail_url;type:varchar(500)" json:"thumbnailUrl"`
	FileSize       string    `gorm:"column:file_size;type:varchar(50)" json:"fileSize"`
	GenerationTime int       `gorm:"column:generation_time;default:0" json:"generationTime"` // 毫秒
	DownloadCount  int       `gorm:"column:download_count;default:0" json:"downloadCount"`
	ShareCount     int       `gorm:"column:share_count;default:0" json:"shareCount"`
	GeneratedBy    *int64    `gorm:"column:generated_by" json:"generatedBy"`
	Status         string    `gorm:"column:status;type:varchar(20);not null;default:active;index" json:"status"` // active/deleted
	CreatedAt      time.Time `gorm:"column:created_at;not null;autoCreateTime;index" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (PosterRecord) TableName() string {
	return "poster_records"
}
