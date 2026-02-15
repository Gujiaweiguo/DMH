/**
 * String utilities
 */

export const truncate = (str, maxLength, suffix = '...') => {
  if (!str) return ''
  if (str.length <= maxLength) return str
  return str.substring(0, maxLength - suffix.length) + suffix
}

export const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase()
}

export const capitalizeWords = (str) => {
  if (!str) return ''
  return str.split(' ').map(capitalize).join(' ')
}

export const trim = (str) => {
  if (!str) return ''
  return str.trim()
}

export const isEmpty = (str) => {
  return !str || str.trim().length === 0
}

export const isNotEmpty = (str) => !isEmpty(str)

export const maskPhone = (phone) => {
  if (!phone || phone.length !== 11) return phone
  return phone.substring(0, 3) + '****' + phone.substring(7)
}

export const maskEmail = (email) => {
  if (!email || !email.includes('@')) return email
  const [local, domain] = email.split('@')
  const maskedLocal = local.length > 2 
    ? local[0] + '***' + local[local.length - 1] 
    : local
  return maskedLocal + '@' + domain
}

export const maskIdCard = (idCard) => {
  if (!idCard || idCard.length < 8) return idCard
  return idCard.substring(0, 4) + '**********' + idCard.substring(idCard.length - 4)
}

export const formatNumber = (num, decimals = 0) => {
  const n = Number(num)
  if (isNaN(n)) return '0'
  return n.toFixed(decimals)
}

export const formatCurrency = (amount, symbol = '¥') => {
  const n = Number(amount)
  if (isNaN(n)) return symbol + '0.00'
  return symbol + n.toFixed(2)
}

export const formatPercent = (value, decimals = 1) => {
  const n = Number(value)
  if (isNaN(n)) return '0%'
  return (n * 100).toFixed(decimals) + '%'
}

export const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + units[i]
}

export const slugify = (str) => {
  if (!str) return ''
  return str
    .toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

export const escapeHtml = (str) => {
  if (!str) return ''
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  }
  return str.replace(/[&<>"']/g, m => map[m])
}

export const unescapeHtml = (str) => {
  if (!str) return ''
  const map = {
    '&amp;': '&',
    '&lt;': '<',
    '&gt;': '>',
    '&quot;': '"',
    '&#039;': "'"
  }
  return str.replace(/&(amp|lt|gt|quot|#039);/g, m => map[m])
}
