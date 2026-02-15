import { describe, it, expect } from 'vitest'
import {
  REQUIRED_RULE,
  createMinLengthRule,
  createMaxLengthRule,
  createRangeRule,
  PHONE_RULE,
  EMAIL_RULE,
  validateField,
  validateForm,
  getFirstError,
  hasErrors
} from '../../src/utils/formValidation.logic.js'

describe('formValidation.logic', () => {
  describe('Constants', () => {
    it('should have required rule', () => {
      expect(REQUIRED_RULE.required).toBe(true)
    })

    it('should have phone pattern', () => {
      expect(PHONE_RULE.pattern.test('13800138000')).toBe(true)
      expect(PHONE_RULE.pattern.test('123')).toBe(false)
    })

    it('should have email pattern', () => {
      expect(EMAIL_RULE.pattern.test('test@example.com')).toBe(true)
      expect(EMAIL_RULE.pattern.test('invalid')).toBe(false)
    })
  })

  describe('createMinLengthRule', () => {
    it('should create rule with min length', () => {
      const rule = createMinLengthRule(5)
      expect(rule.min).toBe(5)
      expect(rule.message).toContain('5')
    })
  })

  describe('createMaxLengthRule', () => {
    it('should create rule with max length', () => {
      const rule = createMaxLengthRule(100)
      expect(rule.max).toBe(100)
      expect(rule.message).toContain('100')
    })
  })

  describe('createRangeRule', () => {
    it('should create rule with range', () => {
      const rule = createRangeRule(1, 100)
      expect(rule.min).toBe(1)
      expect(rule.max).toBe(100)
    })
  })

  describe('validateField', () => {
    it('should return valid for empty rules', () => {
      expect(validateField('test', []).valid).toBe(true)
    })

    it('should return invalid for required field with empty value', () => {
      const result = validateField('', [REQUIRED_RULE])
      expect(result.valid).toBe(false)
    })

    it('should return valid for required field with value', () => {
      const result = validateField('test', [REQUIRED_RULE])
      expect(result.valid).toBe(true)
    })

    it('should return invalid for phone pattern mismatch', () => {
      const result = validateField('123', [PHONE_RULE])
      expect(result.valid).toBe(false)
    })

    it('should return valid for matching phone pattern', () => {
      const result = validateField('13800138000', [PHONE_RULE])
      expect(result.valid).toBe(true)
    })

    it('should return invalid for min length violation', () => {
      const result = validateField('ab', [createMinLengthRule(5)])
      expect(result.valid).toBe(false)
    })
  })

  describe('validateForm', () => {
    const schema = {
      name: [REQUIRED_RULE],
      phone: [PHONE_RULE]
    }

    it('should return valid for correct data', () => {
      const result = validateForm({ name: 'test', phone: '13800138000' }, schema)
      expect(result.valid).toBe(true)
      expect(Object.keys(result.errors).length).toBe(0)
    })

    it('should return invalid for missing required field', () => {
      const result = validateForm({ name: '', phone: '13800138000' }, schema)
      expect(result.valid).toBe(false)
      expect(result.errors.name).toBeDefined()
    })

    it('should return invalid for invalid phone', () => {
      const result = validateForm({ name: 'test', phone: '123' }, schema)
      expect(result.valid).toBe(false)
      expect(result.errors.phone).toBeDefined()
    })
  })

  describe('getFirstError', () => {
    it('should return first error', () => {
      expect(getFirstError({ a: 'error1', b: 'error2' })).toBe('error1')
    })

    it('should return null for empty errors', () => {
      expect(getFirstError({})).toBeNull()
    })
  })

  describe('hasErrors', () => {
    it('should return true for non-empty errors', () => {
      expect(hasErrors({ a: 'error' })).toBe(true)
    })

    it('should return false for empty errors', () => {
      expect(hasErrors({})).toBe(false)
    })
  })
})
