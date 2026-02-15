/**
 * MemberDetail business logic
 */

export const GENDER_MAP = {
  0: '未知',
  1: '男',
  2: '女'
}

export const getGenderText = (gender) => GENDER_MAP[gender] || '未知'

export const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

export const formatAmount = (amount) => {
  const num = Number(amount) || 0
  return num.toFixed(2)
}

export const getStatusText = (status) => {
  return status === 'active' ? '正常' : '禁用'
}

export const isActiveStatus = (status) => status === 'active'

export const getAvatarUrl = (avatar) => avatar || '/default-avatar.png'

export const getNickname = (nickname) => nickname || '未设置昵称'

export const getPhone = (phone) => phone || '未绑定'

export const hasTags = (member) => !!(member && member.tags && member.tags.length > 0)

export const hasBrands = (member) => !!(member && member.brands && member.brands.length > 0)
