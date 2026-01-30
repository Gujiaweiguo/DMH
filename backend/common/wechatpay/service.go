package wechatpay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Config 微信支付配置
type Config struct {
	AppID           string
	MchID           string
	APIKey          string
	APIKeyV3        string
	APIClientCert   string
	APIClientKey    string
	NotifyURL       string
	RefundNotifyURL string
	Sandbox         bool
	CacheTTL        int
}

// NativePayRequest Native支付请求
type NativePayRequest struct {
	AppID          string `json:"appid"`
	MchID          string `json:"mch_id"`
	NonceStr       string `json:"nonce_str"`
	Sign           string `json:"sign"`
	Body           string `json:"body"`
	OutTradeNo     string `json:"out_trade_no"`
	TotalFee       int64  `json:"total_fee"`
	SpbillCreateIP string `json:"spbill_create_ip"`
	TradeType      string `json:"trade_type"`
	Attach         string `json:"attach,omitempty"`
	TimeExpire     string `json:"time_expire,omitempty"`
}

// NativePayResponse Native支付响应
type NativePayResponse struct {
	ReturnCode string `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	AppID      string `json:"appid"`
	MchID      string `json:"mch_id"`
	NonceStr   string `json:"nonce_str"`
	Sign       string `json:"sign"`
	PrepayID   string `json:"prepay_id"`
	CodeURL    string `json:"code_url"`
}

// NotifyRequest 支付通知请求
type NotifyRequest struct {
	ReturnCode    string `xml:"return_code"`
	AppID         string `xml:"appid"`
	MchID         string `xml:"mch_id"`
	DeviceInfo    string `xml:"device_info"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	OpenID        string `xml:"openid"`
	OutTradeNo    string `xml:"out_trade_no"`
	TransactionID string `xml:"transaction_id"`
	TotalFee      int64  `xml:"total_fee"`
	FeeType       string `xml:"fee_type"`
	TimeEnd       string `xml:"time_end"`
	Attach        string `xml:"attach"`
	IsSubscribe   string `xml:"is_subscribe"`
}

// Service 微信支付服务
type Service struct {
	config *Config
}

// NewService 创建微信支付服务
func NewService(config *Config) *Service {
	return &Service{
		config: config,
	}
}

// GenerateMD5Sign 生成MD5签名
func (s *Service) GenerateMD5Sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(params[k])
	}

	data := buf.String() + s.config.APIKey
	hash := md5.Sum([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// VerifyMD5Sign 验证MD5签名
func (s *Service) VerifyMD5Sign(params map[string]string, sign string) bool {
	computedSign := s.GenerateMD5Sign(params)
	return computedSign == sign
}

// CreateNativePay 创建Native支付订单
func (s *Service) CreateNativePay(orderNo string, amount int64, body string) (*NativePayResponse, error) {
	nonce := generateNonce()
	attach := fmt.Sprintf("order=%s&time=%d", orderNo, time.Now().Unix())

	// 构建请求参数
	params := map[string]string{
		"appid":            s.config.AppID,
		"mch_id":           s.config.MchID,
		"nonce_str":        nonce,
		"body":             body,
		"out_trade_no":     orderNo,
		"total_fee":        fmt.Sprintf("%d", amount),
		"spbill_create_ip": "127.0.0.1",
		"trade_type":       "NATIVE",
		"attach":           attach,
		"time_expire":      time.Now().Add(2 * time.Hour).Format("20060102150405"),
	}

	// 生成签名
	sign := s.GenerateMD5Sign(params)
	params["sign"] = sign

	// TODO: 实际调用微信支付API
	// 这里返回模拟响应
	return &NativePayResponse{
		ReturnCode: "SUCCESS",
		AppID:      s.config.AppID,
		MchID:      s.config.MchID,
		NonceStr:   nonce,
		Sign:       sign,
		PrepayID:   fmt.Sprintf("prepay_%s", nonce),
		CodeURL:    fmt.Sprintf("weixin://wxpay/bizpayurl?pr=%s", orderNo),
	}, nil
}

// ParseNotifyRequest 解析支付通知请求
func (s *Service) ParseNotifyRequest(xmlData string) (*NotifyRequest, error) {
	var notify NotifyRequest
	if err := xml.Unmarshal([]byte(xmlData), &notify); err != nil {
		return nil, fmt.Errorf("解析支付通知失败: %w", err)
	}

	if notify.ReturnCode != "SUCCESS" {
		return nil, errors.New("支付失败")
	}

	// 验证签名
	params := map[string]string{
		"return_code":    notify.ReturnCode,
		"appid":          notify.AppID,
		"mch_id":         notify.MchID,
		"device_info":    notify.DeviceInfo,
		"nonce_str":      notify.NonceStr,
		"out_trade_no":   notify.OutTradeNo,
		"transaction_id": notify.TransactionID,
		"total_fee":      fmt.Sprintf("%d", notify.TotalFee),
		"fee_type":       notify.FeeType,
		"time_end":       notify.TimeEnd,
		"attach":         notify.Attach,
		"is_subscribe":   notify.IsSubscribe,
	}

	if !s.VerifyMD5Sign(params, notify.Sign) {
		return nil, errors.New("签名验证失败")
	}

	return &notify, nil
}

// BuildNotifyResponse 构建支付通知响应
func BuildNotifyResponse() string {
	params := map[string]string{
		"return_code": "SUCCESS",
		"return_msg":  "OK",
	}

	xmlData, _ := xml.Marshal(params)
	return string(xmlData)
}

// generateNonce 生成随机字符串
func generateNonce() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().UnixNano())
}
