import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  formatDate,
  formatDateTime,
  formatTime,
  formatDateTimeShort,
  isToday,
  isYesterday,
  isTomorrow,
  getRelativeDate,
  getDaysDiff,
  addDays,
  getStartOfDay,
  getEndOfDay
} from '../../src/utils/dateFormat.logic.js'

describe('dateFormat.logic', () => {
  describe('formatDate', () => {
    it('should return empty string for null', () => {
      expect(formatDate(null)).toBe('')
    })

    it('should return empty string for invalid date', () => {
      expect(formatDate('invalid')).toBe('')
    })

    it('should format date with default format', () => {
      const result = formatDate('2024-01-15T10:30:00Z')
      expect(result).toMatch(/2024-01-15/)
    })

    it('should format with custom format', () => {
      const result = formatDate('2024-01-15T10:30:45Z', 'YYYY-MM-DD HH:mm:ss')
      expect(result).toMatch(/2024-01-15/)
    })
  })

  describe('formatDateTime', () => {
    it('should format with date and time', () => {
      const result = formatDateTime('2024-01-15')
      expect(result).toContain('2024-01-15')
    })
  })

  describe('formatTime', () => {
    it('should format time only', () => {
      const result = formatTime('2024-01-15T10:30:45')
      expect(result).toContain(':')
    })
  })

  describe('formatDateTimeShort', () => {
    it('should format with short format', () => {
      const result = formatDateTimeShort('2024-01-15T10:30:00')
      expect(result).toContain('-')
    })
  })

  describe('isToday', () => {
    beforeEach(() => {
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2024-06-15T12:00:00Z'))
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('should return true for today', () => {
      expect(isToday('2024-06-15T10:00:00Z')).toBe(true)
    })

    it('should return false for other day', () => {
      expect(isToday('2024-06-14T10:00:00Z')).toBe(false)
    })

    it('should return false for null', () => {
      expect(isToday(null)).toBe(false)
    })
  })

  describe('isYesterday', () => {
    beforeEach(() => {
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2024-06-15T12:00:00Z'))
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('should return true for yesterday', () => {
      expect(isYesterday('2024-06-14T10:00:00Z')).toBe(true)
    })

    it('should return false for today', () => {
      expect(isYesterday('2024-06-15T10:00:00Z')).toBe(false)
    })
  })

  describe('isTomorrow', () => {
    beforeEach(() => {
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2024-06-15T12:00:00Z'))
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('should return true for tomorrow', () => {
      expect(isTomorrow('2024-06-16T10:00:00Z')).toBe(true)
    })

    it('should return false for today', () => {
      expect(isTomorrow('2024-06-15T10:00:00Z')).toBe(false)
    })
  })

  describe('getRelativeDate', () => {
    beforeEach(() => {
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2024-06-15T12:00:00Z'))
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('should return 今天 for today', () => {
      expect(getRelativeDate('2024-06-15')).toBe('今天')
    })

    it('should return 昨天 for yesterday', () => {
      expect(getRelativeDate('2024-06-14')).toBe('昨天')
    })

    it('should return 明天 for tomorrow', () => {
      expect(getRelativeDate('2024-06-16')).toBe('明天')
    })

    it('should return empty for null', () => {
      expect(getRelativeDate(null)).toBe('')
    })
  })

  describe('getDaysDiff', () => {
    it('should return positive diff', () => {
      expect(getDaysDiff('2024-06-16', '2024-06-14')).toBe(2)
    })

    it('should return negative diff', () => {
      expect(getDaysDiff('2024-06-14', '2024-06-16')).toBe(-2)
    })

    it('should return 0 for same day', () => {
      expect(getDaysDiff('2024-06-15', '2024-06-15')).toBe(0)
    })
  })

  describe('addDays', () => {
    it('should add positive days', () => {
      const result = addDays('2024-06-15', 3)
      expect(result.getDate()).toBe(18)
    })

    it('should add negative days', () => {
      const result = addDays('2024-06-15', -3)
      expect(result.getDate()).toBe(12)
    })
  })

  describe('getStartOfDay', () => {
    it('should return start of day', () => {
      const result = getStartOfDay('2024-06-15T15:30:00')
      expect(result.getHours()).toBe(0)
      expect(result.getMinutes()).toBe(0)
    })
  })

  describe('getEndOfDay', () => {
    it('should return end of day', () => {
      const result = getEndOfDay('2024-06-15T15:30:00')
      expect(result.getHours()).toBe(23)
      expect(result.getMinutes()).toBe(59)
    })
  })
})
