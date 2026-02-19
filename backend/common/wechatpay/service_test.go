package wechatpay

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestService() *Service {
	return NewService(&Config{
		AppID:       "wx-app",
		MchID:       "mch-1",
		APIKey:      "key123",
		MockEnabled: true,
	})
}

func TestGenerateAndVerifySign(t *testing.T) {
	s := newTestService()
	params := map[string]string{
		"appid":        "wx-app",
		"mch_id":       "mch-1",
		"nonce_str":    "abc",
		"out_trade_no": "o-1",
		"total_fee":    "100",
	}
	sign := s.GenerateMD5Sign(params)
	if sign == "" {
		t.Fatalf("empty sign")
	}
	if !s.VerifyMD5Sign(params, sign) {
		t.Fatalf("verify sign failed")
	}
	if s.VerifyMD5Sign(params, "BAD") {
		t.Fatalf("verify should fail")
	}
}

func TestCreateNativePay(t *testing.T) {
	s := newTestService()
	resp, err := s.CreateNativePay("order-1", 520, "body")
	if err != nil {
		t.Fatalf("create native pay error: %v", err)
	}
	if resp == nil || resp.ReturnCode != "SUCCESS" {
		t.Fatalf("unexpected response")
	}
	if !strings.Contains(resp.CodeURL, "weixin://wxpay") {
		t.Fatalf("unexpected code_url: %s", resp.CodeURL)
	}
}

func TestCreateNativePayRealModeSuccess(t *testing.T) {
	cfg := &Config{
		AppID:       "wx-app",
		MchID:       "mch-1",
		APIKey:      "key123",
		MockEnabled: false,
	}
	s := NewService(cfg)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		nonce := "resp_nonce"
		prepayID := "prepay123"
		codeURL := "weixin://wxpay/bizpayurl?pr=order-1"
		sign := s.GenerateMD5Sign(map[string]string{
			"return_code":  "SUCCESS",
			"return_msg":   "OK",
			"result_code":  "SUCCESS",
			"appid":        cfg.AppID,
			"mch_id":       cfg.MchID,
			"nonce_str":    nonce,
			"prepay_id":    prepayID,
			"code_url":     codeURL,
			"err_code":     "",
			"err_code_des": "",
		})

		xml := fmt.Sprintf("<xml><return_code>SUCCESS</return_code><return_msg>OK</return_msg><result_code>SUCCESS</result_code><appid>%s</appid><mch_id>%s</mch_id><nonce_str>%s</nonce_str><sign>%s</sign><prepay_id>%s</prepay_id><code_url>%s</code_url></xml>", cfg.AppID, cfg.MchID, nonce, sign, prepayID, codeURL)
		_, _ = w.Write([]byte(xml))
	}))
	defer server.Close()

	s.unifiedOrderURL = server.URL

	resp, err := s.CreateNativePay("order-1", 520, "body")
	if err != nil {
		t.Fatalf("create native pay error: %v", err)
	}
	if resp == nil || resp.ReturnCode != "SUCCESS" || resp.PrepayID != "prepay123" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestCreateNativePayRealModeInvalidConfig(t *testing.T) {
	s := NewService(&Config{MockEnabled: false})
	_, err := s.CreateNativePay("order-1", 520, "body")
	if err == nil || !strings.Contains(err.Error(), "微信支付配置不完整") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateNativePayInvalidInput(t *testing.T) {
	s := newTestService()

	if _, err := s.CreateNativePay("", 520, "body"); err == nil {
		t.Fatalf("expected empty order error")
	}
	if _, err := s.CreateNativePay("order-1", 0, "body"); err == nil {
		t.Fatalf("expected amount error")
	}
	if _, err := s.CreateNativePay("order-1", 520, ""); err == nil {
		t.Fatalf("expected body error")
	}
}

func TestParseNotifyRequest(t *testing.T) {
	s := newTestService()
	params := map[string]string{
		"return_code":    "SUCCESS",
		"appid":          "wx-app",
		"mch_id":         "mch-1",
		"device_info":    "d",
		"nonce_str":      "n",
		"out_trade_no":   "o-1",
		"transaction_id": "tx-1",
		"total_fee":      "100",
		"fee_type":       "CNY",
		"time_end":       "20260101010101",
		"attach":         "a",
		"is_subscribe":   "Y",
	}
	sign := s.GenerateMD5Sign(params)
	xml := "<xml>" +
		"<return_code>SUCCESS</return_code>" +
		"<appid>wx-app</appid>" +
		"<mch_id>mch-1</mch_id>" +
		"<device_info>d</device_info>" +
		"<nonce_str>n</nonce_str>" +
		"<sign>" + sign + "</sign>" +
		"<out_trade_no>o-1</out_trade_no>" +
		"<transaction_id>tx-1</transaction_id>" +
		"<total_fee>100</total_fee>" +
		"<fee_type>CNY</fee_type>" +
		"<time_end>20260101010101</time_end>" +
		"<attach>a</attach>" +
		"<is_subscribe>Y</is_subscribe>" +
		"</xml>"

	notify, err := s.ParseNotifyRequest(xml)
	if err != nil {
		t.Fatalf("parse notify error: %v", err)
	}
	if notify == nil || notify.OutTradeNo != "o-1" {
		t.Fatalf("unexpected notify")
	}
}

func TestBuildNotifyResponseAndNonce(t *testing.T) {
	x := BuildNotifyResponse()
	if x != "" {
		t.Fatalf("unexpected notify response: %s", x)
	}
	n1 := generateNonce()
	n2 := generateNonce()
	if n1 == "" || n2 == "" || n1 == n2 {
		t.Fatalf("nonce invalid")
	}
}
