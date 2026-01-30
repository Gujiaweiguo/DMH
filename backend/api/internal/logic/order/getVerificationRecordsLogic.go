// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVerificationRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVerificationRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVerificationRecordsLogic {
	return &GetVerificationRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVerificationRecordsLogic) GetVerificationRecords() (resp *types.VerificationRecordsListResp, err error) {
	// 查询所有核销记录
	var records []model.VerificationRecord
	if err := l.svcCtx.DB.Order("created_at DESC").Find(&records).Error; err != nil {
		l.Errorf("查询核销记录失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	recordList := make([]types.VerificationRecordResp, 0, len(records))
	for _, r := range records {
		verifiedAt := ""
		if r.VerifiedAt != nil {
			verifiedAt = r.VerifiedAt.Format("2006-01-02 15:04:05")
		}

		record := types.VerificationRecordResp{
			Id:                 r.ID,
			OrderId:            r.OrderID,
			VerificationStatus: r.VerificationStatus,
			VerifiedAt:         verifiedAt,
			VerificationCode:   r.VerificationCode,
			VerificationMethod: r.VerificationMethod,
			Remark:             r.Remark,
			CreatedAt:          r.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if r.VerifiedBy != nil {
			record.VerifiedBy = *r.VerifiedBy
		}

		recordList = append(recordList, record)
	}

	resp = &types.VerificationRecordsListResp{
		Total:   int64(len(records)),
		Records: recordList,
	}

	return resp, nil
}
