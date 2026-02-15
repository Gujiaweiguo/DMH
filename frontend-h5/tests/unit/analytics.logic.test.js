import { describe, it, expect } from 'vitest'
import {
  PERIOD_OPTIONS,
  getDefaultCoreMetrics,
  getDefaultChartData,
  getDefaultFunnelData,
  EXPORT_TYPES,
  getExportTypeLabel,
  formatMetricValue,
  calculateChangeIndicator,
  getRankingNumberStyle
} from '../../src/views/brand/analytics.logic.js'

describe('analytics.logic', () => {
  describe('PERIOD_OPTIONS', () => {
    it('has 4 options', () => {
      expect(PERIOD_OPTIONS).toHaveLength(4)
    })

    it('includes today, week, month, quarter', () => {
      const values = PERIOD_OPTIONS.map(o => o.value)
      expect(values).toContain('today')
      expect(values).toContain('week')
      expect(values).toContain('month')
      expect(values).toContain('quarter')
    })
  })

  describe('getDefaultCoreMetrics', () => {
    it('returns default metrics object', () => {
      const metrics = getDefaultCoreMetrics()
      expect(metrics).toEqual({
        totalRevenue: 0,
        totalOrders: 0,
        activePromoters: 0,
        avgOrderValue: 0
      })
    })
  })

  describe('getDefaultChartData', () => {
    it('returns object with empty orders array', () => {
      const chartData = getDefaultChartData()
      expect(chartData.orders).toEqual([])
    })
  })

  describe('getDefaultFunnelData', () => {
    it('returns 5 funnel steps', () => {
      const funnel = getDefaultFunnelData()
      expect(funnel).toHaveLength(5)
    })

    it('first step is 页面访问 with 100%', () => {
      const funnel = getDefaultFunnelData()
      expect(funnel[0].label).toBe('页面访问')
      expect(funnel[0].percentage).toBe(100)
    })
  })

  describe('EXPORT_TYPES', () => {
    it('has 4 export types', () => {
      expect(Object.keys(EXPORT_TYPES)).toHaveLength(4)
    })

    it('has orders, promoters, campaigns, all', () => {
      expect(EXPORT_TYPES.orders).toBe('订单数据')
      expect(EXPORT_TYPES.promoters).toBe('推广员数据')
      expect(EXPORT_TYPES.campaigns).toBe('活动数据')
      expect(EXPORT_TYPES.all).toBe('完整报表')
    })
  })

  describe('getExportTypeLabel', () => {
    it('returns label for known type', () => {
      expect(getExportTypeLabel('orders')).toBe('订单数据')
    })

    it('returns type for unknown', () => {
      expect(getExportTypeLabel('unknown')).toBe('unknown')
    })
  })

  describe('formatMetricValue', () => {
    it('formats number with locale string', () => {
      expect(formatMetricValue(12345)).toBe('12,345')
    })

    it('adds prefix when provided', () => {
      expect(formatMetricValue(100, '¥')).toBe('¥100')
    })

    it('handles zero', () => {
      expect(formatMetricValue(0)).toBe('0')
    })

    it('handles null', () => {
      expect(formatMetricValue(null)).toBe('0')
    })
  })

  describe('calculateChangeIndicator', () => {
    it('calculates positive change', () => {
      const result = calculateChangeIndicator(120, 100)
      expect(result.type).toBe('positive')
      expect(result.value).toBe('20.0')
    })

    it('calculates negative change', () => {
      const result = calculateChangeIndicator(80, 100)
      expect(result.type).toBe('negative')
      expect(result.value).toBe('20.0')
    })

    it('returns neutral for zero previous', () => {
      const result = calculateChangeIndicator(100, 0)
      expect(result.type).toBe('neutral')
    })
  })

  describe('getRankingNumberStyle', () => {
    it('returns gold for index 0', () => {
      expect(getRankingNumberStyle(0)).toBe('gold')
    })

    it('returns silver for index 1', () => {
      expect(getRankingNumberStyle(1)).toBe('silver')
    })

    it('returns bronze for index 2', () => {
      expect(getRankingNumberStyle(2)).toBe('bronze')
    })

    it('returns default for other indices', () => {
      expect(getRankingNumberStyle(3)).toBe('default')
    })
  })
})
