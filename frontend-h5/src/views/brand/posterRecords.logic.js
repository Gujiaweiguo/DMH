/**
 * PosterRecords business logic
 */

export const TYPE_TABS = [
  { value: 'all', label: '全部' },
  { value: 'campaign', label: '活动海报' },
  { value: 'distributor', label: '分销商海报' }
]

export const RECORD_TYPE_LABELS = {
  campaign: '活动海报',
  distributor: '分销商海报'
}

export const formatDateTime = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export const getRecordTypeLabel = (type) => RECORD_TYPE_LABELS[type] || type

export const filterByType = (records, type) => {
  if (type === 'all') return records
  return records.filter(r => r.recordType === type)
}

export const filterByKeyword = (records, keyword) => {
  if (!keyword) return records
  const lowerKeyword = keyword.toLowerCase()
  return records.filter(r =>
    r.campaignName?.toLowerCase().includes(lowerKeyword) ||
    r.distributorName?.toLowerCase().includes(lowerKeyword)
  )
}

export const filterRecords = (records, type, keyword) => {
  let result = filterByType(records, type)
  result = filterByKeyword(result, keyword)
  return result
}

export const calculateStats = (records) => {
  const today = new Date().toLocaleDateString('zh-CN')
  
  return {
    total: records.length,
    campaign: records.filter(r => r.recordType === 'campaign').length,
    distributor: records.filter(r => r.recordType === 'distributor').length,
    today: records.filter(r => {
      const recordDate = new Date(r.createdAt).toLocaleDateString('zh-CN')
      return recordDate === today
    }).length,
    totalDownloads: records.reduce((sum, r) => sum + (r.downloadCount || 0), 0)
  }
}

export const buildPosterFileName = (campaignName) => {
  return `海报_${campaignName}_${Date.now()}.png`
}

export const getDefaultDateRange = () => ({
  start: '',
  end: ''
})

export const getDefaultStats = () => ({
  total: 0,
  campaign: 0,
  distributor: 0,
  today: 0,
  totalDownloads: 0
})

export const getRegenerateButtonText = (regenerating) => {
  return regenerating ? '生成中...' : '重新生成'
}
