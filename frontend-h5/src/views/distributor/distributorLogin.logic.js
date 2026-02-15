/**
 * DistributorLogin business logic
 */

export const getDefaultForm = () => ({
  username: '',
  password: ''
})

export const hasDistributorRole = (data) => {
  return !!(data && data.roles && data.roles.includes('distributor'))
}

export const saveDistributorAuth = (data) => {
  if (!data) return
  localStorage.setItem('dmh_token', data.token)
  localStorage.setItem('dmh_user_role', 'distributor')
  localStorage.setItem('dmh_user_info', JSON.stringify(data))
}

export const isAlreadyLoggedIn = () => {
  const token = localStorage.getItem('dmh_token')
  const userRole = localStorage.getItem('dmh_user_role')
  return !!(token && userRole === 'distributor')
}

export const getLoginErrorMessage = (error) => {
  return error?.message || '登录失败，请重试'
}

export const validateLoginForm = (form) => {
  const errors = []
  
  if (!form?.username || form.username.trim() === '') {
    errors.push('请输入用户名')
  }
  
  if (!form?.password || form.password.trim() === '') {
    errors.push('请输入密码')
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

export const getButtonText = (loading) => {
  return loading ? '登录中...' : '登录'
}
