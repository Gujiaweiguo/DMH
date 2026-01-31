// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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
	// 1. 从核销码中提取订单ID、手机号、时间戳、签名
	orderId, phone, timestamp, signature := parseVerificationCode(req.Code)
	if orderId == 0 {
		l.Errorf("核销码格式无效: %s", req.Code)
		return nil, fmt.Errorf("核销码无效")
	}

	// 2. 验证签名
	expectedSignature := generateSignature(fmt.Sprintf("%d", orderId), phone, timestamp)
	if signature != expectedSignature {
		l.Errorf("核销码签名验证失败: orderId=%d, code=%s", orderId, req.Code)
		return nil, fmt.Errorf("核销码无效")
	}

	// 3. 查询订单
	var order model.Order
	if err := l.svcCtx.DB.Where("id = ? AND pay_status = ?", orderId, "paid").First(&order).Error; err != nil {
		l.Errorf("查询订单失败: %v", err)
		return nil, fmt.Errorf("订单不存在或未支付")
	}

	// 4. 验证订单状态（不允许重复核销）
	if order.VerificationStatus == "verified" {
		return nil, fmt.Errorf("订单已核销")
	}

	// 5. 验证订单手机号是否匹配
	if order.Phone != phone {
		l.Errorf("订单手机号不匹配: orderPhone=%s, codePhone=%s", order.Phone, phone)
		return nil, fmt.Errorf("核销码无效")
	}

	// 5. 开始数据库事务
	tx := l.svcCtx.DB.Begin()

	now := time.Now()

	// 6. 更新订单核销状态
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"verification_status": "verified",
		"verified_at":         &now,
	}).Error; err != nil {
		tx.Rollback()
		l.Errorf("更新订单核销状态失败: %v", err)
		return nil, fmt.Errorf("核销失败: %w", err)
	}

	// 7. 创建核销记录
	verificationRecord := model.VerificationRecord{
		OrderID:            order.Id,
		VerificationStatus: "verified",
		VerifiedAt:         &now,
		VerificationCode:   req.Code,
		VerificationMethod: "manual",
		Remark:             "",
	}
	if err := tx.Create(&verificationRecord).Error; err != nil {
		tx.Rollback()
		l.Errorf("创建核销记录失败: %v", err)
		return nil, fmt.Errorf("核销失败: %w", err)
	}

	// 8. 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("核销失败: %w", err)
	}

	l.Infof("订单核销成功: orderId=%d, code=%s", orderId, req.Code)

	// 9. 返回响应
	resp = &types.VerifyOrderResp{
		OrderId:    order.Id,
		Status:     "verified",
		VerifiedAt: now.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

// parseVerificationCode 解析核销码
// 格式：{orderId}_{phone}_{timestamp}_{signature}
func parseVerificationCode(code string) (orderId int64, phone string, timestamp string, signature string) {
	if code == "" {
		return 0, "", "", ""
	}

	parts := make([]string, 4)
	idx := 0
	start := 0
	for i, c := range code {
		if c == '_' {
			parts[idx] = code[start:i]
			idx++
			start = i + 1
		}
	}
	parts[idx] = code[start:]

	if len(parts) != 4 {
		return 0, "", "", ""
	}

	orderIdVal := int64(0)
	if parts[0] != "" {
		fmt.Sscanf(parts[0], "%d", &orderIdVal)
	}

	return orderIdVal, parts[1], parts[2], parts[3]
}

// verifyCode 验证核销码签名
func verifyCode(code string, orderId int64, phone, timestamp string) bool {
	// 签名：HMAC-SHA1({orderId}_{phone}_{timestamp}, "dmh_secret_key")
	expectedSignature := generateSignature(fmt.Sprintf("%d", orderId), phone, timestamp)

	_, _, _, signature := parseVerificationCode(code)
	if signature == "" {
		return false
	}

	if signature != expectedSignature {
		return false
	}

	return true
}

// generateSignature 生成签名
func generateSignature(orderId, phone, timestamp string) string {
	secret := "dmh_secret_key"
	data := fmt.Sprintf("%s_%s_%s", fmt.Sprintf("%d", orderId), phone, timestamp)
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
