package member

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveExportRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveExportRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveExportRequestLogic {
	return &ApproveExportRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveExportRequestLogic) ApproveExportRequest(requestId int64, req *types.ApproveExportReq) error {
	// 只有平台管理员可以审批
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return err
	}
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user.Role != "platform_admin" {
		return fmt.Errorf("只有平台管理员可以审批导出请求")
	}

	// 查询导出请求
	var exportReq model.ExportRequest
	if err := l.svcCtx.DB.Where("id = ?", requestId).First(&exportReq).Error; err != nil {
		return fmt.Errorf("导出请求不存在: %w", err)
	}

	if exportReq.Status != "pending" {
		return fmt.Errorf("该导出请求已被处理")
	}

	now := time.Now()

	if !req.Approve {
		// 驳回
		exportReq.Status = "rejected"
		exportReq.ApprovedBy = &userId
		exportReq.ApprovedAt = &now
		exportReq.RejectReason = req.Reason

		if err := l.svcCtx.DB.Save(&exportReq).Error; err != nil {
			return fmt.Errorf("更新导出请求失败: %w", err)
		}

		return nil
	}

	// 批准并生成导出
	exportReq.Status = "approved"
	exportReq.ApprovedBy = &userId
	exportReq.ApprovedAt = &now

	if err := l.svcCtx.DB.Save(&exportReq).Error; err != nil {
		return fmt.Errorf("更新导出请求失败: %w", err)
	}

	// 生成导出文件
	if err := l.generateExportFile(&exportReq); err != nil {
		// 更新为失败状态
		exportReq.Status = "failed"
		l.svcCtx.DB.Save(&exportReq)
		return fmt.Errorf("生成导出文件失败: %w", err)
	}

	// 更新为完成状态
	exportReq.Status = "completed"
	if err := l.svcCtx.DB.Save(&exportReq).Error; err != nil {
		return fmt.Errorf("更新导出请求状态失败: %w", err)
	}

	return nil
}

func (l *ApproveExportRequestLogic) generateExportFile(exportReq *model.ExportRequest) error {
	// 查询该品牌的会员数据
	query := l.svcCtx.DB.Model(&model.Member{}).
		Joins("INNER JOIN member_brand_links ON members.id = member_brand_links.member_id").
		Where("member_brand_links.brand_id = ?", exportReq.BrandID)

	// TODO: 根据 exportReq.Filters 添加筛选条件

	var members []model.Member
	if err := query.Find(&members).Error; err != nil {
		return fmt.Errorf("查询会员数据失败: %w", err)
	}

	// 查询会员画像
	memberIds := make([]int64, len(members))
	for i, m := range members {
		memberIds[i] = m.ID
	}

	var profiles []model.MemberProfile
	if len(memberIds) > 0 {
		l.svcCtx.DB.Where("member_id IN ?", memberIds).Find(&profiles)
	}

	profileMap := make(map[int64]model.MemberProfile)
	for _, p := range profiles {
		profileMap[p.MemberID] = p
	}

	// 创建导出目录
	exportDir := "./storage/exports"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return fmt.Errorf("创建导出目录失败: %w", err)
	}

	// 生成文件名
	filename := fmt.Sprintf("members_brand_%d_%s.csv", exportReq.BrandID, time.Now().Format("20060102150405"))
	filepath := filepath.Join(exportDir, filename)

	// 创建 CSV 文件
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	headers := []string{
		"会员ID", "UnionID", "昵称", "手机号", "性别", "来源", "状态",
		"累计订单数", "累计支付金额", "累计奖励金额", "参与活动数",
		"首次下单时间", "最后下单时间", "创建时间",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("写入表头失败: %w", err)
	}

	// 写入数据
	for _, m := range members {
		profile := profileMap[m.ID]

		genderStr := "未知"
		if m.Gender == 1 {
			genderStr = "男"
		} else if m.Gender == 2 {
			genderStr = "女"
		}

		firstOrderAt := ""
		if profile.FirstOrderAt != nil {
			firstOrderAt = profile.FirstOrderAt.Format("2006-01-02 15:04:05")
		}

		lastOrderAt := ""
		if profile.LastOrderAt != nil {
			lastOrderAt = profile.LastOrderAt.Format("2006-01-02 15:04:05")
		}

		row := []string{
			fmt.Sprintf("%d", m.ID),
			m.UnionID,
			m.Nickname,
			m.Phone,
			genderStr,
			m.Source,
			m.Status,
			fmt.Sprintf("%d", profile.TotalOrders),
			fmt.Sprintf("%.2f", profile.TotalPayment),
			fmt.Sprintf("%.2f", profile.TotalReward),
			fmt.Sprintf("%d", profile.ParticipatedCampaigns),
			firstOrderAt,
			lastOrderAt,
			m.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("写入数据失败: %w", err)
		}
	}

	// 更新导出请求
	exportReq.FileUrl = "/exports/" + filename
	exportReq.RecordCount = len(members)

	return nil
}
