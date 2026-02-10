package poster

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PosterServiceTestSuite struct {
	suite.Suite
	service *Service
	testDir string
}

func (suite *PosterServiceTestSuite) SetupTest() {
	suite.testDir = "/tmp/test_posters"
	os.RemoveAll(suite.testDir)

	suite.service = NewService(suite.testDir, "http://localhost:8888")
}

func (suite *PosterServiceTestSuite) TearDownTest() {
	os.RemoveAll(suite.testDir)
}

func (suite *PosterServiceTestSuite) TestNewService() {
	tests := []struct {
		name      string
		posterDir string
		baseURL   string
		wantPanic bool
	}{
		{
			name:      "正常创建",
			posterDir: "/tmp/test_posters_new",
			baseURL:   "http://localhost:8888",
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					NewService(tt.posterDir, tt.baseURL)
				})
			} else {
				service := NewService(tt.posterDir, tt.baseURL)
				assert.NotNil(t, service)
				assert.Equal(t, tt.posterDir, service.posterDir)
				assert.Equal(t, tt.baseURL, service.baseURL)

				// 验证目录已创建
				_, err := os.Stat(tt.posterDir)
				assert.NoError(t, err)

				// 清理
				os.RemoveAll(tt.posterDir)
			}
		})
	}
}

func (suite *PosterServiceTestSuite) TestGenerateUniqueID() {
	id1 := generateUniqueID()
	id2 := generateUniqueID()

	assert.NotZero(suite.T(), id1)
	assert.NotZero(suite.T(), id2)
	assert.NotEqual(suite.T(), id1, id2)
}

func (suite *PosterServiceTestSuite) TestGenerateQRCodeAsBase64() {
	result, err := suite.service.GenerateQRCodeAsBase64("test data")

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
	assert.Contains(suite.T(), result, "data:image/png;base64,")
}

func (suite *PosterServiceTestSuite) TestGenerateQRCode() {
	path, err := suite.service.generateQRCode("https://example.com/test")

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), path)
	assert.Contains(suite.T(), path, "/tmp/qrcode_")

	// 验证文件已创建
	_, err = os.Stat(path)
	assert.NoError(suite.T(), err)

	// 清理
	os.Remove(path)
}

func (suite *PosterServiceTestSuite) TestGenerateCampaignPoster() {
	tests := []struct {
		name            string
		campaignName    string
		campaignDesc    string
		distributorName string
		qrcodeData      string
		wantErr         bool
	}{
		{
			name:            "正常生成海报",
			campaignName:    "Test Campaign",
			campaignDesc:    "Test Description",
			distributorName: "Test Distributor",
			qrcodeData:      "https://example.com/campaign/1",
			wantErr:         false,
		},
		{
			name:            "空描述",
			campaignName:    "Test Campaign",
			campaignDesc:    "",
			distributorName: "Test Distributor",
			qrcodeData:      "https://example.com/campaign/2",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			url, err := suite.service.GenerateCampaignPoster(
				tt.campaignName,
				tt.campaignDesc,
				tt.distributorName,
				tt.qrcodeData,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, url)
			assert.Contains(t, url, "http://localhost:8888/posters/poster_")

			// 验证文件已创建
			filename := filepath.Base(url)
			filepath := filepath.Join(suite.testDir, filename)
			_, err = os.Stat(filepath)
			assert.NoError(t, err)
		})
	}
}

func (suite *PosterServiceTestSuite) TestGenerateDistributorPoster() {
	tests := []struct {
		name            string
		distributorName string
		campaignCount   int
		wantErr         bool
	}{
		{
			name:            "正常生成分销商海报",
			distributorName: "Test Distributor",
			campaignCount:   5,
			wantErr:         false,
		},
		{
			name:            "零活动",
			distributorName: "New Distributor",
			campaignCount:   0,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			url, err := suite.service.GenerateDistributorPoster(
				tt.distributorName,
				tt.campaignCount,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, url)
			assert.Contains(t, url, "http://localhost:8888/posters/distributor_")
		})
	}
}

func TestPosterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PosterServiceTestSuite))
}
