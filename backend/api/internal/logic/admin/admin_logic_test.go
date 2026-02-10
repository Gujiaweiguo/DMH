package admin

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AdminLogicTestSuite struct {
	suite.Suite
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (suite *AdminLogicTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{})

	suite.ctx = context.Background()
	suite.svcCtx = &svc.ServiceContext{
		DB: db,
	}
}

func (suite *AdminLogicTestSuite) TestGetUsers() {
	users := []model.User{
		{Id: 1, Username: "admin", Phone: "13800138000", Status: "active"},
		{Id: 2, Username: "user1", Phone: "13800138001", Status: "active"},
		{Id: 3, Username: "user2", Phone: "13800138002", Status: "disabled"},
	}
	for _, user := range users {
		suite.svcCtx.DB.Create(&user)
	}

	logic := NewGetUsersLogic(suite.ctx, suite.svcCtx)

	tests := []struct {
		name     string
		req      *types.AdminGetUsersReq
		validate func(resp *types.AdminUserListResp)
	}{
		{
			name: "查询所有用户",
			req: &types.AdminGetUsersReq{
				Page:     1,
				PageSize: 10,
			},
			validate: func(resp *types.AdminUserListResp) {
				assert.Equal(suite.T(), int64(3), resp.Total)
				assert.Len(suite.T(), resp.Users, 3)
			},
		},
		{
			name: "按状态筛选",
			req: &types.AdminGetUsersReq{
				Page:     1,
				PageSize: 10,
				Status:   "active",
			},
			validate: func(resp *types.AdminUserListResp) {
				assert.Equal(suite.T(), int64(2), resp.Total)
			},
		},
		{
			name: "关键词搜索",
			req: &types.AdminGetUsersReq{
				Page:     1,
				PageSize: 10,
				Keyword:  "admin",
			},
			validate: func(resp *types.AdminUserListResp) {
				assert.Equal(suite.T(), int64(1), resp.Total)
			},
		},
		{
			name: "分页查询",
			req: &types.AdminGetUsersReq{
				Page:     1,
				PageSize: 2,
			},
			validate: func(resp *types.AdminUserListResp) {
				assert.Equal(suite.T(), int64(3), resp.Total)
				assert.Len(suite.T(), resp.Users, 2)
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			resp, err := logic.GetUsers(tt.req)
			assert.NoError(t, err)
			tt.validate(resp)
		})
	}
}

func TestAdminLogicTestSuite(t *testing.T) {
	suite.Run(t, new(AdminLogicTestSuite))
}
