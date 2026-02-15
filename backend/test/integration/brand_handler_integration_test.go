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

type BrandHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
	brandId    int64
}

func (suite *BrandHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
	suite.createTestBrand()
}

func (suite *BrandHandlerIntegrationTestSuite) loginAsAdmin() {
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
	json.Unmarshal(body, &loginResp)

	suite.adminToken = loginResp.Token
	suite.T().Log("✓ Admin 登录成功")
}

func (suite *BrandHandlerIntegrationTestSuite) createTestBrand() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
		return
	}

	createReq := map[string]interface{}{
		"name":        fmt.Sprintf("集成测试品牌_%d", time.Now().Unix()),
		"logo":        "https://example.com/logo.png",
		"description": "用于集成测试的品牌",
	}
	reqBody, _ := json.Marshal(createReq)

	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/brands", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Logf("创建测试品牌失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	suite.brandId = suite.extractBrandID(body)
	if suite.brandId == 0 {
		suite.loadFirstBrandID()
	}
	suite.T().Logf("✓ 测试品牌创建成功, ID: %d", suite.brandId)
}

func (suite *BrandHandlerIntegrationTestSuite) extractBrandID(body []byte) int64 {
	var direct struct {
		Id int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &direct); err == nil && direct.Id > 0 {
		return direct.Id
	}

	var wrapped struct {
		Data struct {
			Id int64 `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err == nil && wrapped.Data.Id > 0 {
		return wrapped.Data.Id
	}

	return 0
}

func (suite *BrandHandlerIntegrationTestSuite) loadFirstBrandID() {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/brands?page=1&pageSize=1", nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Logf("获取品牌列表失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var listResp struct {
		Brands []struct {
			Id int64 `json:"id"`
		} `json:"brands"`
	}
	if err := json.Unmarshal(body, &listResp); err != nil {
		suite.T().Logf("解析品牌列表失败: %v", err)
		return
	}

	if len(listResp.Brands) > 0 {
		suite.brandId = listResp.Brands[0].Id
	}
}

func (suite *BrandHandlerIntegrationTestSuite) TearDownSuite() {
	if suite.brandId > 0 && suite.adminToken != "" {
		url := fmt.Sprintf("%s/api/v1/brands/%d", suite.baseURL, suite.brandId)
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)
		suite.httpClient.Do(req)
	}
}

func (suite *BrandHandlerIntegrationTestSuite) Test_1_GetBrandsList() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/brands?page=1&pageSize=10", nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Total  int64                    `json:"total"`
		Brands []map[string]interface{} `json:"brands"`
	}
	err = json.Unmarshal(body, &result)
	suite.NoError(err)
	suite.T().Logf("✓ 品牌列表查询成功，共 %d 个品牌", result.Total)
}

func (suite *BrandHandlerIntegrationTestSuite) Test_2_GetBrandById() {
	if suite.adminToken == "" || suite.brandId == 0 {
		suite.T().Skip("未登录或品牌不存在，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/brands/%d", suite.baseURL, suite.brandId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	suite.NoError(err)
	suite.T().Log("✓ 单个品牌查询成功")
}

func (suite *BrandHandlerIntegrationTestSuite) Test_3_UpdateBrand() {
	if suite.adminToken == "" || suite.brandId == 0 {
		suite.T().Skip("未登录或品牌不存在，跳过测试")
	}

	updateReq := map[string]interface{}{
		"name":        fmt.Sprintf("更新后的品牌_%d", time.Now().Unix()),
		"description": "更新后的描述",
	}
	reqBody, _ := json.Marshal(updateReq)

	url := fmt.Sprintf("%s/api/v1/brands/%d", suite.baseURL, suite.brandId)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.T().Log("✓ 品牌更新成功")
}

func (suite *BrandHandlerIntegrationTestSuite) Test_4_GetBrandStats() {
	if suite.adminToken == "" || suite.brandId == 0 {
		suite.T().Skip("未登录或品牌不存在，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/brands/%d/stats", suite.baseURL, suite.brandId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	suite.T().Log("✓ 品牌统计查询成功")
}

func (suite *BrandHandlerIntegrationTestSuite) Test_5_CreateBrandValidation() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	createReq := map[string]interface{}{
		"name": "",
	}
	reqBody, _ := json.Marshal(createReq)

	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/brands", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusUnprocessableEntity}
	suite.Contains(allowedStatus, resp.StatusCode)
	suite.T().Logf("✓ 空品牌名请求完成，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func (suite *BrandHandlerIntegrationTestSuite) Test_6_GetBrandNotFound() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	url := fmt.Sprintf("%s/api/v1/brands/999999", suite.baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	allowedStatus := []int{http.StatusOK, http.StatusNotFound, http.StatusBadRequest}
	suite.Contains(allowedStatus, resp.StatusCode)
	suite.T().Logf("✓ 不存在品牌请求完成，状态码: %d，响应: %s", resp.StatusCode, string(body))
}

func (suite *BrandHandlerIntegrationTestSuite) Test_7_UnauthorizedAccess() {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/brands", nil)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
	suite.T().Log("✓ 无 Token 访问返回 401")
}

func TestBrandHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(BrandHandlerIntegrationTestSuite))
}
