import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  PAGE_SIZE,
  getLevelOptions,
  getTimeOptions,
  calculateDateRange,
  getDefaultFilters,
  getDefaultPagination,
  buildQueryParams,
  formatRewardAmount,
  formatRewardLevel,
  formatRewardRate,
  shouldFinishLoading,
  mergeRewardsList
} from '../../src/views/distributor/distributorRewards.logic.js'

describe('distributorRewards.logic', () => {
  describe('PAGE_SIZE', () => {
    it('is 20', () => {
      expect(PAGE_SIZE).toBe(20)
    })
  })

  describe('getLevelOptions', () => {
    it('returns 4 options', () => {
      const options = getLevelOptions()
      expect(options).toHaveLength(4)
    })

    it('first option is "全部级别" with value 0', () => {
      const options = getLevelOptions()
      expect(options[0]).toEqual({ text: '全部级别', value: 0 })
    })

    it('has level options 1-3', () => {
      const options = getLevelOptions()
      expect(options[1].value).toBe(1)
      expect(options[2].value).toBe(2)
      expect(options[3].value).toBe(3)
    })
  })

  describe('getTimeOptions', () => {
    it('returns 4 options', () => {
      const options = getTimeOptions()
      expect(options).toHaveLength(4)
    })

    it('has expected time values', () => {
      const options = getTimeOptions()
      const values = options.map(o => o.value)
      expect(values).toContain('all')
      expect(values).toContain('month')
      expect(values).toContain('week')
      expect(values).toContain('today')
    })
  })

  describe('calculateDateRange', () => {
    beforeEach(() => {
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2024-06-15T12:00:00Z'))
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('returns null for "all"', () => {
      expect(calculateDateRange('all')).toBeNull()
    })

    it('returns today range for "today"', () => {
      const range = calculateDateRange('today')
      expect(range.startDate).toBe('2024-06-15')
      expect(range.endDate).toBe('2024-06-15')
    })

    it('returns week range for "week"', () => {
      const range = calculateDateRange('week')
      expect(range.startDate).toBe('2024-06-08')
      expect(range.endDate).toBe('2024-06-15')
    })

    it('returns month range for "month"', () => {
      const range = calculateDateRange('month')
      expect(range).not.toBeNull()
      expect(range.startDate).toMatch(/^\d{4}-\d{2}-\d{2}$/)
      expect(range.endDate).toMatch(/^\d{4}-\d{2}-\d{2}$/)
    })
  })

  describe('getDefaultFilters', () => {
    it('returns default filter values', () => {
      const filters = getDefaultFilters()
      expect(filters).toEqual({
        level: 0,
        time: 'all'
      })
    })
  })

  describe('getDefaultPagination', () => {
    it('returns pagination defaults', () => {
      const pagination = getDefaultPagination()
      expect(pagination.page).toBe(1)
      expect(pagination.pageSize).toBe(PAGE_SIZE)
      expect(pagination.finished).toBe(false)
      expect(pagination.refreshing).toBe(false)
    })
  })

  describe('buildQueryParams', () => {
    it('builds basic params', () => {
      const params = buildQueryParams(1, 20, 0, 'all')
      expect(params).toEqual({ page: 1, pageSize: 20 })
    })

    it('includes level when > 0', () => {
      const params = buildQueryParams(1, 20, 2, 'all')
      expect(params.level).toBe(2)
    })

    it('does not include level when 0', () => {
      const params = buildQueryParams(1, 20, 0, 'all')
      expect(params.level).toBeUndefined()
    })

    it('includes date range for non-all time filter', () => {
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2024-06-15T12:00:00Z'))
      
      const params = buildQueryParams(1, 20, 0, 'today')
      expect(params.startDate).toBe('2024-06-15')
      expect(params.endDate).toBe('2024-06-15')
      
      vi.useRealTimers()
    })
  })

  describe('formatRewardAmount', () => {
    it('formats amount to 2 decimal places', () => {
      expect(formatRewardAmount(123.456)).toBe('123.46')
    })

    it('handles zero', () => {
      expect(formatRewardAmount(0)).toBe('0.00')
    })

    it('handles null', () => {
      expect(formatRewardAmount(null)).toBe('0.00')
    })
  })

  describe('formatRewardLevel', () => {
    it('formats level with suffix', () => {
      expect(formatRewardLevel(1)).toBe('1级奖励')
      expect(formatRewardLevel(2)).toBe('2级奖励')
    })
  })

  describe('formatRewardRate', () => {
    it('formats rate with percent sign', () => {
      expect(formatRewardRate(5.5)).toBe('5.5%')
    })
  })

  describe('shouldFinishLoading', () => {
    it('returns true when items < pageSize', () => {
      expect(shouldFinishLoading(15, 20)).toBe(true)
    })

    it('returns false when items >= pageSize', () => {
      expect(shouldFinishLoading(20, 20)).toBe(false)
      expect(shouldFinishLoading(25, 20)).toBe(false)
    })
  })

  describe('mergeRewardsList', () => {
    it('replaces list when refreshing', () => {
      const existing = [{ id: 1 }, { id: 2 }]
      const newItems = [{ id: 3 }, { id: 4 }]
      const result = mergeRewardsList(existing, newItems, true)
      expect(result).toEqual([{ id: 3 }, { id: 4 }])
    })

    it('appends to list when not refreshing', () => {
      const existing = [{ id: 1 }, { id: 2 }]
      const newItems = [{ id: 3 }, { id: 4 }]
      const result = mergeRewardsList(existing, newItems, false)
      expect(result).toEqual([{ id: 1 }, { id: 2 }, { id: 3 }, { id: 4 }])
    })
  })
})
