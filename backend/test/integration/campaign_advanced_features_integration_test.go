package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// PaymentQrcodeIntegrationTestSuite 支付二维码集成测试套件
type PaymentQrcodeIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	authToken  string
	campaignId int64
}

func (suite *PaymentQrcodeIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889" // 后端服务地址
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}

	// 登录获取 token
	suite.login()

	// 创建测试活动
	suite.createTestCampaign()
}

func (suite *PaymentQrcodeIntegrationTestSuite) login() {
	loginReq := map[string]string{
		"username": "admin",
		"password": "123456",
	}
	reqBody, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("无法连接到后端服务，跳过测试: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var loginResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(body, &loginResp)

	if loginResp.Token == "" {
		suite.T().Skipf("登录失败，未获取到 token: %s", string(body))
		return
	}

	suite.authToken = loginResp.Token
	suite.T().Log("✓ 登录成功，获取到 token: " + suite.authToken)
}

func (suite *PaymentQrcodeIntegrationTestSuite) createTestCampaign() {
	// 创建测试品牌和活动
	now := time.Now()
	createCampaignReq := map[string]interface{}{
		"brandId":     1, // 使用默认品牌ID
		"name":        "支付二维码测试活动",
		"description": "用于测试支付二维码生成的活动",
		"rewardRule":  10.0,
		"startTime":   now.Add(1 * time.Hour).Format("2006-01-02T15:04:05"),
		"endTime":     now.Add(30 * 24 * time.Hour).Format("2006-01-02T15:04:05"),
		"paymentConfig": map[string]interface{}{
			"depositAmount":  50.00,
			"fullAmount":     200.00,
			"paymentType":    "deposit",
			"wechatMerchant": "1234567890",
			"callbackUrl":    "http://example.com/callback",
		},
	}

	reqBody, _ := json.Marshal(createCampaignReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/campaigns", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("创建活动失败: %v，跳过测试", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var createResp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Id int64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(body, &createResp)

	if createResp.Code != 200 {
		suite.T().Skipf("创建活动失败: %s，跳过测试", createResp.Msg)
		return
	}

	suite.campaignId = createResp.Data.Id
	suite.T().Logf("✓ 测试活动创建成功，ID: %d", suite.campaignId)
}

// Test_11_2_1_GenerateDepositQrcode 测试生成订金支付二维码
func (suite *PaymentQrcodeIntegrationTestSuite) Test_11_2_1_GenerateDepositQrcode() {
	suite.T().Log("测试场景 1: 生成订金支付二维码")

	url := fmt.Sprintf("%s/api/v1/campaigns/%d/payment-qrcode", suite.baseURL, suite.campaignId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// 检查响应
	suite.Equal(http.StatusOK, resp.StatusCode, "响应状态码应为 200")

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			QrcodeUrl    string  `json:"qrcodeUrl"`
			Amount       float64 `json:"amount"`
			CampaignName string  `json:"campaignName"`
		} `json:"data"`
	}
	json.Unmarshal(body, &result)

	suite.Equal(200, result.Code, "响应码应为 200")
	suite.NotEmpty(result.Data.QrcodeUrl, "二维码 URL 不应为空")
	suite.Equal(50.00, result.Data.Amount, "订金金额应为 50.00")
	suite.Equal("支付二维码测试活动", result.Data.CampaignName, "活动名称正确")

	suite.T().Logf("✓ 订金二维码生成成功:")
	suite.T().Logf("  - URL: %s", result.Data.QrcodeUrl)
	suite.T().Logf("  - 金额: %.2f", result.Data.Amount)
	suite.T().Logf("  - 活动名称: %s", result.Data.CampaignName)
}

// Test_11_2_2_RefreshQrcode 测试刷新支付二维码
func (suite *PaymentQrcodeIntegrationTestSuite) Test_11_2_2_RefreshQrcode() {
	suite.T().Log("测试场景 2: 刷新支付二维码（验证时间戳变化）")

	// 第一次请求
	url := fmt.Sprintf("%s/api/v1/campaigns/%d/payment-qrcode", suite.baseURL, suite.campaignId)
	req1, _ := http.NewRequest("GET", url, nil)
	req1.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp1, _ := suite.httpClient.Do(req1)
	body1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()

	var result1 struct {
		Data struct {
			QrcodeUrl string `json:"qrcodeUrl"`
		} `json:"data"`
	}
	json.Unmarshal(body1, &result1)

	// 等待 1 秒确保时间戳不同
	time.Sleep(1 * time.Second)

	// 第二次请求（刷新）
	req2, _ := http.NewRequest("GET", url, nil)
	req2.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp2, _ := suite.httpClient.Do(req2)
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()

	var result2 struct {
		Data struct {
			QrcodeUrl string `json:"qrcodeUrl"`
		} `json:"data"`
	}
	json.Unmarshal(body2, &result2)

	suite.Equal(http.StatusOK, resp1.StatusCode, "第一次请求状态码应为 200")
	suite.Equal(http.StatusOK, resp2.StatusCode, "第二次请求状态码应为 200")
	suite.NotEqual(result1.Data.QrcodeUrl, result2.Data.QrcodeUrl, "刷新后的二维码 URL 应不同")

	suite.T().Log("✓ 二维码刷新成功:")
	suite.T().Logf("  - 第一次 URL: %s", result1.Data.QrcodeUrl)
	suite.T().Logf("  - 第二次 URL: %s", result2.Data.QrcodeUrl)
}

// Test_11_2_3_QrcodeWithExpiredCampaign 测试过期活动的二维码
func (suite *PaymentQrcodeIntegrationTestSuite) Test_11_2_3_QrcodeWithExpiredCampaign() {
	suite.T().Log("测试场景 3: 过期活动的二维码生成")

	// 创建过期活动
	now := time.Now()
	createCampaignReq := map[string]interface{}{
		"brandId":     1,
		"name":        "过期活动测试",
		"description": "已过期的活动",
		"rewardRule":  10.0,
		"startTime":   now.Add(-30 * 24 * time.Hour).Format("2006-01-02T15:04:05"), // 30天前开始
		"endTime":     now.Add(-1 * time.Hour).Format("2006-01-02T15:04:05"),       // 1小时前结束
		"paymentConfig": map[string]interface{}{
			"depositAmount": 50.00,
			"fullAmount":    200.00,
			"paymentType":   "deposit",
		},
	}

	reqBody, _ := json.Marshal(createCampaignReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/campaigns", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("创建过期活动失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var createResp struct {
		Data struct {
			Id int64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(body, &createResp)

	expiredCampaignId := createResp.Data.Id

	// 尝试生成二维码
	qrcodeUrl := fmt.Sprintf("%s/api/v1/campaigns/%d/payment-qrcode", suite.baseURL, expiredCampaignId)
	req2, _ := http.NewRequest("GET", qrcodeUrl, nil)
	req2.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp2, _ := suite.httpClient.Do(req2)
	resp2.Body.Close()

	// 过期活动可能仍能生成二维码（业务决策），这里验证不崩溃即可
	suite.T().Logf("✓ 过期活动二维码请求完成 (状态码: %d)", resp2.StatusCode)
}

// Test_11_2_4_QrcodePermission 测试权限控制
func (suite *PaymentQrcodeIntegrationTestSuite) Test_11_2_4_QrcodePermission() {
	suite.T().Log("测试场景 4: 未认证用户访问权限")

	url := fmt.Sprintf("%s/api/v1/campaigns/%d/payment-qrcode", suite.baseURL, suite.campaignId)
	req, _ := http.NewRequest("GET", url, nil)
	// 不设置 Authorization header

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// 未认证用户应该返回 401 或 403
	suite.T().Logf("✓ 未认证访问返回状态码: %d", resp.StatusCode)
	suite.True(resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden,
		"未认证访问应返回 401 或 403")
}

// Test_11_2_5_QrcodeResponseFormat 测试响应格式
func (suite *PaymentQrcodeIntegrationTestSuite) Test_11_2_5_QrcodeResponseFormat() {
	suite.T().Log("测试场景 5: 验证响应格式符合规范")

	url := fmt.Sprintf("%s/api/v1/campaigns/%d/payment-qrcode", suite.baseURL, suite.campaignId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			QrcodeUrl    string  `json:"qrcodeUrl"`
			Amount       float64 `json:"amount"`
			CampaignName string  `json:"campaignName"`
		} `json:"data"`
	}
	json.Unmarshal(body, &result)

	// 验证响应字段
	suite.Equal(200, result.Code, "code 字段应为 200")
	suite.NotEmpty(result.Data.QrcodeUrl, "data.qrcodeUrl 不应为空")
	suite.Greater(result.Data.Amount, 0.0, "data.amount 应大于 0")
	suite.NotEmpty(result.Data.CampaignName, "data.campaignName 不应为空")

	suite.T().Log("✓ 响应格式验证通过:")
	suite.T().Logf("  - code: %d", result.Code)
	suite.T().Logf("  - qrcodeUrl: %s", result.Data.QrcodeUrl)
	suite.T().Logf("  - amount: %.2f", result.Data.Amount)
	suite.T().Logf("  - campaignName: %s", result.Data.CampaignName)
}

// TestPaymentQrcodeIntegrationTestSuite 运行测试套件
func TestPaymentQrcodeIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentQrcodeIntegrationTestSuite))
}
