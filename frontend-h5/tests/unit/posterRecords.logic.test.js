import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  TYPE_TABS,
  RECORD_TYPE_LABELS,
  formatDateTime,
  getRecordTypeLabel,
  filterByType,
  filterByKeyword,
  filterRecords,
  calculateStats,
  buildPosterFileName,
  getDefaultDateRange,
  getDefaultStats,
  getRegenerateButtonText
} from '../../src/views/brand/posterRecords.logic.js'

describe('posterRecords.logic', () => {
  describe('TYPE_TABS', () => {
    it('should have correct tab values', () => {
      expect(TYPE_TABS).toHaveLength(3)
      expect(TYPE_TABS[0].value).toBe('all')
      expect(TYPE_TABS[1].value).toBe('campaign')
      expect(TYPE_TABS[2].value).toBe('distributor')
    })
  })

  describe('RECORD_TYPE_LABELS', () => {
    it('should have correct type labels', () => {
      expect(RECORD_TYPE_LABELS['campaign']).toBe('活动海报')
      expect(RECORD_TYPE_LABELS['distributor']).toBe('分销商海报')
    })
  })

  describe('formatDateTime', () => {
    it('should return empty string for null', () => {
      expect(formatDateTime(null)).toBe('')
    })

    it('should format date correctly', () => {
      const result = formatDateTime('2024-01-15T10:30:00Z')
      expect(result).toContain('2024')
    })
  })

  describe('getRecordTypeLabel', () => {
    it('should return correct label for known type', () => {
      expect(getRecordTypeLabel('campaign')).toBe('活动海报')
      expect(getRecordTypeLabel('distributor')).toBe('分销商海报')
    })

    it('should return type for unknown type', () => {
      expect(getRecordTypeLabel('unknown')).toBe('unknown')
    })
  })

  describe('filterByType', () => {
    const records = [
      { id: 1, recordType: 'campaign' },
      { id: 2, recordType: 'distributor' },
      { id: 3, recordType: 'campaign' }
    ]

    it('should return all records for all type', () => {
      expect(filterByType(records, 'all')).toHaveLength(3)
    })

    it('should filter by campaign type', () => {
      const result = filterByType(records, 'campaign')
      expect(result).toHaveLength(2)
      expect(result.every(r => r.recordType === 'campaign')).toBe(true)
    })

    it('should filter by distributor type', () => {
      const result = filterByType(records, 'distributor')
      expect(result).toHaveLength(1)
    })
  })

  describe('filterByKeyword', () => {
    const records = [
      { id: 1, campaignName: '春节活动', distributorName: '张三' },
      { id: 2, campaignName: '夏季促销', distributorName: '李四' }
    ]

    it('should return all records for empty keyword', () => {
      expect(filterByKeyword(records, '')).toHaveLength(2)
    })

    it('should filter by campaign name', () => {
      const result = filterByKeyword(records, '春节')
      expect(result).toHaveLength(1)
      expect(result[0].campaignName).toBe('春节活动')
    })

    it('should filter by distributor name', () => {
      const result = filterByKeyword(records, '李四')
      expect(result).toHaveLength(1)
      expect(result[0].distributorName).toBe('李四')
    })

    it('should be case insensitive', () => {
      const result = filterByKeyword(records, 'ZHANG')
      expect(result).toHaveLength(0)
    })
  })

  describe('filterRecords', () => {
    const records = [
      { id: 1, recordType: 'campaign', campaignName: '春节活动' },
      { id: 2, recordType: 'distributor', campaignName: '夏季促销' }
    ]

    it('should combine type and keyword filters', () => {
      const result = filterRecords(records, 'campaign', '春节')
      expect(result).toHaveLength(1)
      expect(result[0].id).toBe(1)
    })

    it('should return empty when no match', () => {
      const result = filterRecords(records, 'campaign', '夏季')
      expect(result).toHaveLength(0)
    })
  })

  describe('calculateStats', () => {
    beforeEach(() => {
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2024-01-15T12:00:00Z'))
    })

    it('should calculate correct stats', () => {
      const records = [
        { id: 1, recordType: 'campaign', createdAt: '2024-01-15T10:00:00Z', downloadCount: 5 },
        { id: 2, recordType: 'distributor', createdAt: '2024-01-14T10:00:00Z', downloadCount: 3 },
        { id: 3, recordType: 'campaign', createdAt: '2024-01-15T11:00:00Z', downloadCount: 2 }
      ]

      const stats = calculateStats(records)
      expect(stats.total).toBe(3)
      expect(stats.campaign).toBe(2)
      expect(stats.distributor).toBe(1)
      expect(stats.totalDownloads).toBe(10)
    })
  })

  describe('buildPosterFileName', () => {
    it('should build filename with campaign name and timestamp', () => {
      const result = buildPosterFileName('春节活动')
      expect(result).toContain('春节活动')
      expect(result).toMatch(/\.png$/)
    })
  })

  describe('getDefaultDateRange', () => {
    it('should return empty date range', () => {
      const result = getDefaultDateRange()
      expect(result).toEqual({ start: '', end: '' })
    })
  })

  describe('getDefaultStats', () => {
    it('should return zero stats', () => {
      const result = getDefaultStats()
      expect(result).toEqual({
        total: 0,
        campaign: 0,
        distributor: 0,
        today: 0,
        totalDownloads: 0
      })
    })
  })

  describe('getRegenerateButtonText', () => {
    it('should return generating text when regenerating', () => {
      expect(getRegenerateButtonText(true)).toBe('生成中...')
    })

    it('should return normal text when not regenerating', () => {
      expect(getRegenerateButtonText(false)).toBe('重新生成')
    })
  })
})
