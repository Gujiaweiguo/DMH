/**
 * CampaignList business logic
 */

export const STATUS_TEXT_MAP = {
  active: '进行中',
  paused: '已暂停',
  ended: '已结束'
}

export const EMPTY_TEXT_MAP = {
  all: '活动',
  ongoing: '进行中的活动',
  ended: '已结束的活动'
}

export const statusText = (status) => STATUS_TEXT_MAP[status] || status

export const emptyText = (tab) => EMPTY_TEXT_MAP[tab] || '活动'

export const formatDate = (time) => {
  if (!time) return ''
  return time.substring(0, 10)
}

export const formatReward = (reward) => {
  const num = Number(reward) || 0
  return num.toFixed(2)
}

export const buildTabs = (campaigns) => {
  const ongoing = campaigns.filter(c => c.status === 'active')
  const ended = campaigns.filter(c => c.status === 'ended' || c.status === 'paused')
  
  return [
    { key: 'all', label: '全部', count: campaigns.length },
    { key: 'ongoing', label: '进行中', count: ongoing.length },
    { key: 'ended', label: '已结束', count: ended.length }
  ]
}

export const filterCampaigns = (campaigns, activeTab, onlyUnregistered = false) => {
  if (activeTab === 'all') {
    return campaigns
  }
  
  if (activeTab === 'ongoing') {
    return campaigns
      .filter(c => c.status === 'active')
      .filter(c => !onlyUnregistered || !c.isRegistered)
  }
  
  if (activeTab === 'ended') {
    return campaigns.filter(c => c.status === 'ended' || c.status === 'paused')
  }
  
  return campaigns
}

export const buildSourceData = (query) => ({
  c_id: query.c_id || '',
  u_id: query.u_id || ''
})

export const saveSourceToStorage = (source) => {
  try {
    localStorage.setItem('dmh_source', JSON.stringify(source))
    return true
  } catch (e) {
    return false
  }
}

export const loadMyPhone = () => {
  try {
    return localStorage.getItem('dmh_my_phone') || ''
  } catch (e) {
    return ''
  }
}

export const markRegisteredCampaigns = (campaigns, orders) => {
  return campaigns.map(campaign => {
    const isRegistered = orders.some(order => order.campaignId === campaign.id)
    return { ...campaign, isRegistered }
  })
}
