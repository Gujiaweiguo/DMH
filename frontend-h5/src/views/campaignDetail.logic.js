/**
 * CampaignDetail business logic
 */

export const formatTime = (time) => {
  if (!time) return ''
  return time.replace('T', ' ').replace('Z', '').substring(0, 16)
}

export const formatReward = (reward) => {
  const num = Number(reward) || 0
  return num.toFixed(2)
}

export const formatPaymentAmount = (amount) => {
  const num = Number(amount) || 0
  return num.toFixed(2)
}

export const getPaymentTypeText = (type) => {
  return type === 'deposit' ? '订金支付' : '全款支付'
}

export const getPaymentLabel = (type) => {
  return type === 'deposit' ? '订金金额' : '支付金额'
}

export const isPaymentRequired = (config) => {
  return !!(config && config.requirePayment)
}

export const calculateConversionRate = (participantCount, orderCount) => {
  if (!participantCount || participantCount === 0) return '0.0'
  return ((orderCount || 0) / participantCount * 100).toFixed(1)
}

export const isBrandAdmin = () => {
  const userRole = localStorage.getItem('dmh_user_role')
  return userRole === 'brand_admin'
}

export const buildSourceData = (query) => ({
  c_id: query.c_id || '',
  u_id: query.u_id || ''
})

export const saveSourceToStorage = (source) => {
  try {
    localStorage.setItem('dmh_source', JSON.stringify(source))
    return true
  } catch (e) {
    return false
  }
}

export const buildPosterDownloadName = (campaignId) => `poster_${campaignId}.png`

export const buildVerifyRoute = (campaignId) => ({
  path: '/verify',
  query: { campaignId }
})

export const buildPaymentRoute = (campaignId) => `/campaign/${campaignId}/payment`

export const buildFormRoute = (campaignId, isBrandAdmin) => {
  if (isBrandAdmin) {
    return { path: '/verify', query: { campaignId } }
  }
  return `/campaign/${campaignId}/form`
}
