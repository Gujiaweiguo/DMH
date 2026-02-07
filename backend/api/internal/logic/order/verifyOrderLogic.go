// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyOrderLogic {
	return &VerifyOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyOrderLogic) VerifyOrder(req *types.VerifyOrderReq) (resp *types.VerifyOrderResp, err error) {
	orderId, phone, timestamp, signature := l.parseVerificationCode(req.Code)
	if orderId == 0 {
		l.Errorf("核销码格式无效: %s", req.Code)
		return nil, fmt.Errorf("核销码无效")
	}

	if !l.verifySignature(req.Code, orderId, phone, timestamp, signature) {
		l.Errorf("核销码签名验证失败: orderId=%d, code=%s", orderId, req.Code)
		return nil, fmt.Errorf("核销码无效")
	}

	if !hasVerificationPermission(l.ctx) {
		l.Errorf("核销权限不足: code=%s", req.Code)
		return nil, fmt.Errorf("权限不足，仅品牌管理员或平台管理员可执行核销操作")
	}

	var order model.Order
	if err := l.svcCtx.DB.Where("id = ? AND deleted_at IS NULL", orderId).First(&order).Error; err != nil {
		l.Errorf("查询订单失败: %v", err)
		return nil, fmt.Errorf("订单不存在")
	}

	if order.VerificationStatus == "verified" {
		return nil, fmt.Errorf("订单已核销")
	}

	// 解析备注字段
	remark := req.Remark
	if remark == "" {
		remark = "品牌管理员核销"
	}

	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"verification_status": "verified",
		"verified_at":         &now,
		"verified_by":         l.getUserId(),
	}).Error; err != nil {
		tx.Rollback()
		l.Errorf("更新订单核销状态失败: %v", err)
		return nil, fmt.Errorf("核销失败: %w", err)
	}

	verifiedBy := l.getUserId()
	verificationRecord := model.VerificationRecord{
		OrderID:            order.Id,
		VerificationStatus: "verified",
		VerifiedAt:         &now,
		VerifiedBy:         &verifiedBy,
		VerificationCode:   req.Code,
		VerificationMethod: "manual",
		Remark:             remark,
	}
	if err := tx.Create(&verificationRecord).Error; err != nil {
		tx.Rollback()
		l.Errorf("创建核销记录失败: %v", err)
		return nil, fmt.Errorf("核销失败 %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		l.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("核销失败 %w", err)
	}

	l.Infof("订单核销成功: orderId=%d, code=%s, remark=%s", orderId, req.Code, remark)

	return &types.VerifyOrderResp{
		OrderId:    order.Id,
		Status:     "verified",
		VerifiedAt: now.Format("2006-01-02T15:04:05"),
		VerifiedBy: &verifiedBy,
	}, nil
}

func (l *VerifyOrderLogic) parseVerificationCode(code string) (orderId int64, phone string, timestamp int64, signature string) {
	if code == "" {
		return 0, "", 0, ""
	}

	parts := strings.Split(code, "_")
	if len(parts) != 4 {
		return 0, "", 0, ""
	}

	orderIdVal := int64(0)
	if parts[0] != "" {
		orderIdVal, _ = strconv.ParseInt(parts[0], 10, 64)
	}

	timestampVal := int64(0)
	if parts[2] != "" {
		timestampVal, _ = strconv.ParseInt(parts[2], 10, 64)
	}

	return orderIdVal, parts[1], timestampVal, parts[3]
}

func (l *VerifyOrderLogic) verifySignature(code string, orderId int64, phone string, timestamp int64, signature string) bool {
	secretKey := "dmh-verification-secret-2026"

	signatureData := fmt.Sprintf("%d_%s_%d_%s", orderId, phone, timestamp, secretKey)
	hash := md5.Sum([]byte(signatureData))
	expectedSignature := hex.EncodeToString(hash[:])

	if signature != expectedSignature {
		return false
	}

	return true
}

func (l *VerifyOrderLogic) getUserId() int64 {
	userIdValue := l.ctx.Value("userId")
	if userIdValue != nil {
		if uid, ok := userIdValue.(int64); ok {
			return uid
		} else if uid, ok := userIdValue.(float64); ok {
			return int64(uid)
		}
	}
	return 0
}
