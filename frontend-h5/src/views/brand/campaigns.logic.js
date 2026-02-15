export const filterCampaignsByStatus = (campaigns, status) => {
  if (status === 'all') {
    return campaigns
  }
  return campaigns.filter((campaign) => campaign.status === status)
}

export const findCampaignStatusLabel = (statusTabs, currentStatus) => {
  const status = statusTabs.find((item) => item.value === currentStatus)
  return status ? status.label : ''
}

export const getCampaignStatusText = (status) => {
  const statusMap = {
    active: '进行中',
    paused: '已暂停',
    ended: '已结束',
  }
  return statusMap[status] || status
}

export const getCampaignStatusTagType = (status) => {
  const typeMap = {
    active: 'success',
    paused: 'warning',
    ended: 'danger',
  }
  return typeMap[status] || 'default'
}

export const formatCampaignDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
  })
}

export const getFallbackCampaigns = () => [
  {
    id: 1,
    name: '春节特惠活动',
    description: '新春佳节，推荐好友享双重奖励',
    status: 'active',
    rewardRule: 88,
    startTime: '2026-02-01 00:00:00',
    endTime: '2026-02-15 23:59:59',
    participantCount: 156,
  },
  {
    id: 2,
    name: '会员招募计划',
    description: '招募品牌会员，享受专属优惠',
    status: 'active',
    rewardRule: 66,
    startTime: '2026-01-01 00:00:00',
    endTime: '2026-12-31 23:59:59',
    participantCount: 89,
  },
  {
    id: 3,
    name: '元宵节活动',
    description: '元宵佳节，猜灯谜赢大奖',
    status: 'paused',
    rewardRule: 50,
    startTime: '2026-02-28 00:00:00',
    endTime: '2026-03-01 23:59:59',
    participantCount: 23,
  },
]

export const getNextCampaignStatus = (currentStatus) => {
  return currentStatus === 'active' ? 'paused' : 'active'
}

export const getCampaignStatusActionText = (nextStatus) => {
  return nextStatus === 'active' ? '启用' : '暂停'
}
