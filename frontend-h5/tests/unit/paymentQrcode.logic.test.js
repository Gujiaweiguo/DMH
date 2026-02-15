import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  formatExpireTime,
  formatAmount,
  isExpired,
  getDefaultQrcodeData,
  buildQrcodeData,
  validateCampaignId,
  getMinutesUntilExpiry
} from '../../src/views/paymentQrcode.logic.js'

describe('paymentQrcode.logic', () => {
  describe('formatExpireTime', () => {
    it('should return dash for null', () => {
      expect(formatExpireTime(null)).toBe('-')
    })

    it('should return dash for empty string', () => {
      expect(formatExpireTime('')).toBe('-')
    })

    it('should return expired for past time', () => {
      const pastTime = new Date(Date.now() - 1000 * 60 * 10).toISOString()
      expect(formatExpireTime(pastTime)).toBe('已过期')
    })

    it('should format minutes correctly', () => {
      const futureTime = new Date(Date.now() + 1000 * 60 * 5).toISOString()
      const result = formatExpireTime(futureTime)
      expect(result).toMatch(/\d+分钟后过期/)
    })

    it('should format hours and minutes correctly', () => {
      const futureTime = new Date(Date.now() + 1000 * 60 * 65).toISOString()
      const result = formatExpireTime(futureTime)
      expect(result).toMatch(/1小时\d+分钟后过期/)
    })

    it('should return original string for invalid date', () => {
      expect(formatExpireTime('invalid')).toBe('invalid')
    })

    it('should handle date parsing error', () => {
      const badDate = { toString: () => { throw new Error('bad date') } }
      const result = formatExpireTime(badDate)
      expect(result).toBe(badDate)
    })
  })

  describe('formatAmount', () => {
    it('should format number to 2 decimal places', () => {
      expect(formatAmount(100)).toBe('100.00')
      expect(formatAmount(99.9)).toBe('99.90')
      expect(formatAmount(99.999)).toBe('100.00')
    })

    it('should handle zero', () => {
      expect(formatAmount(0)).toBe('0.00')
    })

    it('should handle null', () => {
      expect(formatAmount(null)).toBe('0.00')
    })

    it('should handle string numbers', () => {
      expect(formatAmount('100.5')).toBe('100.50')
    })
  })

  describe('isExpired', () => {
    it('should return false for null', () => {
      expect(isExpired(null)).toBe(false)
    })

    it('should return true for past time', () => {
      const pastTime = new Date(Date.now() - 1000).toISOString()
      expect(isExpired(pastTime)).toBe(true)
    })

    it('should return false for future time', () => {
      const futureTime = new Date(Date.now() + 1000 * 60 * 10).toISOString()
      expect(isExpired(futureTime)).toBe(false)
    })

    it('should return false for invalid date', () => {
      expect(isExpired('invalid')).toBe(false)
    })

    it('should return false on date comparison error', () => {
      const badInput = { 
        toString: () => { throw new Error('bad') }
      }
      expect(isExpired(badInput)).toBe(false)
    })
  })

  describe('getDefaultQrcodeData', () => {
    it('should return default qrcode data structure', () => {
      const result = getDefaultQrcodeData()
      expect(result).toEqual({
        qrcodeBase64: null,
        campaignName: null,
        amount: null,
        expireAt: null
      })
    })
  })

  describe('buildQrcodeData', () => {
    it('should return null for null response', () => {
      expect(buildQrcodeData(null)).toBeNull()
    })

    it('should return null for response without data', () => {
      expect(buildQrcodeData({})).toBeNull()
    })

    it('should build qrcode data from response', () => {
      const response = {
        data: {
          qrcodeUrl: 'data:image/png;base64,abc',
          campaignName: '测试活动',
          amount: 100,
          expireAt: '2024-12-31T23:59:59Z'
        }
      }
      const result = buildQrcodeData(response)
      expect(result.qrcodeBase64).toBe('data:image/png;base64,abc')
      expect(result.campaignName).toBe('测试活动')
      expect(result.amount).toBe(100)
      expect(result.expireAt).toBe('2024-12-31T23:59:59Z')
    })

    it('should use default values for missing fields', () => {
      const response = { data: {} }
      const result = buildQrcodeData(response)
      expect(result.campaignName).toBe('未知活动')
      expect(result.amount).toBe(0)
      expect(result.expireAt).toBeDefined()
    })
  })

  describe('validateCampaignId', () => {
    it('should return error for null id', () => {
      const result = validateCampaignId(null)
      expect(result.valid).toBe(false)
      expect(result.error).toBe('活动ID无效')
    })

    it('should return error for empty id', () => {
      const result = validateCampaignId('')
      expect(result.valid).toBe(false)
      expect(result.error).toBe('活动ID无效')
    })

    it('should return valid for valid id', () => {
      const result = validateCampaignId('123')
      expect(result.valid).toBe(true)
      expect(result.error).toBeNull()
    })
  })

  describe('getMinutesUntilExpiry', () => {
    it('should return 0 for null', () => {
      expect(getMinutesUntilExpiry(null)).toBe(0)
    })

    it('should return 0 for past time', () => {
      const pastTime = new Date(Date.now() - 1000 * 60).toISOString()
      expect(getMinutesUntilExpiry(pastTime)).toBe(0)
    })

    it('should return minutes for future time', () => {
      const futureTime = new Date(Date.now() + 1000 * 60 * 10).toISOString()
      const result = getMinutesUntilExpiry(futureTime)
      expect(result).toBeGreaterThanOrEqual(9)
      expect(result).toBeLessThanOrEqual(10)
    })

    it('should return 0 for invalid date', () => {
      expect(getMinutesUntilExpiry('invalid')).toBe(0)
    })
  })
})
