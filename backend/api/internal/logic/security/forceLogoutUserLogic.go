// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForceLogoutUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForceLogoutUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForceLogoutUserLogic {
	return &ForceLogoutUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForceLogoutUserLogic) ForceLogoutUser(userID int64, req *types.ForceLogoutReq) (resp *types.CommonResp, err error) {
	if userID <= 0 {
		return nil, errors.New("用户ID无效")
	}

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	reason := "管理员强制下线"
	if req != nil {
		if trimmed := strings.TrimSpace(req.Reason); trimmed != "" {
			reason = trimmed
		}
	}

	now := time.Now()
	result := l.svcCtx.DB.Model(&model.UserSession{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Updates(map[string]interface{}{
			"status":     "revoked",
			"updated_at": now,
		})
	if result.Error != nil {
		l.Errorf("强制下线失败: %v", result.Error)
		return nil, errors.New("强制下线失败")
	}

	operatorID, hasOperator := userIDFromContext(l.ctx)
	username, _ := l.ctx.Value("username").(string)
	var auditUserID *int64
	if hasOperator {
		auditUserID = &operatorID
	}

	auditDetails := reason
	if result.RowsAffected > 0 {
		auditDetails = reason + "，撤销活跃会话数: " + strconv.FormatInt(result.RowsAffected, 10)
	}

	if err := l.svcCtx.DB.Create(&model.AuditLog{
		UserID:     auditUserID,
		Username:   username,
		Action:     "force_logout_user",
		Resource:   "user_session",
		ResourceID: strconv.FormatInt(userID, 10),
		Details:    auditDetails,
		Status:     "success",
	}).Error; err != nil {
		l.Errorf("记录强制下线审计日志失败: %v", err)
	}

	if result.RowsAffected == 0 {
		return &types.CommonResp{Message: "用户当前无活跃会话"}, nil
	}

	return &types.CommonResp{Message: "用户已强制下线"}, nil
}
