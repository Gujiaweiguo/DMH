import { describe, expect, it } from 'vitest'
import {
  buildPromoterLink,
  buildPromoterLinkForm,
  calculatePromoterStats,
  filterAndSortPromoters,
  formatPromoterTime,
  getPromoterStatusText,
} from '../../src/views/brand/promoters.logic.js'

describe('promoters logic', () => {
  const promoters = [
    {
      id: 1,
      name: '张推广',
      phone: '138****1234',
      status: 'active',
      level: 'VIP',
      totalRewards: 1200,
      conversionRate: 16,
      joinDate: '2026-02-12',
      todayOrders: 3,
    },
    {
      id: 2,
      name: '李推广',
      phone: '139****5678',
      status: 'inactive',
      level: '银牌',
      totalRewards: 400,
      conversionRate: 8,
      joinDate: '2025-12-01',
      todayOrders: 0,
    },
    {
      id: 3,
      name: '王推广',
      phone: '137****9999',
      status: 'active',
      level: '金牌',
      totalRewards: 1600,
      conversionRate: 12,
      joinDate: '2026-02-09',
      todayOrders: 1,
    },
  ]

  it('maps promoter status text', () => {
    expect(getPromoterStatusText('active')).toBe('活跃')
    expect(getPromoterStatusText('inactive')).toBe('不活跃')
    expect(getPromoterStatusText('blocked')).toBe('已封禁')
    expect(getPromoterStatusText('other')).toBe('other')
  })

  it('formats promoter activity time', () => {
    expect(formatPromoterTime('2026-02-13 10:00:00')).toContain('2')
  })

  it('filters by active and sorts by rewards desc', () => {
    const result = filterAndSortPromoters(promoters, 'active', '')
    expect(result.map((item) => item.id)).toEqual([3, 1])
  })

  it('filters top promoters and supports keyword search', () => {
    const top = filterAndSortPromoters(promoters, 'top', '')
    expect(top.map((item) => item.id)).toEqual([3, 1])

    const byKeyword = filterAndSortPromoters(promoters, 'all', '139****')
    expect(byKeyword).toHaveLength(1)
    expect(byKeyword[0].id).toBe(2)
  })

  it('filters new promoters by one week window', () => {
    const now = new Date('2026-02-13T00:00:00Z')
    const result = filterAndSortPromoters(promoters, 'new', '', now)
    expect(result.map((item) => item.id)).toEqual([3, 1])
  })

  it('calculates promoter stats with safe zero handling', () => {
    const stats = calculatePromoterStats(promoters)
    expect(stats).toEqual({
      active: 2,
      totalRewards: 3200,
      todayOrders: 4,
      conversionRate: 12,
    })

    const emptyStats = calculatePromoterStats([])
    expect(emptyStats.conversionRate).toBe(0)
  })

  it('builds link form and promo url', () => {
    const form = buildPromoterLinkForm({ id: 9, name: '赵推广' })
    expect(form).toEqual({ promoterId: 9, promoterName: '赵推广', campaignId: '' })

    expect(buildPromoterLink('https://dmh.test', 11, 9)).toBe('https://dmh.test/campaign/11?ref=9')
    expect(buildPromoterLink('https://dmh.test', '', 9)).toBe('')
  })
})
