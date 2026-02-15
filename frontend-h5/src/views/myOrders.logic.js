/**
 * MyOrders business logic
 */

export const PHONE_REGEX = /^1[3-9]\d{9}$/

export const isValidPhone = (phone) => PHONE_REGEX.test(phone)

export const STATUS_TEXT_MAP = {
  paid: '已支付',
  pending: '待支付',
  cancelled: '已取消'
}

export const statusText = (status) => STATUS_TEXT_MAP[status] || status

export const formatDate = (time) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export const formatAmount = (amount) => {
  const num = Number(amount) || 0
  return num.toFixed(2)
}

export const buildCampaignMap = (campaigns) => {
  const map = {}
  if (Array.isArray(campaigns)) {
    campaigns.forEach(c => {
      map[c.id] = c
    })
  }
  return map
}

export const getCampaignName = (campaignId, campaignMap) => {
  return campaignMap[campaignId]?.name || '未知活动'
}

export const filterOrdersByPhone = (orders, phone) => {
  if (!Array.isArray(orders)) return []
  return orders.filter(order => order.phone === phone)
}

export const savePhoneToStorage = (phone) => {
  try {
    localStorage.setItem('dmh_my_phone', phone)
    return true
  } catch (e) {
    return false
  }
}

export const loadPhoneFromStorage = () => {
  try {
    return localStorage.getItem('dmh_my_phone') || ''
  } catch (e) {
    return ''
  }
}

export const hasFormData = (order) => {
  return !!(order && order.formData && Object.keys(order.formData).length > 0)
}
