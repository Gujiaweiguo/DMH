import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  PHONE_REGEX,
  isValidPhone,
  STATUS_TEXT_MAP,
  statusText,
  formatDate,
  formatAmount,
  buildCampaignMap,
  getCampaignName,
  filterOrdersByPhone,
  savePhoneToStorage,
  loadPhoneFromStorage,
  hasFormData
} from '../../src/views/myOrders.logic.js'

describe('myOrders.logic', () => {
  describe('PHONE_REGEX', () => {
    it('matches valid phone numbers', () => {
      expect(PHONE_REGEX.test('13812345678')).toBe(true)
      expect(PHONE_REGEX.test('15912345678')).toBe(true)
      expect(PHONE_REGEX.test('18812345678')).toBe(true)
    })

    it('rejects invalid phone numbers', () => {
      expect(PHONE_REGEX.test('12345678901')).toBe(false)
      expect(PHONE_REGEX.test('1381234567')).toBe(false)
      expect(PHONE_REGEX.test('138123456789')).toBe(false)
    })
  })

  describe('isValidPhone', () => {
    it('returns true for valid phone', () => {
      expect(isValidPhone('13812345678')).toBe(true)
    })

    it('returns false for invalid phone', () => {
      expect(isValidPhone('12345')).toBe(false)
    })
  })

  describe('STATUS_TEXT_MAP', () => {
    it('has correct mappings', () => {
      expect(STATUS_TEXT_MAP.paid).toBe('已支付')
      expect(STATUS_TEXT_MAP.pending).toBe('待支付')
      expect(STATUS_TEXT_MAP.cancelled).toBe('已取消')
    })
  })

  describe('statusText', () => {
    it('returns correct text for paid', () => {
      expect(statusText('paid')).toBe('已支付')
    })

    it('returns original status for unknown', () => {
      expect(statusText('unknown')).toBe('unknown')
    })
  })

  describe('formatDate', () => {
    it('formats date string', () => {
      const result = formatDate('2024-06-15T10:30:00')
      expect(result).toContain('2024')
    })

    it('returns empty for null', () => {
      expect(formatDate(null)).toBe('')
    })
  })

  describe('formatAmount', () => {
    it('formats to 2 decimals', () => {
      expect(formatAmount(99.5)).toBe('99.50')
    })

    it('handles zero', () => {
      expect(formatAmount(0)).toBe('0.00')
    })

    it('handles null', () => {
      expect(formatAmount(null)).toBe('0.00')
    })
  })

  describe('buildCampaignMap', () => {
    it('builds map from campaigns', () => {
      const campaigns = [
        { id: 1, name: 'Campaign 1' },
        { id: 2, name: 'Campaign 2' }
      ]
      const map = buildCampaignMap(campaigns)
      expect(map[1].name).toBe('Campaign 1')
      expect(map[2].name).toBe('Campaign 2')
    })

    it('returns empty object for null', () => {
      expect(buildCampaignMap(null)).toEqual({})
    })
  })

  describe('getCampaignName', () => {
    it('returns campaign name from map', () => {
      const map = { 1: { name: 'Test Campaign' } }
      expect(getCampaignName(1, map)).toBe('Test Campaign')
    })

    it('returns 未知活动 for missing campaign', () => {
      expect(getCampaignName(999, {})).toBe('未知活动')
    })
  })

  describe('filterOrdersByPhone', () => {
    const orders = [
      { id: 1, phone: '13812345678' },
      { id: 2, phone: '13987654321' },
      { id: 3, phone: '13812345678' }
    ]

    it('filters orders by phone', () => {
      const result = filterOrdersByPhone(orders, '13812345678')
      expect(result).toHaveLength(2)
    })

    it('returns empty for no matches', () => {
      expect(filterOrdersByPhone(orders, '11111111111')).toEqual([])
    })

    it('handles null orders', () => {
      expect(filterOrdersByPhone(null, '13812345678')).toEqual([])
    })
  })

  describe('savePhoneToStorage', () => {
    beforeEach(() => {
      vi.stubGlobal('localStorage', {
        setItem: vi.fn(() => true)
      })
    })

    it('saves phone to storage', () => {
      expect(savePhoneToStorage('13812345678')).toBe(true)
    })
  })

  describe('loadPhoneFromStorage', () => {
    it('loads phone from storage', () => {
      vi.stubGlobal('localStorage', {
        getItem: vi.fn(() => '13812345678')
      })
      expect(loadPhoneFromStorage()).toBe('13812345678')
    })

    it('returns empty string when not set', () => {
      vi.stubGlobal('localStorage', {
        getItem: vi.fn(() => null)
      })
      expect(loadPhoneFromStorage()).toBe('')
    })
  })

  describe('hasFormData', () => {
    it('returns true for order with formData', () => {
      expect(hasFormData({ formData: { name: 'test' } })).toBe(true)
    })

    it('returns false for order without formData', () => {
      expect(hasFormData({ formData: null })).toBe(false)
      expect(hasFormData({})).toBe(false)
    })

    it('returns false for empty formData', () => {
      expect(hasFormData({ formData: {} })).toBe(false)
    })
  })
})
