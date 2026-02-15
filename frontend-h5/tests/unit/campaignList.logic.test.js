import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  STATUS_TEXT_MAP,
  EMPTY_TEXT_MAP,
  statusText,
  emptyText,
  formatDate,
  formatReward,
  buildTabs,
  filterCampaigns,
  buildSourceData,
  saveSourceToStorage,
  loadMyPhone,
  markRegisteredCampaigns
} from '../../src/views/campaignList.logic.js'

describe('campaignList.logic', () => {
  describe('STATUS_TEXT_MAP', () => {
    it('has correct mappings', () => {
      expect(STATUS_TEXT_MAP.active).toBe('进行中')
      expect(STATUS_TEXT_MAP.paused).toBe('已暂停')
      expect(STATUS_TEXT_MAP.ended).toBe('已结束')
    })
  })

  describe('EMPTY_TEXT_MAP', () => {
    it('has correct mappings', () => {
      expect(EMPTY_TEXT_MAP.all).toBe('活动')
      expect(EMPTY_TEXT_MAP.ongoing).toBe('进行中的活动')
      expect(EMPTY_TEXT_MAP.ended).toBe('已结束的活动')
    })
  })

  describe('statusText', () => {
    it('returns correct text for active', () => {
      expect(statusText('active')).toBe('进行中')
    })

    it('returns original status for unknown', () => {
      expect(statusText('unknown')).toBe('unknown')
    })
  })

  describe('emptyText', () => {
    it('returns correct text for tab', () => {
      expect(emptyText('all')).toBe('活动')
      expect(emptyText('ongoing')).toBe('进行中的活动')
    })
  })

  describe('formatDate', () => {
    it('extracts date part', () => {
      expect(formatDate('2024-06-15 10:30:00')).toBe('2024-06-15')
    })

    it('returns empty for null', () => {
      expect(formatDate(null)).toBe('')
    })

    it('returns empty for undefined', () => {
      expect(formatDate(undefined)).toBe('')
    })
  })

  describe('formatReward', () => {
    it('formats to 2 decimal places', () => {
      expect(formatReward(88.5)).toBe('88.50')
    })

    it('handles zero', () => {
      expect(formatReward(0)).toBe('0.00')
    })

    it('handles null', () => {
      expect(formatReward(null)).toBe('0.00')
    })
  })

  describe('buildTabs', () => {
    it('builds tabs with counts', () => {
      const campaigns = [
        { id: 1, status: 'active' },
        { id: 2, status: 'active' },
        { id: 3, status: 'ended' }
      ]
      const tabs = buildTabs(campaigns)
      expect(tabs[0].count).toBe(3)
      expect(tabs[1].count).toBe(2)
      expect(tabs[2].count).toBe(1)
    })

    it('handles empty array', () => {
      const tabs = buildTabs([])
      expect(tabs[0].count).toBe(0)
    })
  })

  describe('filterCampaigns', () => {
    const campaigns = [
      { id: 1, status: 'active', isRegistered: false },
      { id: 2, status: 'active', isRegistered: true },
      { id: 3, status: 'ended', isRegistered: false }
    ]

    it('returns all for "all" tab', () => {
      expect(filterCampaigns(campaigns, 'all')).toHaveLength(3)
    })

    it('filters active for "ongoing" tab', () => {
      expect(filterCampaigns(campaigns, 'ongoing')).toHaveLength(2)
    })

    it('filters unregistered when onlyUnregistered is true', () => {
      const result = filterCampaigns(campaigns, 'ongoing', true)
      expect(result).toHaveLength(1)
      expect(result[0].id).toBe(1)
    })

    it('filters ended for "ended" tab', () => {
      expect(filterCampaigns(campaigns, 'ended')).toHaveLength(1)
    })

    it('returns all for unknown tab', () => {
      expect(filterCampaigns(campaigns, 'unknown')).toHaveLength(3)
    })
  })

  describe('buildSourceData', () => {
    it('builds source from query', () => {
      const query = { c_id: 'camp123', u_id: 'user456' }
      expect(buildSourceData(query)).toEqual({
        c_id: 'camp123',
        u_id: 'user456'
      })
    })

    it('defaults to empty strings', () => {
      expect(buildSourceData({})).toEqual({ c_id: '', u_id: '' })
    })
  })

  describe('saveSourceToStorage', () => {
    let localStorageMock

    beforeEach(() => {
      localStorageMock = {
        setItem: vi.fn(() => true)
      }
      vi.stubGlobal('localStorage', localStorageMock)
    })

    it('saves source to localStorage', () => {
      const result = saveSourceToStorage({ c_id: '123', u_id: '456' })
      expect(result).toBe(true)
    })

    it('returns false on error', () => {
      localStorageMock.setItem.mockImplementation(() => { throw new Error('QuotaExceededError') })
      const result = saveSourceToStorage({ c_id: '123', u_id: '456' })
      expect(result).toBe(false)
    })
  })

  describe('loadMyPhone', () => {
    let localStorageMock

    beforeEach(() => {
      localStorageMock = {
        getItem: vi.fn(() => '13812345678')
      }
      vi.stubGlobal('localStorage', localStorageMock)
    })

    it('returns phone from storage', () => {
      expect(loadMyPhone()).toBe('13812345678')
    })

    it('returns empty string when not set', () => {
      localStorageMock.getItem.mockReturnValue(null)
      expect(loadMyPhone()).toBe('')
    })

    it('returns empty string on error', () => {
      localStorageMock.getItem.mockImplementation(() => { throw new Error('Storage error') })
      expect(loadMyPhone()).toBe('')
    })
  })

  describe('markRegisteredCampaigns', () => {
    it('marks campaigns as registered', () => {
      const campaigns = [
        { id: 1, name: 'Campaign 1' },
        { id: 2, name: 'Campaign 2' }
      ]
      const orders = [
        { campaignId: 1 },
        { campaignId: 3 }
      ]
      const result = markRegisteredCampaigns(campaigns, orders)
      expect(result[0].isRegistered).toBe(true)
      expect(result[1].isRegistered).toBe(false)
    })
  })
})
