/**
 * DistributorSubordinates business logic
 */

export const PAGE_SIZE = 20

export const getDefaultSubordinate = () => ({
  id: 0,
  username: '',
  level: 1,
  createdAt: '',
  totalOrders: 0,
  totalEarnings: 0
})

export const formatSubordinateName = (sub) => {
  const username = sub.username || '用户'
  return `${username}${sub.id}`
}

export const formatSubordinateInfo = (sub) => {
  return `${sub.level}级分销商 · 加入时间: ${sub.createdAt}`
}

export const formatSubordinateEarnings = (earnings) => {
  const num = Number(earnings) || 0
  return num.toFixed(2)
}

export const shouldFinishLoading = (newItemsCount, pageSize) => {
  return newItemsCount < pageSize
}

export const mergeSubordinatesList = (existingList, newList, isRefreshing) => {
  if (isRefreshing) {
    return [...newList]
  }
  return [...existingList, ...newList]
}

export const buildQueryParams = (page, pageSize) => ({
  page,
  pageSize
})
