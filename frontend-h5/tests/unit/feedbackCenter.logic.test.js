import { describe, it, expect } from 'vitest'
import {
  CATEGORY_OPTIONS,
  PRIORITY_OPTIONS,
  RATING_OPTIONS,
  STATUS_TEXT_MAP,
  getDefaultForm,
  statusText,
  formatDate,
  validateFeedbackForm,
  buildFeedbackPayload,
  resetForm,
  getSubmitButtonText,
  getHelpfulCount
} from '../../src/views/feedbackCenter.logic.js'

describe('feedbackCenter.logic', () => {
  describe('CATEGORY_OPTIONS', () => {
    it('has 4 options', () => {
      expect(CATEGORY_OPTIONS).toHaveLength(4)
    })
  })

  describe('PRIORITY_OPTIONS', () => {
    it('has 3 options', () => {
      expect(PRIORITY_OPTIONS).toHaveLength(3)
    })
  })

  describe('RATING_OPTIONS', () => {
    it('has 6 options', () => {
      expect(RATING_OPTIONS).toHaveLength(6)
    })
  })

  describe('STATUS_TEXT_MAP', () => {
    it('has correct mappings', () => {
      expect(STATUS_TEXT_MAP.pending).toBe('待处理')
      expect(STATUS_TEXT_MAP.reviewing).toBe('处理中')
      expect(STATUS_TEXT_MAP.resolved).toBe('已解决')
      expect(STATUS_TEXT_MAP.closed).toBe('已关闭')
    })
  })

  describe('getDefaultForm', () => {
    it('returns default form', () => {
      const form = getDefaultForm()
      expect(form.category).toBe('poster')
      expect(form.priority).toBe('medium')
      expect(form.title).toBe('')
      expect(form.content).toBe('')
      expect(form.rating).toBeNull()
    })
  })

  describe('statusText', () => {
    it('returns correct text for pending', () => {
      expect(statusText('pending')).toBe('待处理')
    })

    it('returns original status for unknown', () => {
      expect(statusText('unknown')).toBe('unknown')
    })
  })

  describe('formatDate', () => {
    it('formats date string', () => {
      const result = formatDate('2024-06-15T10:30:00')
      expect(result).toContain('2024')
    })

    it('returns "-" for null', () => {
      expect(formatDate(null)).toBe('-')
    })

    it('returns "-" for undefined', () => {
      expect(formatDate(undefined)).toBe('-')
    })
  })

  describe('validateFeedbackForm', () => {
    it('returns valid for complete form', () => {
      const form = { title: 'Test', content: 'Content' }
      const result = validateFeedbackForm(form)
      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    it('returns error for missing title', () => {
      const form = { title: '', content: 'Content' }
      const result = validateFeedbackForm(form)
      expect(result.valid).toBe(false)
    })

    it('returns error for missing content', () => {
      const form = { title: 'Test', content: '' }
      const result = validateFeedbackForm(form)
      expect(result.valid).toBe(false)
    })

    it('returns error for null form', () => {
      const result = validateFeedbackForm(null)
      expect(result.valid).toBe(false)
    })
  })

  describe('buildFeedbackPayload', () => {
    it('builds payload from form', () => {
      const form = {
        category: 'payment',
        priority: 'high',
        title: 'Test',
        content: 'Content',
        rating: 5
      }
      const payload = buildFeedbackPayload(form)
      expect(payload.category).toBe('payment')
      expect(payload.priority).toBe('high')
      expect(payload.title).toBe('Test')
      expect(payload.content).toBe('Content')
      expect(payload.rating).toBe(5)
      expect(payload.featureUseCase).toBe('h5_feedback_center')
    })
  })

  describe('resetForm', () => {
    it('returns default form', () => {
      const form = resetForm()
      expect(form).toEqual(getDefaultForm())
    })
  })

  describe('getSubmitButtonText', () => {
    it('returns 提交反馈 when not submitting', () => {
      expect(getSubmitButtonText(false)).toBe('提交反馈')
    })

    it('returns 提交中... when submitting', () => {
      expect(getSubmitButtonText(true)).toBe('提交中...')
    })
  })

  describe('getHelpfulCount', () => {
    it('returns helpfulCount for helpful type', () => {
      const faq = { helpfulCount: 10, notHelpfulCount: 2 }
      expect(getHelpfulCount(faq, 'helpful')).toBe(10)
    })

    it('returns notHelpfulCount for not_helpful type', () => {
      const faq = { helpfulCount: 10, notHelpfulCount: 2 }
      expect(getHelpfulCount(faq, 'not_helpful')).toBe(2)
    })

    it('returns 0 for missing counts', () => {
      expect(getHelpfulCount({}, 'helpful')).toBe(0)
    })
  })
})
