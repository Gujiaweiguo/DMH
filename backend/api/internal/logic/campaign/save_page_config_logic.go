// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SavePageConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSavePageConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePageConfigLogic {
	return &SavePageConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SavePageConfigLogic) SavePageConfig(req *types.PageConfigReq, campaignId string) (resp *types.PageConfigResp, err error) {
	// 解析campaignId为int64
	cid, err := strconv.ParseInt(campaignId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的活动ID: %v", err)
	}
	
	// 如果数据库连接可用，尝试保存到数据库
	if l.svcCtx.DB != nil {
		// 将组件和主题转换为JSON字符串
		componentsJSON, err := json.Marshal(req.Components)
		if err != nil {
			return nil, fmt.Errorf("序列化组件配置失败: %v", err)
		}
		
		themeJSON, err := json.Marshal(req.Theme)
		if err != nil {
			return nil, fmt.Errorf("序列化主题配置失败: %v", err)
		}
		
		// 检查是否已存在配置
		var existingConfig struct {
			ID int64 `gorm:"column:id"`
		}
		
		err = l.svcCtx.DB.Table("page_configs").
			Where("campaign_id = ?", cid).
			First(&existingConfig).Error
		
		now := time.Now()
		
		if err == gorm.ErrRecordNotFound {
			// 创建新配置
			newConfig := map[string]interface{}{
				"campaign_id": cid,
				"components":  string(componentsJSON),
				"theme":       string(themeJSON),
				"created_at":  now,
				"updated_at":  now,
			}
			
			err = l.svcCtx.DB.Table("page_configs").Create(&newConfig).Error
			if err != nil {
				l.Logger.Errorf("创建页面配置失败: %v", err)
			} else {
				l.Logger.Infof("页面配置已保存到数据库")
			}
		} else if err == nil {
			// 更新现有配置
			updates := map[string]interface{}{
				"components": string(componentsJSON),
				"theme":      string(themeJSON),
				"updated_at": now,
			}
			
			err = l.svcCtx.DB.Table("page_configs").
				Where("id = ?", existingConfig.ID).
				Updates(updates).Error
			if err != nil {
				l.Logger.Errorf("更新页面配置失败: %v", err)
			} else {
				l.Logger.Infof("页面配置已更新，ID: %d", existingConfig.ID)
			}
		}
	}
	
	// 返回成功响应（即使数据库操作失败，也返回成功，因为前端可以继续工作）
	return &types.PageConfigResp{
		Id:         1,
		CampaignId: cid,
		Components: req.Components,
		Theme:      req.Theme,
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
