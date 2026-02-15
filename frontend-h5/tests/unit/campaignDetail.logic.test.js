import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  formatTime,
  formatReward,
  formatPaymentAmount,
  getPaymentTypeText,
  getPaymentLabel,
  isPaymentRequired,
  calculateConversionRate,
  isBrandAdmin,
  buildSourceData,
  saveSourceToStorage,
  buildPosterDownloadName,
  buildVerifyRoute,
  buildPaymentRoute,
  buildFormRoute
} from '../../src/views/campaignDetail.logic.js'

describe('campaignDetail.logic', () => {
  describe('formatTime', () => {
    it('formats ISO time string', () => {
      expect(formatTime('2024-06-15T10:30:00Z')).toBe('2024-06-15 10:30')
    })

    it('returns empty for null', () => {
      expect(formatTime(null)).toBe('')
    })

    it('returns empty for undefined', () => {
      expect(formatTime(undefined)).toBe('')
    })
  })

  describe('formatReward', () => {
    it('formats to 2 decimals', () => {
      expect(formatReward(88.5)).toBe('88.50')
    })

    it('handles zero', () => {
      expect(formatReward(0)).toBe('0.00')
    })

    it('handles null', () => {
      expect(formatReward(null)).toBe('0.00')
    })
  })

  describe('formatPaymentAmount', () => {
    it('formats amount to 2 decimals', () => {
      expect(formatPaymentAmount(99.9)).toBe('99.90')
    })
  })

  describe('getPaymentTypeText', () => {
    it('returns 订金支付 for deposit', () => {
      expect(getPaymentTypeText('deposit')).toBe('订金支付')
    })

    it('returns 全款支付 for full', () => {
      expect(getPaymentTypeText('full')).toBe('全款支付')
    })
  })

  describe('getPaymentLabel', () => {
    it('returns 订金金额 for deposit', () => {
      expect(getPaymentLabel('deposit')).toBe('订金金额')
    })

    it('returns 支付金额 for full', () => {
      expect(getPaymentLabel('full')).toBe('支付金额')
    })
  })

  describe('isPaymentRequired', () => {
    it('returns true when requirePayment is true', () => {
      expect(isPaymentRequired({ requirePayment: true })).toBe(true)
    })

    it('returns false when requirePayment is false', () => {
      expect(isPaymentRequired({ requirePayment: false })).toBe(false)
    })

    it('returns false for null config', () => {
      expect(isPaymentRequired(null)).toBe(false)
    })
  })

  describe('calculateConversionRate', () => {
    it('calculates conversion rate', () => {
      expect(calculateConversionRate(100, 25)).toBe('25.0')
    })

    it('returns 0.0 for zero participants', () => {
      expect(calculateConversionRate(0, 10)).toBe('0.0')
    })

    it('handles null participant count', () => {
      expect(calculateConversionRate(null, 10)).toBe('0.0')
    })
  })

  describe('isBrandAdmin', () => {
    it('returns true for brand_admin role', () => {
      vi.stubGlobal('localStorage', {
        getItem: vi.fn(() => 'brand_admin')
      })
      expect(isBrandAdmin()).toBe(true)
    })

    it('returns false for other roles', () => {
      vi.stubGlobal('localStorage', {
        getItem: vi.fn(() => 'distributor')
      })
      expect(isBrandAdmin()).toBe(false)
    })
  })

  describe('buildSourceData', () => {
    it('builds source data from query', () => {
      const query = { c_id: 'camp1', u_id: 'user1' }
      expect(buildSourceData(query)).toEqual({ c_id: 'camp1', u_id: 'user1' })
    })

    it('defaults to empty strings', () => {
      expect(buildSourceData({})).toEqual({ c_id: '', u_id: '' })
    })
  })

  describe('saveSourceToStorage', () => {
    beforeEach(() => {
      vi.stubGlobal('localStorage', {
        setItem: vi.fn(() => true)
      })
    })

    it('saves source to storage', () => {
      const result = saveSourceToStorage({ c_id: '123', u_id: '456' })
      expect(result).toBe(true)
    })
  })

  describe('buildPosterDownloadName', () => {
    it('builds poster filename', () => {
      expect(buildPosterDownloadName(123)).toBe('poster_123.png')
    })
  })

  describe('buildVerifyRoute', () => {
    it('builds verify route with campaignId', () => {
      const route = buildVerifyRoute(456)
      expect(route.path).toBe('/verify')
      expect(route.query.campaignId).toBe(456)
    })
  })

  describe('buildPaymentRoute', () => {
    it('builds payment route', () => {
      expect(buildPaymentRoute(789)).toBe('/campaign/789/payment')
    })
  })

  describe('buildFormRoute', () => {
    it('builds verify route for brand admin', () => {
      const route = buildFormRoute(123, true)
      expect(route.path).toBe('/verify')
    })

    it('builds form route for non-admin', () => {
      expect(buildFormRoute(123, false)).toBe('/campaign/123/form')
    })
  })
})
