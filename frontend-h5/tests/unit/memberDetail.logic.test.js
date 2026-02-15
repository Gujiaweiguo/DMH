import { describe, it, expect } from 'vitest'
import {
  GENDER_MAP,
  getGenderText,
  formatDateTime,
  formatAmount,
  getStatusText,
  isActiveStatus,
  getAvatarUrl,
  getNickname,
  getPhone,
  hasTags,
  hasBrands
} from '../../src/views/brand/memberDetail.logic.js'

describe('memberDetail.logic', () => {
  describe('GENDER_MAP', () => {
    it('should have correct gender mappings', () => {
      expect(GENDER_MAP[0]).toBe('未知')
      expect(GENDER_MAP[1]).toBe('男')
      expect(GENDER_MAP[2]).toBe('女')
    })
  })

  describe('getGenderText', () => {
    it('should return correct text for valid gender', () => {
      expect(getGenderText(0)).toBe('未知')
      expect(getGenderText(1)).toBe('男')
      expect(getGenderText(2)).toBe('女')
    })

    it('should return unknown for invalid gender', () => {
      expect(getGenderText(99)).toBe('未知')
      expect(getGenderText(null)).toBe('未知')
    })
  })

  describe('formatDateTime', () => {
    it('should return dash for null', () => {
      expect(formatDateTime(null)).toBe('-')
    })

    it('should format date correctly', () => {
      const result = formatDateTime('2024-01-15T10:30:00Z')
      expect(result).toMatch(/2024/)
    })
  })

  describe('formatAmount', () => {
    it('should format number to 2 decimal places', () => {
      expect(formatAmount(100)).toBe('100.00')
      expect(formatAmount(99.9)).toBe('99.90')
    })

    it('should handle zero', () => {
      expect(formatAmount(0)).toBe('0.00')
    })

    it('should handle null', () => {
      expect(formatAmount(null)).toBe('0.00')
    })
  })

  describe('getStatusText', () => {
    it('should return correct text for active status', () => {
      expect(getStatusText('active')).toBe('正常')
    })

    it('should return disabled for other status', () => {
      expect(getStatusText('inactive')).toBe('禁用')
      expect(getStatusText(null)).toBe('禁用')
    })
  })

  describe('isActiveStatus', () => {
    it('should return true for active status', () => {
      expect(isActiveStatus('active')).toBe(true)
    })

    it('should return false for other status', () => {
      expect(isActiveStatus('inactive')).toBe(false)
      expect(isActiveStatus(null)).toBe(false)
    })
  })

  describe('getAvatarUrl', () => {
    it('should return avatar if provided', () => {
      expect(getAvatarUrl('https://example.com/avatar.png')).toBe('https://example.com/avatar.png')
    })

    it('should return default avatar if not provided', () => {
      expect(getAvatarUrl(null)).toBe('/default-avatar.png')
      expect(getAvatarUrl('')).toBe('/default-avatar.png')
    })
  })

  describe('getNickname', () => {
    it('should return nickname if provided', () => {
      expect(getNickname('张三')).toBe('张三')
    })

    it('should return default text if not provided', () => {
      expect(getNickname(null)).toBe('未设置昵称')
      expect(getNickname('')).toBe('未设置昵称')
    })
  })

  describe('getPhone', () => {
    it('should return phone if provided', () => {
      expect(getPhone('13800138000')).toBe('13800138000')
    })

    it('should return default text if not provided', () => {
      expect(getPhone(null)).toBe('未绑定')
      expect(getPhone('')).toBe('未绑定')
    })
  })

  describe('hasTags', () => {
    it('should return true for member with tags', () => {
      expect(hasTags({ tags: [{ id: 1 }] })).toBe(true)
    })

    it('should return false for member without tags', () => {
      expect(hasTags({ tags: [] })).toBe(false)
      expect(hasTags({})).toBe(false)
      expect(hasTags(null)).toBe(false)
    })
  })

  describe('hasBrands', () => {
    it('should return true for member with brands', () => {
      expect(hasBrands({ brands: [{ brandId: 1 }] })).toBe(true)
    })

    it('should return false for member without brands', () => {
      expect(hasBrands({ brands: [] })).toBe(false)
      expect(hasBrands({})).toBe(false)
      expect(hasBrands(null)).toBe(false)
    })
  })
})
