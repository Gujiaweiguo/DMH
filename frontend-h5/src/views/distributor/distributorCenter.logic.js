/**
 * DistributorCenter business logic
 */

export const getStatusType = (status) => {
  const types = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return types[status] || 'default'
}

export const getStatusText = (status) => {
  const texts = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return texts[status] || status
}

export const getDefaultStatistics = () => ({
  totalEarnings: 0,
  totalOrders: 0,
  subordinatesCount: 0,
  todayEarnings: 0,
  monthEarnings: 0,
  balance: 0
})

export const clearDistributorAuth = () => {
  localStorage.removeItem('dmh_token')
  localStorage.removeItem('dmh_user_role')
  localStorage.removeItem('dmh_user_info')
  localStorage.removeItem('dmh_current_brand_id')
}

export const buildPromotionRoute = (brandId) => `/distributor/promotion?brandId=${brandId}`

export const buildRewardsRoute = (brandId) => `/distributor/rewards?brandId=${brandId}`

export const buildSubordinatesRoute = (brandId) => `/distributor/subordinates?brandId=${brandId}`

export const buildStatisticsRoute = (brandId) => `/distributor/statistics?brandId=${brandId}`

export const buildWithdrawalsRoute = (brandId) => `/distributor/withdrawals?brandId=${brandId}`

export const buildPosterRoute = (campaignId) => `/poster-generator/${campaignId}`

export const formatMoney = (value) => {
  const num = Number(value) || 0
  return num.toFixed(2)
}
