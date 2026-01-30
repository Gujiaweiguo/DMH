// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package poster

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPosterRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPosterRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPosterRecordsLogic {
	return &GetPosterRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPosterRecordsLogic) GetPosterRecords() (resp *types.PosterRecordsListResp, err error) {
	// 查询所有海报记录
	var records []model.PosterRecord
	if err := l.svcCtx.DB.Order("created_at DESC").Find(&records).Error; err != nil {
		l.Errorf("查询海报记录失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	recordList := make([]types.PosterRecordResp, 0, len(records))
	for _, r := range records {
		record := types.PosterRecordResp{
			Id:             r.ID,
			RecordType:     r.RecordType,
			CampaignId:     r.CampaignID,
			DistributorId:  r.DistributorID,
			TemplateName:   r.TemplateName,
			PosterUrl:      r.PosterUrl,
			ThumbnailUrl:   r.ThumbnailUrl,
			FileSize:       r.FileSize,
			GenerationTime: r.GenerationTime,
			DownloadCount:  r.DownloadCount,
			ShareCount:     r.ShareCount,
			Status:         r.Status,
			CreatedAt:      r.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if r.GeneratedBy != nil {
			record.GeneratedBy = *r.GeneratedBy
		}

		recordList = append(recordList, record)
	}

	resp = &types.PosterRecordsListResp{
		Total:   int64(len(records)),
		Records: recordList,
	}

	return resp, nil
}
