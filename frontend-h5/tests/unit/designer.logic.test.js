import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  formatTime,
  formatCountdown,
  padZero,
  isExpired,
  getImageUrl,
  truncateText,
  parseJsonSafe,
  buildQueryString,
  DEFAULT_COLORS,
  validateRequired,
  validatePhone,
  validateEmail
} from '../../src/components/designer/designer.logic.js'

describe('designer.logic', () => {
  describe('formatTime', () => {
    it('should return empty string for null', () => {
      expect(formatTime(null)).toBe('')
    })

    it('should format timestamp correctly', () => {
      const result = formatTime('2024-01-15T10:30:00Z')
      expect(result).toContain('2024')
    })
  })

  describe('formatCountdown', () => {
    it('should return zeros for negative diff', () => {
      expect(formatCountdown(-1)).toEqual({ days: 0, hours: 0, minutes: 0, seconds: 0 })
    })

    it('should format diff correctly', () => {
      const diff = 2 * 24 * 60 * 60 * 1000 + 3 * 60 * 60 * 1000 + 30 * 60 * 1000 + 45 * 1000
      const result = formatCountdown(diff)
      expect(result.days).toBe(2)
      expect(result.hours).toBe(3)
      expect(result.minutes).toBe(30)
      expect(result.seconds).toBe(45)
    })
  })

  describe('padZero', () => {
    it('should pad single digit', () => {
      expect(padZero(5)).toBe('05')
    })

    it('should not pad double digit', () => {
      expect(padZero(15)).toBe('15')
    })
  })

  describe('isExpired', () => {
    beforeEach(() => {
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2024-06-15T12:00:00Z'))
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('should return true for past time', () => {
      expect(isExpired('2024-01-01T00:00:00Z')).toBe(true)
    })

    it('should return false for future time', () => {
      expect(isExpired('2024-12-31T23:59:59Z')).toBe(false)
    })

    it('should return false for null', () => {
      expect(isExpired(null)).toBe(false)
    })
  })

  describe('getImageUrl', () => {
    it('should return url if provided', () => {
      expect(getImageUrl('https://example.com/image.png')).toBe('https://example.com/image.png')
    })

    it('should return empty string for null', () => {
      expect(getImageUrl(null)).toBe('')
    })
  })

  describe('truncateText', () => {
    it('should truncate long text', () => {
      expect(truncateText('1234567890', 5)).toBe('12345...')
    })

    it('should not truncate short text', () => {
      expect(truncateText('123', 5)).toBe('123')
    })

    it('should return empty for null', () => {
      expect(truncateText(null, 5)).toBe('')
    })
  })

  describe('parseJsonSafe', () => {
    it('should parse valid JSON', () => {
      expect(parseJsonSafe('{"a":1}')).toEqual({ a: 1 })
    })

    it('should return default for invalid JSON', () => {
      expect(parseJsonSafe('invalid')).toEqual({})
    })

    it('should return default for null', () => {
      expect(parseJsonSafe(null)).toEqual({})
    })
  })

  describe('buildQueryString', () => {
    it('should build query string', () => {
      expect(buildQueryString({ a: 1, b: 'test' })).toBe('a=1&b=test')
    })

    it('should skip null values', () => {
      expect(buildQueryString({ a: 1, b: null })).toBe('a=1')
    })

    it('should skip empty strings', () => {
      expect(buildQueryString({ a: 1, b: '' })).toBe('a=1')
    })
  })

  describe('DEFAULT_COLORS', () => {
    it('should have primary color', () => {
      expect(DEFAULT_COLORS.primary).toBe('#4f46e5')
    })
  })

  describe('validateRequired', () => {
    it('should return error for empty value', () => {
      const result = validateRequired('', '名称')
      expect(result.valid).toBe(false)
      expect(result.error).toBe('名称不能为空')
    })

    it('should return valid for non-empty value', () => {
      const result = validateRequired('test', '名称')
      expect(result.valid).toBe(true)
    })
  })

  describe('validatePhone', () => {
    it('should return error for empty phone', () => {
      expect(validatePhone('').valid).toBe(false)
    })

    it('should return error for invalid phone', () => {
      expect(validatePhone('12345678901').valid).toBe(false)
    })

    it('should return valid for correct phone', () => {
      expect(validatePhone('13800138000').valid).toBe(true)
    })
  })

  describe('validateEmail', () => {
    it('should return valid for empty email', () => {
      expect(validateEmail('').valid).toBe(true)
    })

    it('should return error for invalid email', () => {
      expect(validateEmail('invalid').valid).toBe(false)
    })

    it('should return valid for correct email', () => {
      expect(validateEmail('test@example.com').valid).toBe(true)
    })
  })
})
