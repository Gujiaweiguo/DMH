/**
 * Analytics business logic
 */

export const PERIOD_OPTIONS = [
  { value: 'today', label: '今日' },
  { value: 'week', label: '本周' },
  { value: 'month', label: '本月' },
  { value: 'quarter', label: '本季度' }
]

export const getDefaultCoreMetrics = () => ({
  totalRevenue: 0,
  totalOrders: 0,
  activePromoters: 0,
  avgOrderValue: 0
})

export const getDefaultChartData = () => ({
  orders: []
})

export const getDefaultFunnelData = () => [
  { label: '页面访问', value: 10000, percentage: 100 },
  { label: '点击活动', value: 3500, percentage: 35 },
  { label: '填写表单', value: 1200, percentage: 12 },
  { label: '提交订单', value: 800, percentage: 8 },
  { label: '支付成功', value: 650, percentage: 6.5 }
]

export const EXPORT_TYPES = {
  orders: '订单数据',
  promoters: '推广员数据',
  campaigns: '活动数据',
  all: '完整报表'
}

export const getExportTypeLabel = (type) => EXPORT_TYPES[type] || type

export const formatMetricValue = (value, prefix = '') => {
  const num = Number(value) || 0
  return prefix + num.toLocaleString()
}

export const calculateChangeIndicator = (current, previous) => {
  if (!previous || previous === 0) {
    return { value: 0, type: 'neutral' }
  }
  const change = ((current - previous) / previous) * 100
  return {
    value: Math.abs(change).toFixed(1),
    type: change >= 0 ? 'positive' : 'negative'
  }
}

export const getRankingNumberStyle = (index) => {
  if (index === 0) return 'gold'
  if (index === 1) return 'silver'
  if (index === 2) return 'bronze'
  return 'default'
}
