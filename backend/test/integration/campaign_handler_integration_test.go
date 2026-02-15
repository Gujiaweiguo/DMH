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

type CampaignHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
	campaignId int64
}

func (suite *CampaignHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
	suite.createTestCampaign()
}

func (suite *CampaignHandlerIntegrationTestSuite) loginAsAdmin() {
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

	body, _ := io.ReadAll(resp.Body)
	var loginResp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(body, &loginResp)

	if loginResp.Token == "" {
		suite.T().Skipf("登录失败: %s", string(body))
		return
	}

	suite.adminToken = loginResp.Token
	suite.T().Log("✓ Admin 登录成功")
}

func (suite *CampaignHandlerIntegrationTestSuite) createCampaignWithName(name string) (int64, int, []byte) {
	now := time.Now()
	createReq := map[string]interface{}{
		"brandId":     1,
		"name":        name,
		"description": "用于集成测试的活动",
		"rewardRule":  10,
		"startTime":   now.Add(-1 * time.Hour).Format("2006-01-02T15:04:05"),
		"endTime":     now.Add(24 * time.Hour).Format("2006-01-02T15:04:05"),
		"formFields": []map[string]interface{}{
			{
				"type":     "text",
				"name":     "name",
				"label":    "姓名",
				"required": true,
			},
		},
	}

	reqBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/campaigns", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, resp.StatusCode, body
	}

	var createResp struct {
		Id int64 `json:"id"`
	}
	_ = json.Unmarshal(body, &createResp)

	return createResp.Id, resp.StatusCode, body
}

func (suite *CampaignHandlerIntegrationTestSuite) createTestCampaign() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
		return
	}

	id, status, body := suite.createCampaignWithName(fmt.Sprintf("集成测试活动_%d", time.Now().Unix()))
	if status != http.StatusOK || id == 0 {
		suite.T().Skipf("创建测试活动失败，status=%d, body=%s", status, string(body))
		return
	}

	suite.campaignId = id
	suite.T().Logf("✓ 测试活动创建成功，ID: %d", suite.campaignId)
}

func (suite *CampaignHandlerIntegrationTestSuite) TearDownSuite() {
	if suite.campaignId > 0 && suite.adminToken != "" {
		url := fmt.Sprintf("%s/api/v1/campaigns/%d", suite.baseURL, suite.campaignId)
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)
		_, _ = suite.httpClient.Do(req)
	}
}

func (suite *CampaignHandlerIntegrationTestSuite) Test_1_GetCampaignsList() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/campaigns?page=1&pageSize=10", nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.T().Logf("✓ 活动列表查询成功，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func (suite *CampaignHandlerIntegrationTestSuite) Test_2_GetCampaignById() {
	if suite.adminToken == "" || suite.campaignId == 0 {
		suite.T().Skip("未登录或活动不存在，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/campaigns/%d", suite.baseURL, suite.campaignId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.T().Logf("✓ 活动详情查询成功，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func (suite *CampaignHandlerIntegrationTestSuite) Test_3_UpdateCampaign() {
	if suite.adminToken == "" || suite.campaignId == 0 {
		suite.T().Skip("未登录或活动不存在，跳过测试")
	}

	updateReq := map[string]interface{}{
		"name": fmt.Sprintf("更新后的活动_%d", time.Now().Unix()),
	}
	reqBody, _ := json.Marshal(updateReq)

	url := fmt.Sprintf("%s/api/v1/campaigns/%d", suite.baseURL, suite.campaignId)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.T().Logf("✓ 活动更新成功，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func (suite *CampaignHandlerIntegrationTestSuite) Test_4_GetPaymentQrcode() {
	if suite.adminToken == "" || suite.campaignId == 0 {
		suite.T().Skip("未登录或活动不存在，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/campaigns/%d/payment-qrcode", suite.baseURL, suite.campaignId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.T().Logf("✓ 活动支付二维码获取成功，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func (suite *CampaignHandlerIntegrationTestSuite) Test_5_CreateCampaignInvalidTime() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	createReq := map[string]interface{}{
		"brandId":     1,
		"name":        "非法时间活动",
		"description": "测试时间格式校验",
		"rewardRule":  10,
		"startTime":   "invalid-time",
		"endTime":     time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05"),
		"formFields": []map[string]interface{}{
			{
				"type":     "text",
				"name":     "name",
				"label":    "姓名",
				"required": true,
			},
		},
	}

	reqBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/campaigns", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	suite.NotEqual(http.StatusOK, resp.StatusCode)
	suite.T().Logf("✓ 非法时间活动创建被拒绝，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func (suite *CampaignHandlerIntegrationTestSuite) Test_6_GetCampaignNotFound() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/campaigns/999999", suite.baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	allowedStatus := []int{http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, resp.StatusCode)
	suite.T().Logf("✓ 不存在活动查询完成，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func (suite *CampaignHandlerIntegrationTestSuite) Test_7_DeleteCampaign() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	tempID, status, body := suite.createCampaignWithName(fmt.Sprintf("待删除活动_%d", time.Now().UnixNano()))
	if status != http.StatusOK || tempID == 0 {
		suite.T().Skipf("创建待删除活动失败，status=%d, body=%s", status, string(body))
	}

	url := fmt.Sprintf("%s/api/v1/campaigns/%d", suite.baseURL, tempID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.T().Logf("✓ 活动删除成功，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
}

func (suite *CampaignHandlerIntegrationTestSuite) Test_8_UnauthorizedAccess() {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/campaigns?page=1&pageSize=10", nil)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
	suite.T().Logf("✓ 无 Token 访问被拒绝，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func TestCampaignHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CampaignHandlerIntegrationTestSuite))
}
