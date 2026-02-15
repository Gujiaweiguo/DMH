import { describe, it, expect } from 'vitest'
import {
  PAGE_SIZE,
  getDefaultSubordinate,
  formatSubordinateName,
  formatSubordinateInfo,
  formatSubordinateEarnings,
  shouldFinishLoading,
  mergeSubordinatesList,
  buildQueryParams
} from '../../src/views/distributor/distributorSubordinates.logic.js'

describe('distributorSubordinates.logic', () => {
  describe('PAGE_SIZE', () => {
    it('is 20', () => {
      expect(PAGE_SIZE).toBe(20)
    })
  })

  describe('getDefaultSubordinate', () => {
    it('returns default subordinate object', () => {
      const sub = getDefaultSubordinate()
      expect(sub).toEqual({
        id: 0,
        username: '',
        level: 1,
        createdAt: '',
        totalOrders: 0,
        totalEarnings: 0
      })
    })
  })

  describe('formatSubordinateName', () => {
    it('formats name with id', () => {
      const sub = { id: 123, username: '张三' }
      expect(formatSubordinateName(sub)).toBe('张三123')
    })

    it('uses fallback for missing username', () => {
      const sub = { id: 456, username: null }
      expect(formatSubordinateName(sub)).toBe('用户456')
    })

    it('handles empty username', () => {
      const sub = { id: 789, username: '' }
      expect(formatSubordinateName(sub)).toBe('用户789')
    })
  })

  describe('formatSubordinateInfo', () => {
    it('formats level and date', () => {
      const sub = { level: 2, createdAt: '2024-06-15' }
      expect(formatSubordinateInfo(sub)).toBe('2级分销商 · 加入时间: 2024-06-15')
    })
  })

  describe('formatSubordinateEarnings', () => {
    it('formats earnings to 2 decimals', () => {
      expect(formatSubordinateEarnings(123.456)).toBe('123.46')
    })

    it('handles zero', () => {
      expect(formatSubordinateEarnings(0)).toBe('0.00')
    })

    it('handles null', () => {
      expect(formatSubordinateEarnings(null)).toBe('0.00')
    })

    it('handles string number', () => {
      expect(formatSubordinateEarnings('99.9')).toBe('99.90')
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

  describe('mergeSubordinatesList', () => {
    it('replaces list when refreshing', () => {
      const existing = [{ id: 1 }, { id: 2 }]
      const newList = [{ id: 3 }, { id: 4 }]
      const result = mergeSubordinatesList(existing, newList, true)
      expect(result).toEqual([{ id: 3 }, { id: 4 }])
    })

    it('appends to list when not refreshing', () => {
      const existing = [{ id: 1 }, { id: 2 }]
      const newList = [{ id: 3 }, { id: 4 }]
      const result = mergeSubordinatesList(existing, newList, false)
      expect(result).toEqual([{ id: 1 }, { id: 2 }, { id: 3 }, { id: 4 }])
    })
  })

  describe('buildQueryParams', () => {
    it('builds query params', () => {
      expect(buildQueryParams(1, 20)).toEqual({ page: 1, pageSize: 20 })
    })
  })
})
