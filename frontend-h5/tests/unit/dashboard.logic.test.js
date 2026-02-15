import { describe, expect, it } from 'vitest'
import {
  getDashboardStatusText,
  getDefaultBrandInfo,
  getDefaultTodayStats,
  getMockDashboardData,
} from '../../src/views/brand/dashboard.logic.js'

describe('dashboard logic', () => {
  it('provides default brand info and stats', () => {
    expect(getDefaultBrandInfo()).toMatchObject({
      name: '示例品牌',
    })
    expect(getDefaultTodayStats()).toEqual({
      orders: 0,
      rewards: 0,
      promoters: 0,
      campaigns: 0,
    })
  })

  it('maps dashboard campaign status text', () => {
    expect(getDashboardStatusText('active')).toBe('进行中')
    expect(getDashboardStatusText('paused')).toBe('已暂停')
    expect(getDashboardStatusText('ended')).toBe('已结束')
    expect(getDashboardStatusText('custom')).toBe('custom')
  })

  it('provides mock dashboard data', () => {
    const data = getMockDashboardData()
    expect(data.todayStats).toMatchObject({
      orders: 23,
      rewards: 1580,
    })
    expect(data.recentCampaigns).toHaveLength(2)
    expect(data.recentCampaigns[0]).toMatchObject({
      status: 'active',
    })
  })
})
