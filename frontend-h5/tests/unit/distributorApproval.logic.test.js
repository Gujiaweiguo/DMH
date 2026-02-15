import { describe, it, expect } from 'vitest'
import {
  formatTime,
  getStatusText,
  getDefaultApprovalForm,
  getDefaultPagination,
  buildTabs,
  validateRejection,
  buildApprovePayload,
  buildRejectPayload,
  getMockPendingApplications,
  getAvatarText
} from '../../src/views/brand/distributorApproval.logic.js'

describe('distributorApproval.logic', () => {
  describe('formatTime', () => {
    it('returns "-" for null', () => {
      expect(formatTime(null)).toBe('-')
    })

    it('returns "-" for undefined', () => {
      expect(formatTime(undefined)).toBe('-')
    })

    it('returns "-" for empty string', () => {
      expect(formatTime('')).toBe('-')
    })

    it('formats ISO date string', () => {
      const result = formatTime('2024-06-15T14:30:00.000Z')
      expect(result).toMatch(/\d{2}\/\d{2}/) // MM/DD format
    })
  })

  describe('getStatusText', () => {
    it('returns "待审批" for pending', () => {
      expect(getStatusText('pending')).toBe('待审批')
    })

    it('returns "已通过" for approved', () => {
      expect(getStatusText('approved')).toBe('已通过')
    })

    it('returns "已拒绝" for rejected', () => {
      expect(getStatusText('rejected')).toBe('已拒绝')
    })

    it('returns original status for unknown status', () => {
      expect(getStatusText('unknown')).toBe('unknown')
    })
  })

  describe('getDefaultApprovalForm', () => {
    it('returns form with level 1', () => {
      const form = getDefaultApprovalForm()
      expect(form.level).toBe(1)
    })

    it('returns form with empty notes', () => {
      const form = getDefaultApprovalForm()
      expect(form.notes).toBe('')
    })
  })

  describe('getDefaultPagination', () => {
    it('returns currentPage 1', () => {
      const pagination = getDefaultPagination()
      expect(pagination.currentPage).toBe(1)
    })

    it('returns pageSize 20', () => {
      const pagination = getDefaultPagination()
      expect(pagination.pageSize).toBe(20)
    })

    it('returns hasMore false', () => {
      const pagination = getDefaultPagination()
      expect(pagination.hasMore).toBe(false)
    })
  })

  describe('buildTabs', () => {
    it('builds tabs with counts', () => {
      const tabs = buildTabs(5, 10)
      expect(tabs).toHaveLength(2)
      expect(tabs[0]).toEqual({ key: 'pending', label: '待审批', count: 5 })
      expect(tabs[1]).toEqual({ key: 'processed', label: '已处理', count: 10 })
    })

    it('defaults counts to 0', () => {
      const tabs = buildTabs()
      expect(tabs[0].count).toBe(0)
      expect(tabs[1].count).toBe(0)
    })
  })

  describe('validateRejection', () => {
    it('returns false for null form', () => {
      expect(validateRejection(null)).toBe(false)
    })

    it('returns false for form without notes', () => {
      expect(validateRejection({ level: 1 })).toBe(false)
    })

    it('returns false for form with empty notes', () => {
      expect(validateRejection({ level: 1, notes: '' })).toBe(false)
    })

    it('returns false for form with whitespace-only notes', () => {
      expect(validateRejection({ level: 1, notes: '   ' })).toBe(false)
    })

    it('returns true for form with valid notes', () => {
      expect(validateRejection({ level: 1, notes: '不符合条件' })).toBe(true)
    })
  })

  describe('buildApprovePayload', () => {
    it('builds approve payload', () => {
      const payload = buildApprovePayload(123, { level: 2, notes: '欢迎加入' })
      expect(payload).toEqual({
        applicationId: 123,
        action: 'approve',
        level: 2,
        reason: '欢迎加入'
      })
    })

    it('defaults level to 1', () => {
      const payload = buildApprovePayload(123, { notes: 'ok' })
      expect(payload.level).toBe(1)
    })

    it('defaults notes to empty string', () => {
      const payload = buildApprovePayload(123, {})
      expect(payload.reason).toBe('')
    })
  })

  describe('buildRejectPayload', () => {
    it('builds reject payload', () => {
      const payload = buildRejectPayload(456, '资质不符')
      expect(payload).toEqual({
        applicationId: 456,
        action: 'reject',
        reason: '资质不符'
      })
    })

    it('defaults reason to empty string', () => {
      const payload = buildRejectPayload(456)
      expect(payload.reason).toBe('')
    })
  })

  describe('getMockPendingApplications', () => {
    it('returns array of applications', () => {
      const apps = getMockPendingApplications()
      expect(Array.isArray(apps)).toBe(true)
      expect(apps.length).toBe(2)
    })

    it('each application has required fields', () => {
      const apps = getMockPendingApplications()
      apps.forEach(app => {
        expect(app).toHaveProperty('id')
        expect(app).toHaveProperty('username')
        expect(app).toHaveProperty('status')
        expect(app.status).toBe('pending')
      })
    })
  })

  describe('getAvatarText', () => {
    it('returns first character of username', () => {
      expect(getAvatarText('张三')).toBe('张')
    })

    it('returns fallback for null username', () => {
      expect(getAvatarText(null)).toBe('申')
    })

    it('returns fallback for undefined username', () => {
      expect(getAvatarText(undefined)).toBe('申')
    })

    it('returns custom fallback', () => {
      expect(getAvatarText(null, 'X')).toBe('X')
    })
  })
})
