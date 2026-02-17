// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginAttemptsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoginAttemptsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginAttemptsLogic {
	return &GetLoginAttemptsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoginAttemptsLogic) GetLoginAttempts() (resp *types.LoginAttemptListResp, err error) {
	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return &types.LoginAttemptListResp{Total: 0, Attempts: []types.LoginAttemptResp{}}, nil
	}

	var total int64
	if err := l.svcCtx.DB.Model(&model.LoginAttempt{}).Count(&total).Error; err != nil {
		l.Errorf("查询登录尝试总数失败: %v", err)
		return nil, err
	}

	if total == 0 {
		return &types.LoginAttemptListResp{Total: 0, Attempts: []types.LoginAttemptResp{}}, nil
	}

	var attempts []model.LoginAttempt
	if err := l.svcCtx.DB.Order("created_at DESC, id DESC").Find(&attempts).Error; err != nil {
		l.Errorf("查询登录尝试失败: %v", err)
		return nil, err
	}

	respAttempts := make([]types.LoginAttemptResp, 0, len(attempts))
	for _, attempt := range attempts {
		var userID int64
		if attempt.UserID != nil {
			userID = *attempt.UserID
		}

		respAttempts = append(respAttempts, types.LoginAttemptResp{
			Id:         attempt.ID,
			UserId:     userID,
			Username:   attempt.Username,
			ClientIp:   attempt.ClientIP,
			UserAgent:  attempt.UserAgent,
			Success:    attempt.Success,
			FailReason: attempt.FailReason,
			CreatedAt:  attempt.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.LoginAttemptListResp{Total: total, Attempts: respAttempts}, nil
}
