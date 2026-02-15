/**
 * Designer components shared logic
 */

export const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

export const formatCountdown = (diff) => {
  if (diff <= 0) return { days: 0, hours: 0, minutes: 0, seconds: 0 }
  
  const days = Math.floor(diff / (24 * 60 * 60 * 1000))
  const hours = Math.floor((diff % (24 * 60 * 60 * 1000)) / (60 * 60 * 1000))
  const minutes = Math.floor((diff % (60 * 60 * 1000)) / (60 * 1000))
  const seconds = Math.floor((diff % (60 * 1000)) / 1000)
  
  return { days, hours, minutes, seconds }
}

export const padZero = (num) => String(num).padStart(2, '0')

export const isExpired = (endTime) => {
  if (!endTime) return false
  return new Date(endTime) < new Date()
}

export const getImageUrl = (url) => url || ''

export const truncateText = (text, maxLength) => {
  if (!text) return ''
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

export const parseJsonSafe = (str, defaultValue = {}) => {
  try {
    return JSON.parse(str || '{}')
  } catch {
    return defaultValue
  }
}

export const buildQueryString = (params) => {
  return Object.entries(params)
    .filter(([_, v]) => v !== null && v !== undefined && v !== '')
    .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
    .join('&')
}

export const DEFAULT_COLORS = {
  primary: '#4f46e5',
  secondary: '#6366f1',
  success: '#10b981',
  warning: '#f59e0b',
  danger: '#ef4444',
  text: '#333333',
  textSecondary: '#666666',
  background: '#ffffff'
}

export const validateRequired = (value, fieldName) => {
  if (!value || (typeof value === 'string' && !value.trim())) {
    return { valid: false, error: `${fieldName}不能为空` }
  }
  return { valid: true, error: null }
}

export const validatePhone = (phone) => {
  if (!phone) return { valid: false, error: '请输入手机号' }
  if (!/^1[3-9]\d{9}$/.test(phone)) return { valid: false, error: '手机号格式不正确' }
  return { valid: true, error: null }
}

export const validateEmail = (email) => {
  if (!email) return { valid: true, error: null }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return { valid: false, error: '邮箱格式不正确' }
  return { valid: true, error: null }
}
