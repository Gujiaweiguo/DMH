export const getBrandIdFromStorage = (storedBrandId, userInfoRaw) => {
  if (storedBrandId) {
    const parsed = Number.parseInt(storedBrandId, 10)
    return Number.isNaN(parsed) ? null : parsed
  }

  if (!userInfoRaw) return null

  try {
    const userInfo = JSON.parse(userInfoRaw)
    if (Array.isArray(userInfo.brandIds) && userInfo.brandIds.length > 0) {
      return userInfo.brandIds[0]
    }
  } catch (_error) {
    return null
  }

  return null
}

export const buildMembersQueryParams = ({ page, pageSize, brandId, keyword, source, status }) => {
  const params = new URLSearchParams()
  params.append('page', String(page))
  params.append('pageSize', String(pageSize))

  if (brandId) params.append('brandId', String(brandId))
  if (keyword) params.append('keyword', keyword)
  if (source) params.append('source', source)
  if (status) params.append('status', status)

  return params
}

export const mergeMembersPage = ({ currentMembers, incomingMembers, page, pageSize }) => {
  const list = incomingMembers || []
  const members = page === 1 ? list : [...currentMembers, ...list]
  const finished = list.length < pageSize
  const nextPage = finished ? page : page + 1

  return {
    members,
    finished,
    nextPage,
  }
}

export const buildExportRequestBody = ({ brandId, reason, keyword, source, status }) => ({
  brandId,
  reason,
  filters: JSON.stringify({
    keyword,
    source,
    status,
  }),
})

export const formatMemberDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(
    date.getDate(),
  ).padStart(2, '0')}`
}
