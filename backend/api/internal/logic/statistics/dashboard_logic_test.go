package statistics

import (
	"context"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestGetDashboardStatsLogic_InvalidBrandId(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	req := &types.GetDashboardStatsReq{
		BrandId: 0,
		Period:  "month",
	}

	resp, err := logic.GetDashboardStats(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "brandId is required")
}

func TestDateFilterQuery(t *testing.T) {
	tests := []struct {
		name     string
		period   string
		expected string
	}{
		{
			name:     "默认周期",
			period:   "",
			expected: "DATE_SUB(NOW(), INTERVAL 90 DAY)",
		},
		{
			name:     "周周期",
			period:   "week",
			expected: "DATE_SUB(NOW(), INTERVAL 7 DAY)",
		},
		{
			name:     "月周期",
			period:   "month",
			expected: "DATE_SUB(NOW(), INTERVAL 30 DAY)",
		},
		{
			name:     "年周期",
			period:   "year",
			expected: "DATE_SUB(NOW(), INTERVAL 365 DAY)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, dateFilterQuery(&types.GetDashboardStatsReq{Period: tt.period}))
		})
	}
}

func TestGetDashboardStatsLogic_Validation(t *testing.T) {
	tests := []struct {
		name    string
		brandId int64
		wantErr bool
		errMsg  string
	}{
		{
			name:    "zero brandId",
			brandId: 0,
			wantErr: true,
			errMsg:  "brandId is required",
		},
		{
			name:    "negative brandId",
			brandId: -1,
			wantErr: true,
			errMsg:  "brandId is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svcCtx := &svc.ServiceContext{}
			logic := NewGetDashboardStatsLogic(ctx, svcCtx)

			req := &types.GetDashboardStatsReq{
				BrandId: tt.brandId,
				Period:  "month",
			}

			resp, err := logic.GetDashboardStats(req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetDashboardStatsLogic_Construct(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	assert.NotNil(t, logic)
	assert.NotNil(t, logic.Logger)
	assert.NotNil(t, logic.ctx)
	assert.NotNil(t, logic.svcCtx)
}
