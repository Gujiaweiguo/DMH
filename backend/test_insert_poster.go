package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PosterTemplate struct {
	Id            int64
	Type          string
	CampaignId    int64
	DistributorId int64
	TemplateUrl   string
	PosterData    string
	CreatedAt     string
}

func (PosterTemplate) TableName() string {
	return "poster_templates"
}

func main() {
	dsn := "root:#Admin168@tcp(127.0.0.1:3306)/dmh?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		return
	}
	defer db.Close()

	// 查询所有海报
	var posters []PosterTemplate
	result := db.Find(&posters)
	if result.Error != nil {
		fmt.Printf("查询海报失败: %v\n", result.Error)
		return
	}

	fmt.Printf("找到 %d 个海报记录:\n", len(posters))
	for _, p := range posters {
		fmt.Printf("ID: %d, Type: %s, URL: %s\n", p.Id, p.Type, p.TemplateUrl)
	}
}
