import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  getStatusType,
  getStatusText,
  getDefaultStatistics,
  clearDistributorAuth,
  buildPromotionRoute,
  buildRewardsRoute,
  buildSubordinatesRoute,
  buildStatisticsRoute,
  buildWithdrawalsRoute,
  buildPosterRoute,
  formatMoney
} from '../../src/views/distributor/distributorCenter.logic.js'

describe('distributorCenter.logic', () => {
  describe('getStatusType', () => {
    it('returns "warning" for pending', () => {
      expect(getStatusType('pending')).toBe('warning')
    })

    it('returns "success" for approved', () => {
      expect(getStatusType('approved')).toBe('success')
    })

    it('returns "danger" for rejected', () => {
      expect(getStatusType('rejected')).toBe('danger')
    })

    it('returns "default" for unknown status', () => {
      expect(getStatusType('unknown')).toBe('default')
    })
  })

  describe('getStatusText', () => {
    it('returns "待审核" for pending', () => {
      expect(getStatusText('pending')).toBe('待审核')
    })

    it('returns "已通过" for approved', () => {
      expect(getStatusText('approved')).toBe('已通过')
    })

    it('returns "已拒绝" for rejected', () => {
      expect(getStatusText('rejected')).toBe('已拒绝')
    })

    it('returns original status for unknown', () => {
      expect(getStatusText('unknown')).toBe('unknown')
    })
  })

  describe('getDefaultStatistics', () => {
    it('returns default statistics object', () => {
      const stats = getDefaultStatistics()
      expect(stats).toEqual({
        totalEarnings: 0,
        totalOrders: 0,
        subordinatesCount: 0,
        todayEarnings: 0,
        monthEarnings: 0,
        balance: 0
      })
    })
  })

  describe('clearDistributorAuth', () => {
    beforeEach(() => {
      vi.stubGlobal('localStorage', {
        removeItem: vi.fn()
      })
    })

    it('removes all auth items from localStorage', () => {
      clearDistributorAuth()
      expect(localStorage.removeItem).toHaveBeenCalledWith('dmh_token')
      expect(localStorage.removeItem).toHaveBeenCalledWith('dmh_user_role')
      expect(localStorage.removeItem).toHaveBeenCalledWith('dmh_user_info')
      expect(localStorage.removeItem).toHaveBeenCalledWith('dmh_current_brand_id')
    })
  })

  describe('route builders', () => {
    it('buildPromotionRoute returns correct path', () => {
      expect(buildPromotionRoute(123)).toBe('/distributor/promotion?brandId=123')
    })

    it('buildRewardsRoute returns correct path', () => {
      expect(buildRewardsRoute(456)).toBe('/distributor/rewards?brandId=456')
    })

    it('buildSubordinatesRoute returns correct path', () => {
      expect(buildSubordinatesRoute(789)).toBe('/distributor/subordinates?brandId=789')
    })

    it('buildStatisticsRoute returns correct path', () => {
      expect(buildStatisticsRoute(100)).toBe('/distributor/statistics?brandId=100')
    })

    it('buildWithdrawalsRoute returns correct path', () => {
      expect(buildWithdrawalsRoute(200)).toBe('/distributor/withdrawals?brandId=200')
    })

    it('buildPosterRoute returns correct path', () => {
      expect(buildPosterRoute(5)).toBe('/poster-generator/5')
    })
  })

  describe('formatMoney', () => {
    it('formats number to 2 decimal places', () => {
      expect(formatMoney(123.456)).toBe('123.46')
    })

    it('handles zero', () => {
      expect(formatMoney(0)).toBe('0.00')
    })

    it('handles null', () => {
      expect(formatMoney(null)).toBe('0.00')
    })

    it('handles undefined', () => {
      expect(formatMoney(undefined)).toBe('0.00')
    })

    it('handles string number', () => {
      expect(formatMoney('99.9')).toBe('99.90')
    })
  })
})
