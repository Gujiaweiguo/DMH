import { describe, it, expect } from 'vitest'
import {
  DISTRIBUTOR_DESCRIPTION,
  getDefaultForm,
  buildBrandOptions,
  canSubmit,
  buildApplyPayload,
  validateForm,
  TERMS_CONTENT
} from '../../src/views/distributor/distributorApply.logic.js'

describe('distributorApply.logic', () => {
  describe('DISTRIBUTOR_DESCRIPTION', () => {
    it('is a non-empty string', () => {
      expect(typeof DISTRIBUTOR_DESCRIPTION).toBe('string')
      expect(DISTRIBUTOR_DESCRIPTION.length).toBeGreaterThan(0)
    })
  })

  describe('getDefaultForm', () => {
    it('returns form with null brandId', () => {
      const form = getDefaultForm()
      expect(form.brandId).toBeNull()
    })

    it('returns form with empty reason', () => {
      const form = getDefaultForm()
      expect(form.reason).toBe('')
    })
  })

  describe('buildBrandOptions', () => {
    it('converts brands to options', () => {
      const brands = [
        { id: 1, name: 'Brand A' },
        { id: 2, name: 'Brand B' }
      ]
      const options = buildBrandOptions(brands)
      expect(options).toEqual([
        { text: 'Brand A', value: 1 },
        { text: 'Brand B', value: 2 }
      ])
    })

    it('returns empty array for null', () => {
      expect(buildBrandOptions(null)).toEqual([])
    })

    it('returns empty array for undefined', () => {
      expect(buildBrandOptions(undefined)).toEqual([])
    })

    it('returns empty array for non-array', () => {
      expect(buildBrandOptions('not array')).toEqual([])
    })

    it('handles empty array', () => {
      expect(buildBrandOptions([])).toEqual([])
    })
  })

  describe('canSubmit', () => {
    it('returns true for valid form and agreed terms', () => {
      const form = { brandId: 1, reason: 'I want to promote' }
      expect(canSubmit(form, true)).toBe(true)
    })

    it('returns false when brandId is null', () => {
      const form = { brandId: null, reason: 'test' }
      expect(canSubmit(form, true)).toBe(false)
    })

    it('returns false when reason is empty', () => {
      const form = { brandId: 1, reason: '' }
      expect(canSubmit(form, true)).toBe(false)
    })

    it('returns false when reason is whitespace only', () => {
      const form = { brandId: 1, reason: '   ' }
      expect(canSubmit(form, true)).toBe(false)
    })

    it('returns false when terms not agreed', () => {
      const form = { brandId: 1, reason: 'test' }
      expect(canSubmit(form, false)).toBe(false)
    })

    it('returns false for null form', () => {
      expect(canSubmit(null, true)).toBe(false)
    })
  })

  describe('buildApplyPayload', () => {
    it('builds payload from form', () => {
      const form = { brandId: 123, reason: 'test reason' }
      expect(buildApplyPayload(form)).toEqual({
        brandId: 123,
        reason: 'test reason'
      })
    })

    it('handles null form', () => {
      expect(buildApplyPayload(null)).toEqual({
        brandId: undefined,
        reason: ''
      })
    })
  })

  describe('validateForm', () => {
    it('returns valid for complete form', () => {
      const form = { brandId: 1, reason: 'valid reason' }
      const result = validateForm(form, true)
      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    it('returns error for missing brandId', () => {
      const form = { brandId: null, reason: 'test' }
      const result = validateForm(form, true)
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('请选择品牌')
    })

    it('returns error for empty reason', () => {
      const form = { brandId: 1, reason: '' }
      const result = validateForm(form, true)
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('请填写申请理由')
    })

    it('returns error for not agreeing terms', () => {
      const form = { brandId: 1, reason: 'test' }
      const result = validateForm(form, false)
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('请先同意分销商协议')
    })

    it('returns multiple errors', () => {
      const result = validateForm({ brandId: null, reason: '' }, false)
      expect(result.errors.length).toBeGreaterThan(1)
    })
  })

  describe('TERMS_CONTENT', () => {
    it('has title', () => {
      expect(TERMS_CONTENT.title).toBe('分销商协议')
    })

    it('has 4 sections', () => {
      expect(TERMS_CONTENT.sections).toHaveLength(4)
    })

    it('each section has title and items', () => {
      TERMS_CONTENT.sections.forEach(section => {
        expect(section.title).toBeTruthy()
        expect(Array.isArray(section.items)).toBe(true)
        expect(section.items.length).toBeGreaterThan(0)
      })
    })
  })
})
