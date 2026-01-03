// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersLogic {
	return &GetUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersLogic) GetUsers(req *types.AdminGetUsersReq) (resp *types.AdminUserListResp, err error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建查询条件
	query := l.svcCtx.DB.Model(&model.User{})

	// 角色过滤
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR phone LIKE ? OR real_name LIKE ?", 
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		l.Logger.Errorf("获取用户总数失败: %v", err)
		return nil, fmt.Errorf("获取用户列表失败")
	}

	// 分页查询
	var users []model.User
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&users).Error; err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, fmt.Errorf("获取用户列表失败")
	}

	// 获取用户的品牌关联信息
	userIds := make([]int64, len(users))
	for i, user := range users {
		userIds[i] = user.Id
	}

	var userBrands []model.UserBrand
	if len(userIds) > 0 {
		l.svcCtx.DB.Where("user_id IN ?", userIds).Find(&userBrands)
	}

	// 构建用户品牌映射
	userBrandMap := make(map[int64][]int64)
	for _, ub := range userBrands {
		userBrandMap[ub.UserId] = append(userBrandMap[ub.UserId], ub.BrandId)
	}

	// 转换为响应格式
	userList := make([]types.UserInfoResp, len(users))
	for i, user := range users {
		userList[i] = types.UserInfoResp{
			Id:        user.Id,
			Username:  user.Username,
			Phone:     user.Phone,
			Email:     user.Email,
			RealName:  user.RealName,
			Status:    user.Status,
			Roles:     []string{user.Role},
			BrandIds:  userBrandMap[user.Id],
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &types.AdminUserListResp{
		Total: total,
		Users: userList,
	}, nil
}
