export const getDefaultBrandInfo = () => ({
  name: '示例品牌',
  logo: 'https://api.dicebear.com/7.x/initials/svg?seed=Brand',
})

export const getDefaultTodayStats = () => ({
  orders: 0,
  rewards: 0,
  promoters: 0,
  campaigns: 0,
})

export const getDashboardStatusText = (status) => {
  const statusMap = {
    active: '进行中',
    paused: '已暂停',
    ended: '已结束',
  }
  return statusMap[status] || status
}

export const getMockDashboardData = () => ({
  todayStats: {
    orders: 23,
    rewards: 1580,
    promoters: 8,
    campaigns: 3,
  },
  recentCampaigns: [
    {
      id: 1,
      name: '春节特惠活动',
      description: '新春佳节，推荐好友享双重奖励',
      status: 'active',
      orders: 156,
      rewards: 3120,
    },
    {
      id: 2,
      name: '会员招募计划',
      description: '招募品牌会员，享受专属优惠',
      status: 'active',
      orders: 89,
      rewards: 1780,
    },
  ],
})
