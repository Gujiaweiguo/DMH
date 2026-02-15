/**
 * DistributorWithdrawals business logic
 */

export const PAY_TYPES = [
  { label: '微信', value: 'wechat' },
  { label: '支付宝', value: 'alipay' },
  { label: '银行卡', value: 'bank' }
]

export const STATUS_TYPES = {
  pending: 'warning',
  approved: 'primary',
  rejected: 'danger',
  processing: 'success',
  completed: 'success',
  failed: 'danger'
}

export const STATUS_TEXTS = {
  pending: '待审核',
  approved: '已批准',
  rejected: '已拒绝',
  processing: '处理中',
  completed: '已完成',
  failed: '失败'
}

export const PAY_TYPE_TEXTS = {
  wechat: '微信',
  alipay: '支付宝',
  bank: '银行卡'
}

export const PAY_ACCOUNT_LABELS = {
  wechat: '微信号',
  alipay: '支付宝账号',
  bank: '银行卡号'
}

export const PAY_REAL_NAME_LABELS = {
  wechat: '微信昵称',
  alipay: '支付宝姓名',
  bank: '银行卡持卡人'
}

export const getStatusType = (status) => STATUS_TYPES[status] || 'default'

export const getStatusText = (status) => STATUS_TEXTS[status] || status

export const getPayTypeText = (payType) => PAY_TYPE_TEXTS[payType] || payType

export const getPayAccountLabel = (payType) => PAY_ACCOUNT_LABELS[payType] || '提现账号'

export const getPayRealNameLabel = (payType) => PAY_REAL_NAME_LABELS[payType] || '真实姓名'

export const getDefaultWithdrawalForm = () => ({
  amount: '',
  payType: 'wechat',
  payAccount: '',
  payRealName: ''
})

export const validateAmount = (value, balance) => {
  const amount = parseFloat(value)
  if (isNaN(amount) || amount <= 0) {
    return '提现金额必须大于0'
  }
  if (amount > balance) {
    return '提现金额不能超过可用余额'
  }
  return true
}

export const buildWithdrawalPayload = (form) => ({
  amount: parseFloat(form.amount),
  payType: form.payType,
  payAccount: form.payAccount,
  payRealName: form.payRealName
})

export const formatBalance = (balance) => {
  const num = Number(balance) || 0
  return num.toFixed(2)
}

export const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const year = d.getFullYear()
  const month = (d.getMonth() + 1).toString().padStart(2, '0')
  const day = d.getDate().toString().padStart(2, '0')
  const hours = d.getHours().toString().padStart(2, '0')
  const minutes = d.getMinutes().toString().padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

export const calculateNewBalance = (currentBalance, withdrawalAmount) => {
  return Math.max(0, currentBalance - withdrawalAmount)
}
