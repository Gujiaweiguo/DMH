/**
 * Brand Login business logic
 */

export const TEST_USERNAME = 'brand_manager'
export const TEST_PASSWORD = '123456'

export const getDefaultForm = () => ({
  username: '',
  password: ''
})

export const quickFillTestAccount = (form) => {
  form.username = TEST_USERNAME
  form.password = TEST_PASSWORD
  return form
}

export const validateLoginForm = (form) => {
  const errors = []
  
  if (!form || !form.username || form.username.trim() === '') {
    errors.push('请输入用户名')
  }
  
  if (!form || !form.password || form.password.trim() === '') {
    errors.push('请输入密码')
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

export const hasBrandAccess = (data) => {
  return !!(data && data.brandIds && Array.isArray(data.brandIds) && data.brandIds.length > 0)
}

export const getFirstBrandId = (brandIds) => {
  if (Array.isArray(brandIds) && brandIds.length > 0) {
    return brandIds[0]
  }
  return null
}

export const saveLoginInfo = (data, firstBrandId) => {
  try {
    localStorage.setItem('dmh_token', data.token)
    localStorage.setItem('dmh_user_role', 'participant')
    localStorage.setItem('dmh_user_info', JSON.stringify(data))
    localStorage.setItem('dmh_current_brand_id', String(firstBrandId))
    return true
  } catch (e) {
    console.error('保存登录信息失败', e)
    return false
  }
}

export const clearLoginInfo = () => {
  try {
    localStorage.removeItem('dmh_token')
    localStorage.removeItem('dmh_user_role')
    localStorage.removeItem('dmh_user_info')
    localStorage.removeItem('dmh_current_brand_id')
    return true
  } catch (e) {
    return false
  }
}

export const getLoginButtonText = (loading) => {
  return loading ? '登录中...' : '登录'
}

export const buildLoginError = (error) => {
  if (!error) return '登录失败，请重试'
  if (typeof error === 'string') return error
  return error.message || '登录失败，请重试'
}
