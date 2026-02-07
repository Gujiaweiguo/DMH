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

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnverifyOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnverifyOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnverifyOrderLogic {
	return &UnverifyOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnverifyOrderLogic) UnverifyOrder(req *types.UnverifyOrderReq) (resp *types.UnverifyOrderResp, err error) {
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
		l.Errorf("取消核销权限不足: code=%s", req.Code)
		return nil, fmt.Errorf("权限不足，仅品牌管理员或平台管理员可执行取消核销操作")
	}

	var order model.Order
	if err := l.svcCtx.DB.Where("id = ? AND deleted_at IS NULL", orderId).First(&order).Error; err != nil {
		l.Errorf("查询订单失败: %v", err)
		return nil, fmt.Errorf("订单不存在")
	}

	reason := "品牌管理员取消核销"

	if order.VerificationStatus != "verified" {
		return nil, fmt.Errorf("订单尚未核销，无法取消核销")
	}

	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"verification_status": "cancelled",
		"verified_at":         nil,
		"verified_by":         nil,
	}).Error; err != nil {
		tx.Rollback()
		l.Errorf("取消订单核销状态失败: %v", err)
		return nil, fmt.Errorf("取消核销失败: %w", err)
	}

	verificationRecord := model.VerificationRecord{
		OrderID:            order.Id,
		VerificationStatus: "cancelled",
		VerifiedAt:         nil,
		VerifiedBy:         nil,
		VerificationCode:   req.Code,
		VerificationMethod: "manual",
		Remark:             reason,
	}
	if err := tx.Create(&verificationRecord).Error; err != nil {
		tx.Rollback()
		l.Errorf("创建取消核销记录失败: %v", err)
		return nil, fmt.Errorf("取消核销失败: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		l.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("取消核销失败: %w", err)
	}

	l.Infof("订单取消核销成功: orderId=%d, code=%s, reason=%s", orderId, req.Code, reason)

	return &types.UnverifyOrderResp{
		OrderId: order.Id,
		Status:  "unverified",
	}, nil
}

func (l *UnverifyOrderLogic) parseVerificationCode(code string) (orderId int64, phone string, timestamp int64, signature string) {
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

func (l *UnverifyOrderLogic) verifySignature(code string, orderId int64, phone string, timestamp int64, signature string) bool {
	secretKey := "dmh-verification-secret-2026"

	signatureData := fmt.Sprintf("%d_%s_%d_%s", orderId, phone, timestamp, secretKey)
	hash := md5.Sum([]byte(signatureData))
	expectedSignature := hex.EncodeToString(hash[:])

	return signature == expectedSignature
}
