/**
 * PaymentQrcode business logic
 */

export const formatExpireTime = (timeStr) => {
  if (!timeStr) return '-'
  try {
    const date = new Date(timeStr)
    if (isNaN(date.getTime())) return timeStr
    
    const now = new Date()
    const diff = date - now
    
    if (diff < 0) {
      return '已过期'
    }
    
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(minutes / 60)
    
    if (hours > 0) {
      return `${hours}小时${minutes % 60}分钟后过期`
    } else {
      return `${minutes}分钟后过期`
    }
  } catch (e) {
    return timeStr
  }
}

export const formatAmount = (amount) => {
  const num = Number(amount) || 0
  return num.toFixed(2)
}

export const isExpired = (expireAt) => {
  if (!expireAt) return false
  try {
    const date = new Date(expireAt)
    const now = new Date()
    return date < now
  } catch (e) {
    return false
  }
}

export const getDefaultQrcodeData = () => ({
  qrcodeBase64: null,
  campaignName: null,
  amount: null,
  expireAt: null
})

export const buildQrcodeData = (response) => {
  if (!response || !response.data) {
    return null
  }
  
  return {
    qrcodeBase64: response.data.qrcodeUrl,
    campaignName: response.data.campaignName || '未知活动',
    amount: response.data.amount || 0,
    expireAt: response.data.expireAt || new Date(Date.now() + 2 * 60 * 1000).toISOString()
  }
}

export const validateCampaignId = (id) => {
  if (!id) {
    return { valid: false, error: '活动ID无效' }
  }
  return { valid: true, error: null }
}

export const getMinutesUntilExpiry = (expireAt) => {
  if (!expireAt) return 0
  try {
    const date = new Date(expireAt)
    if (isNaN(date.getTime())) return 0
    
    const now = new Date()
    const diff = date - now
    return Math.max(0, Math.floor(diff / 60000))
  } catch (e) {
    return 0
  }
}
