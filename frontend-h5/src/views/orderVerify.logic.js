/**
 * OrderVerify business logic
 */

export const VERIFICATION_STATUS_MAP = {
  'unverified': '未核销',
  'verified': '已核销',
  'cancelled': '已取消'
}

export const getVerificationStatusText = (status) => {
  return VERIFICATION_STATUS_MAP[status] || 'unknown'
}

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

export const formatFormData = (formDataStr) => {
  try {
    const formData = JSON.parse(formDataStr || '{}')
    return Object.entries(formData)
      .map(([key, value]) => `${key}: ${value}`)
      .join(', ')
  } catch (error) {
    return formDataStr || ''
  }
}

export const canVerify = (status) => status === 'unverified'

export const canUnverify = (status) => status === 'verified'

export const isValidVerificationCode = (code) => {
  return !!(code && code.trim().length > 0)
}

export const buildScanButtonText = (scanning) => {
  return scanning ? '扫描中...' : '扫码'
}
