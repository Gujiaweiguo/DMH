package statistics

import (
	"dmh/api/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
