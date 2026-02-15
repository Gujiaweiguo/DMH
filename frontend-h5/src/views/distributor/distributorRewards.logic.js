/**
 * DistributorRewards business logic
 */

export const PAGE_SIZE = 20

export const getLevelOptions = () => [
  { text: '全部级别', value: 0 },
  { text: '一级奖励', value: 1 },
  { text: '二级奖励', value: 2 },
  { text: '三级奖励', value: 3 }
]

export const getTimeOptions = () => [
  { text: '全部时间', value: 'all' },
  { text: '本月', value: 'month' },
  { text: '本周', value: 'week' },
  { text: '今日', value: 'today' }
]

export const calculateDateRange = (filterTime) => {
  if (filterTime === 'all') {
    return null
  }

  const now = new Date()
  let startDate, endDate

  if (filterTime === 'today') {
    startDate = now.toISOString().split('T')[0]
    endDate = startDate
  } else if (filterTime === 'week') {
    const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
    startDate = weekAgo.toISOString().split('T')[0]
    endDate = now.toISOString().split('T')[0]
  } else if (filterTime === 'month') {
    const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)
    startDate = monthStart.toISOString().split('T')[0]
    endDate = now.toISOString().split('T')[0]
  }

  return startDate && endDate ? { startDate, endDate } : null
}

export const getDefaultFilters = () => ({
  level: 0,
  time: 'all'
})

export const getDefaultPagination = () => ({
  page: 1,
  pageSize: PAGE_SIZE,
  finished: false,
  refreshing: false
})

export const buildQueryParams = (page, pageSize, filterLevel, filterTime) => {
  const params = { page, pageSize }

  if (filterLevel > 0) {
    params.level = filterLevel
  }

  const dateRange = calculateDateRange(filterTime)
  if (dateRange) {
    params.startDate = dateRange.startDate
    params.endDate = dateRange.endDate
  }

  return params
}

export const formatRewardAmount = (amount) => {
  const num = Number(amount) || 0
  return num.toFixed(2)
}

export const formatRewardLevel = (level) => `${level}级奖励`

export const formatRewardRate = (rate) => `${rate}%`

export const shouldFinishLoading = (newItemsCount, pageSize) => {
  return newItemsCount < pageSize
}

export const mergeRewardsList = (existingRewards, newRewards, isRefreshing) => {
  if (isRefreshing) {
    return [...newRewards]
  }
  return [...existingRewards, ...newRewards]
}
