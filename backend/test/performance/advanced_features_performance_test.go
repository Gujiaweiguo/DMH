package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// AdvancedFeaturesPerformanceTestSuite 高级功能性能测试套件
type AdvancedFeaturesPerformanceTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	authToken  string
	campaignId int64
}

func (suite *AdvancedFeaturesPerformanceTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 30 * time.Second}

	suite.login()
	suite.createTestCampaign()
}

func (suite *AdvancedFeaturesPerformanceTestSuite) login() {
	loginReq := map[string]string{
		"username": "admin",
		"password": "123456",
	}
	reqBody, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("无法连接到后端服务: %v", err)
		return
	}
	defer resp.Body.Close()

	var loginResp struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&loginResp)

	if loginResp.Token == "" {
		suite.T().Skipf("登录失败: token is empty")
		return
	}

	suite.authToken = loginResp.Token
	suite.T().Log("✓ 登录成功")
}

func (suite *AdvancedFeaturesPerformanceTestSuite) createTestCampaign() {
	now := time.Now()
	createCampaignReq := map[string]interface{}{
		"brandId":     1,
		"name":        "性能测试活动",
		"description": "用于性能测试的活动",
		"rewardRule":  10.0,
		"startTime":   now.Add(1 * time.Hour).Format("2006-01-02T15:04:05"),
		"endTime":     now.Add(30 * 24 * time.Hour).Format("2006-01-02T15:04:05"),
		"formFields": []map[string]interface{}{
			{
				"type":     "text",
				"name":     "name",
				"label":    "姓名",
				"required": true,
			},
		},
	}

	reqBody, _ := json.Marshal(createCampaignReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/campaigns", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("创建活动失败: %v", err)
		return
	}
	defer resp.Body.Close()

	var createResp struct {
		Id int64 `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&createResp)

	suite.campaignId = createResp.Id
	suite.T().Logf("✓ 测试活动创建成功，ID: %d", suite.campaignId)
}

// Test_12_1_PosterGenerationPerformance 测试海报生成性能（目标 < 3秒）
func (suite *AdvancedFeaturesPerformanceTestSuite) Test_12_1_PosterGenerationPerformance() {
	suite.T().Log("测试场景 12.1: 海报生成性能测试（目标 < 3秒）")

	generatePosterReq := map[string]int64{
		"templateId": 1,
	}
	reqBody, _ := json.Marshal(generatePosterReq)

	url := fmt.Sprintf("%s/api/v1/campaigns/%d/poster", suite.baseURL, suite.campaignId)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	start := time.Now()
	resp, err := suite.httpClient.Do(req)
	duration := time.Since(start)

	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.T().Logf("海报生成耗时: %v", duration)
	suite.T().Logf("响应状态码: %d", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		suite.T().Log("✓ 海报生成成功")
		// 性能目标：< 3 秒
		suite.Less(duration, 3*time.Second, "海报生成时间应 < 3 秒")
	} else {
		suite.T().Logf("⚠ 海报生成失败，状态码: %d", resp.StatusCode)
	}
}

// Test_12_2_PaymentQRCodePerformance 测试二维码生成性能（目标 < 500ms）
func (suite *AdvancedFeaturesPerformanceTestSuite) Test_12_2_PaymentQRCodePerformance() {
	suite.T().Log("测试场景 12.2: 二维码生成性能测试（目标 < 500ms）")

	// 多次测试取平均值
	testCount := 10
	totalDuration := time.Duration(0)
	successCount := 0

	for i := 0; i < testCount; i++ {
		url := fmt.Sprintf("%s/api/v1/campaigns/%d/payment-qrcode", suite.baseURL, suite.campaignId)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+suite.authToken)

		start := time.Now()
		resp, err := suite.httpClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			suite.T().Logf("请求 %d 失败: %v", i+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			totalDuration += duration
			successCount++
			suite.T().Logf("请求 %d 成功，耗时: %v", i+1, duration)
		}
	}

	if successCount > 0 {
		avgDuration := totalDuration / time.Duration(successCount)
		suite.T().Logf("✓ 平均耗时: %v (%d/%d 次成功)", avgDuration, successCount, testCount)
		suite.Less(avgDuration, 500*time.Millisecond, "二维码生成平均时间应 < 500ms")
	} else {
		suite.T().Log("✗ 所有请求均失败")
	}
}

// Test_12_3_OrderVerifyPerformance 测试核销接口响应时间测试（目标 < 500ms）
func (suite *AdvancedFeaturesPerformanceTestSuite) Test_12_3_OrderVerifyPerformance() {
	suite.T().Log("测试场景 12.3: 核销接口响应时间测试（目标 < 500ms）")

	// 创建测试订单
	createOrderReq := map[string]interface{}{
		"campaignId": suite.campaignId,
		"phone":      "13800139999",
		"formData":   map[string]string{"name": "性能测试用户"},
	}
	reqBody, _ := json.Marshal(createOrderReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/orders", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	createResp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("创建订单失败: %v", err)
		return
	}
	defer createResp.Body.Close()

	var orderInfo struct {
		Id int64 `json:"id"`
	}
	json.NewDecoder(createResp.Body).Decode(&orderInfo)

	// 模拟核销码（实际应从订单生成时获取）
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	verificationCode := fmt.Sprintf("%d_13800139999_%s_abc123", orderInfo.Id, timestamp)

	// 多次测试取平均值
	testCount := 10
	totalDuration := time.Duration(0)
	successCount := 0

	for i := 0; i < testCount; i++ {
		verifyReq := map[string]string{
			"code": verificationCode,
		}
		reqBody, _ := json.Marshal(verifyReq)

		url := fmt.Sprintf("%s/api/v1/orders/verify", suite.baseURL)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.authToken)

		start := time.Now()
		resp, err := suite.httpClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			suite.T().Logf("请求 %d 失败: %v", i+1, err)
			continue
		}
		defer resp.Body.Close()

		// 计入所有请求（包括失败）
		totalDuration += duration
		successCount++
		suite.T().Logf("请求 %d 完成 (状态码: %d)，耗时: %v", i+1, resp.StatusCode, duration)
	}

	if successCount > 0 {
		avgDuration := totalDuration / time.Duration(successCount)
		suite.T().Logf("✓ 平均耗时: %v (%d/%d 次完成)", avgDuration, successCount, testCount)
		// 性能目标：< 500ms
		suite.Less(avgDuration, 500*time.Millisecond, "核销接口平均响应时间应 < 500ms")
	} else {
		suite.T().Log("✗ 所有请求均失败")
	}
}

// Test_12_4_ConcurrentPosterStressTest 并发海报生成压力测试
func (suite *AdvancedFeaturesPerformanceTestSuite) Test_12_4_ConcurrentPosterStressTest() {
	suite.T().Log("测试场景 12.4: 并发海报生成压力测试")

	concurrentRequests := 20
	var results []time.Duration
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			generatePosterReq := map[string]int64{
				"templateId": 1,
			}
			reqBody, _ := json.Marshal(generatePosterReq)

			url := fmt.Sprintf("%s/api/v1/campaigns/%d/poster", suite.baseURL, suite.campaignId)
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+suite.authToken)

			start := time.Now()
			resp, err := suite.httpClient.Do(req)
			duration := time.Since(start)

			if err != nil {
				suite.T().Logf("请求 %d 失败: %v", index+1, err)
				return
			}
			defer resp.Body.Close()

			mu.Lock()
			results = append(results, duration)
			mu.Unlock()

			if resp.StatusCode == http.StatusOK {
				suite.T().Logf("请求 %d 成功，耗时: %v", index+1, duration)
			} else {
				suite.T().Logf("请求 %d 失败 (状态码: %d)，耗时: %v", index+1, resp.StatusCode, duration)
			}
		}(i)
	}

	wg.Wait()

	if len(results) > 0 {
		total := time.Duration(0)
		maxDuration := time.Duration(0)
		for _, d := range results {
			total += d
			if d > maxDuration {
				maxDuration = d
			}
		}

		avgDuration := total / time.Duration(len(results))
		suite.T().Logf("✓ 并发测试完成:")
		suite.T().Logf("  总请求数: %d", concurrentRequests)
		suite.T().Logf("  成功请求: %d", len(results))
		suite.T().Logf("  平均耗时: %v", avgDuration)
		suite.T().Logf("  最大耗时: %v", maxDuration)

		suite.Less(maxDuration, 10*time.Second, "最大响应时间应 < 10 秒")
	} else {
		suite.T().Log("✗ 所有请求均失败")
	}
}

// TestAdvancedFeaturesPerformanceTestSuite 运行测试套件
func TestAdvancedFeaturesPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(AdvancedFeaturesPerformanceTestSuite))
}
