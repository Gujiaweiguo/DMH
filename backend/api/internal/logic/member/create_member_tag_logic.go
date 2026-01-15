package member

import (
	"context"
	"fmt"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMemberTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMemberTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMemberTagLogic {
	return &CreateMemberTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMemberTagLogic) CreateMemberTag(req *types.CreateMemberTagReq) (*types.MemberTagResp, error) {
	// 只有平台管理员可以创建标签
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user.Role != "platform_admin" {
		return nil, fmt.Errorf("只有平台管理员可以创建标签")
	}

	// 检查标签名称是否已存在
	var count int64
	l.svcCtx.DB.Model(&model.MemberTag{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("标签名称已存在")
	}

	// 创建标签
	tag := model.MemberTag{
		Name:        req.Name,
		Category:    req.Category,
		Color:       req.Color,
		Description: req.Description,
	}

	if err := l.svcCtx.DB.Create(&tag).Error; err != nil {
		return nil, fmt.Errorf("创建标签失败: %w", err)
	}

	return &types.MemberTagResp{
		Id:          tag.ID,
		Name:        tag.Name,
		Category:    tag.Category,
		Color:       tag.Color,
		Description: tag.Description,
	}, nil
}
