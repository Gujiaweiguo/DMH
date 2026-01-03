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

type GetPageConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPageConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPageConfigLogic {
	return &GetPageConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPageConfigLogic) GetPageConfig(campaignId string) (resp *types.PageConfigResp, err error) {
	// 解析campaignId为int64
	cid, err := strconv.ParseInt(campaignId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的活动ID: %v", err)
	}
	
	// 如果数据库连接可用，尝试从数据库加载
	if l.svcCtx.DB != nil {
		var pageConfig struct {
			ID         int64     `gorm:"column:id"`
			Components string    `gorm:"column:components"`
			Theme      string    `gorm:"column:theme"`
			CreatedAt  time.Time `gorm:"column:created_at"`
			UpdatedAt  time.Time `gorm:"column:updated_at"`
		}
		
		err = l.svcCtx.DB.Table("page_configs").
			Where("campaign_id = ?", cid).
			First(&pageConfig).Error
		
		if err == nil {
			// 解析JSON
			var components []map[string]interface{}
			if err := json.Unmarshal([]byte(pageConfig.Components), &components); err != nil {
				l.Logger.Errorf("解析组件配置失败: %v", err)
			} else {
				var theme map[string]interface{}
				if err := json.Unmarshal([]byte(pageConfig.Theme), &theme); err != nil {
					l.Logger.Errorf("解析主题配置失败: %v", err)
				} else {
					// 成功从数据库加载
					return &types.PageConfigResp{
						Id:         pageConfig.ID,
						CampaignId: cid,
						Components: components,
						Theme:      theme,
						CreatedAt:  pageConfig.CreatedAt.Format("2006-01-02 15:04:05"),
						UpdatedAt:  pageConfig.UpdatedAt.Format("2006-01-02 15:04:05"),
					}, nil
				}
			}
		} else if err != gorm.ErrRecordNotFound {
			l.Logger.Errorf("查询页面配置失败: %v", err)
		}
	}
	
	// 数据库中没有配置，返回404
	return nil, fmt.Errorf("页面配置不存在")
}
