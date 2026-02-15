export const parseCurrentBrandId = (storageValue, userInfoRaw) => {
  const fromStorage = Number(storageValue)
  if (Number.isFinite(fromStorage) && fromStorage > 0) return fromStorage

  try {
    const info = JSON.parse(userInfoRaw || '{}')
    const firstBrandId = Array.isArray(info.brandIds) && info.brandIds.length > 0 ? Number(info.brandIds[0]) : 0
    if (Number.isFinite(firstBrandId) && firstBrandId > 0) {
      return firstBrandId
    }
  } catch {
    return 0
  }
  return 0
}

export const formatDistributorTime = (timeString) => {
  if (!timeString) return '-'
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

export const getDistributorStatusText = (status) => {
  const map = {
    active: '启用',
    suspended: '停用',
  }
  return map[status] || status || '-'
}

export const countDistributorsByStatus = (distributors, status) => {
  return distributors.filter((distributor) => distributor.status === status).length
}

export const buildDistributorListQuery = ({ page, pageSize, brandId, keyword, status, level }) => {
  const query = new URLSearchParams()
  query.set('page', String(page))
  query.set('pageSize', String(pageSize))
  query.set('brandId', String(brandId))
  if (keyword?.trim()) query.set('keyword', keyword.trim())
  if (status) query.set('status', status)
  if (level > 0) query.set('level', String(level))
  return query
}

export const mergeDistributorList = ({ currentList, incomingList, replace, total }) => {
  const distributors = replace ? incomingList : currentList.concat(incomingList)
  const hasMore = distributors.length < total
  return {
    distributors,
    hasMore,
  }
}
