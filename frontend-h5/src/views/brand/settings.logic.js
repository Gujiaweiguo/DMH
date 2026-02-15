export const getDefaultBrandInfo = () => ({
  name: '示例品牌',
  description: '专业的数字营销服务品牌',
  logo: 'https://api.dicebear.com/7.x/initials/svg?seed=Brand',
  phone: '400-123-4567',
  email: 'contact@brand.com',
})

export const getDefaultRewardSettings = () => ({
  defaultRate: 20,
  minWithdraw: 100,
  settlementType: 'instant',
})

export const getDefaultNotificationSettings = () => ({
  newOrder: true,
  newPromoter: true,
  dailyReport: false,
  email: 'admin@brand.com',
})

export const getDefaultSyncSettings = () => ({
  status: 'connected',
  frequency: 'realtime',
  dataTypes: ['orders', 'users'],
})

export const getDefaultPasswordForm = () => ({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

export const getSyncStatusText = (status) => {
  const statusMap = {
    connected: '已连接',
    disconnected: '未连接',
    error: '连接错误',
  }
  return statusMap[status] || status
}

export const validatePasswordForm = (passwordForm) => {
  if (!passwordForm.oldPassword || !passwordForm.newPassword) {
    return '请填写完整的密码信息'
  }

  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    return '两次输入的新密码不一致'
  }

  return ''
}
