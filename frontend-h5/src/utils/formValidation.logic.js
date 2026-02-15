/**
 * Form validation utilities
 */

export const REQUIRED_RULE = { required: true, message: '此字段必填' }

export const createMinLengthRule = (min) => ({
  min,
  message: `最少${min}个字符`
})

export const createMaxLengthRule = (max) => ({
  max,
  message: `最多${max}个字符`
})

export const createRangeRule = (min, max) => ({
  min,
  max,
  message: `取值范围: ${min} - ${max}`
})

export const PHONE_RULE = {
  pattern: /^1[3-9]\d{9}$/,
  message: '请输入正确的手机号'
}

export const EMAIL_RULE = {
  pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  message: '请输入正确的邮箱格式'
}

export const validateField = (value, rules) => {
  if (!rules || rules.length === 0) return { valid: true, error: null }
  
  for (const rule of rules) {
    if (rule.required && (!value || (typeof value === 'string' && !value.trim()))) {
      return { valid: false, error: rule.message }
    }
    
    if (rule.min !== undefined && typeof value === 'string' && value.length < rule.min) {
      return { valid: false, error: rule.message }
    }
    
    if (rule.max !== undefined && typeof value === 'string' && value.length > rule.max) {
      return { valid: false, error: rule.message }
    }
    
    if (rule.pattern && !rule.pattern.test(value)) {
      return { valid: false, error: rule.message }
    }
  }
  
  return { valid: true, error: null }
}

export const validateForm = (formData, schema) => {
  const errors = {}
  let isValid = true
  
  for (const [field, rules] of Object.entries(schema)) {
    const result = validateField(formData[field], rules)
    if (!result.valid) {
      errors[field] = result.error
      isValid = false
    }
  }
  
  return { valid: isValid, errors }
}

export const getFirstError = (errors) => {
  const values = Object.values(errors)
  return values.length > 0 ? values[0] : null
}

export const hasErrors = (errors) => {
  return Object.keys(errors).length > 0
}
