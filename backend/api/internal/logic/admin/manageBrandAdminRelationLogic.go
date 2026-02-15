// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageBrandAdminRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageBrandAdminRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageBrandAdminRelationLogic {
	return &ManageBrandAdminRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageBrandAdminRelationLogic) ManageBrandAdminRelation(req *types.BrandAdminRelationReq) (resp *types.CommonResp, err error) {
	if req.UserId <= 0 {
		return nil, errors.New("用户ID无效")
	}

	if len(req.BrandIds) == 0 {
		return nil, errors.New("品牌ID列表不能为空")
	}

	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("user_id = ?", req.UserId).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("删除旧关系失败")
	}

	for _, brandId := range req.BrandIds {
		userRole := &model.UserRole{
			UserID: req.UserId,
			RoleID: brandId,
		}
		if err := tx.Create(userRole).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("创建关系失败")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("提交事务失败")
	}

	resp = &types.CommonResp{
		Message: "操作成功",
	}

	return resp, nil
}
