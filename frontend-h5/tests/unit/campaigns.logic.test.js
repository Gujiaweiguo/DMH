import { describe, expect, it } from 'vitest'
import {
  filterCampaignsByStatus,
  findCampaignStatusLabel,
  formatCampaignDate,
  getCampaignStatusActionText,
  getCampaignStatusTagType,
  getCampaignStatusText,
  getFallbackCampaigns,
  getNextCampaignStatus,
} from '../../src/views/brand/campaigns.logic.js'

describe('campaigns logic', () => {
  const campaigns = [
    { id: 1, status: 'active' },
    { id: 2, status: 'paused' },
    { id: 3, status: 'active' },
  ]

  it('filters campaigns by status', () => {
    expect(filterCampaignsByStatus(campaigns, 'all')).toEqual(campaigns)
    expect(filterCampaignsByStatus(campaigns, 'active').map((item) => item.id)).toEqual([1, 3])
    expect(filterCampaignsByStatus(campaigns, 'paused').map((item) => item.id)).toEqual([2])
  })

  it('finds current status label', () => {
    const tabs = [
      { value: 'all', label: '全部' },
      { value: 'active', label: '进行中' },
    ]
    expect(findCampaignStatusLabel(tabs, 'active')).toBe('进行中')
    expect(findCampaignStatusLabel(tabs, 'missing')).toBe('')
  })

  it('maps campaign status text and tag type', () => {
    expect(getCampaignStatusText('active')).toBe('进行中')
    expect(getCampaignStatusText('paused')).toBe('已暂停')
    expect(getCampaignStatusText('ended')).toBe('已结束')
    expect(getCampaignStatusText('custom')).toBe('custom')

    expect(getCampaignStatusTagType('active')).toBe('success')
    expect(getCampaignStatusTagType('paused')).toBe('warning')
    expect(getCampaignStatusTagType('ended')).toBe('danger')
    expect(getCampaignStatusTagType('custom')).toBe('default')
  })

  it('formats campaign date', () => {
    expect(formatCampaignDate('2026-02-13 10:00:00')).toContain('2')
  })

  it('returns fallback campaigns', () => {
    const fallback = getFallbackCampaigns()
    expect(fallback).toHaveLength(3)
    expect(fallback[0]).toMatchObject({ name: '春节特惠活动', status: 'active' })
  })

  it('calculates next status and action text', () => {
    expect(getNextCampaignStatus('active')).toBe('paused')
    expect(getNextCampaignStatus('paused')).toBe('active')
    expect(getCampaignStatusActionText('active')).toBe('启用')
    expect(getCampaignStatusActionText('paused')).toBe('暂停')
  })
})
