/**
 * CampaignForm business logic
 */

export const PHONE_REGEX = /^1[3-9]\d{9}$/

export const validatePhone = (phone) => {
  if (!phone) {
    return { valid: false, error: '请输入手机号' }
  }
  if (!PHONE_REGEX.test(phone)) {
    return { valid: false, error: '请输入正确的手机号' }
  }
  return { valid: true, error: null }
}

export const validateFormField = (fieldName, value) => {
  if (!value || !value.trim()) {
    return { valid: false, error: `请填写${fieldName}` }
  }
  return { valid: true, error: null }
}

export const validateForm = (phone, formData, formFields) => {
  const errors = []
  
  const phoneResult = validatePhone(phone)
  if (!phoneResult.valid) {
    errors.push(phoneResult.error)
  }
  
  if (Array.isArray(formFields)) {
    for (const field of formFields) {
      const fieldResult = validateFormField(field, formData[field])
      if (!fieldResult.valid) {
        errors.push(fieldResult.error)
      }
    }
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

export const loadSourceFromStorage = () => {
  try {
    const saved = localStorage.getItem('dmh_source')
    if (saved) {
      return JSON.parse(saved)
    }
  } catch (e) {
    console.error('读取来源信息失败', e)
  }
  return { c_id: '', u_id: '' }
}

export const buildOrderPayload = (campaignId, phone, formData, source) => ({
  campaignId: Number(campaignId),
  phone,
  formData: { ...formData },
  referrerId: source.u_id ? Number(source.u_id) : 0
})

export const initializeFormData = (formFields) => {
  const data = {}
  if (Array.isArray(formFields)) {
    formFields.forEach(field => {
      data[field] = ''
    })
  }
  return data
}

export const getSubmitButtonText = (submitting) => {
  return submitting ? '提交中...' : '提交报名'
}
