import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  PHONE_REGEX,
  validatePhone,
  validateFormField,
  validateForm,
  loadSourceFromStorage,
  buildOrderPayload,
  initializeFormData,
  getSubmitButtonText
} from '../../src/views/campaignForm.logic.js'

describe('campaignForm.logic', () => {
  describe('PHONE_REGEX', () => {
    it('should match valid Chinese phone numbers', () => {
      expect(PHONE_REGEX.test('13800138000')).toBe(true)
      expect(PHONE_REGEX.test('15912345678')).toBe(true)
      expect(PHONE_REGEX.test('18888888888')).toBe(true)
    })

    it('should not match invalid phone numbers', () => {
      expect(PHONE_REGEX.test('12800138000')).toBe(false)
      expect(PHONE_REGEX.test('1380013800')).toBe(false)
      expect(PHONE_REGEX.test('138001380001')).toBe(false)
      expect(PHONE_REGEX.test('abc')).toBe(false)
      expect(PHONE_REGEX.test('')).toBe(false)
    })
  })

  describe('validatePhone', () => {
    it('should return error when phone is empty', () => {
      const result = validatePhone('')
      expect(result.valid).toBe(false)
      expect(result.error).toBe('请输入手机号')
    })

    it('should return error when phone is null', () => {
      const result = validatePhone(null)
      expect(result.valid).toBe(false)
      expect(result.error).toBe('请输入手机号')
    })

    it('should return error when phone format is invalid', () => {
      const result = validatePhone('12345678901')
      expect(result.valid).toBe(false)
      expect(result.error).toBe('请输入正确的手机号')
    })

    it('should return valid when phone is correct', () => {
      const result = validatePhone('13800138000')
      expect(result.valid).toBe(true)
      expect(result.error).toBeNull()
    })
  })

  describe('validateFormField', () => {
    it('should return error when value is empty', () => {
      const result = validateFormField('姓名', '')
      expect(result.valid).toBe(false)
      expect(result.error).toBe('请填写姓名')
    })

    it('should return error when value is only whitespace', () => {
      const result = validateFormField('姓名', '   ')
      expect(result.valid).toBe(false)
      expect(result.error).toBe('请填写姓名')
    })

    it('should return valid when value is provided', () => {
      const result = validateFormField('姓名', '张三')
      expect(result.valid).toBe(true)
      expect(result.error).toBeNull()
    })
  })

  describe('validateForm', () => {
    it('should return errors for empty form', () => {
      const result = validateForm('', {}, ['姓名', '地址'])
      expect(result.valid).toBe(false)
      expect(result.errors.length).toBeGreaterThan(0)
    })

    it('should return valid when all fields are filled', () => {
      const result = validateForm('13800138000', { 姓名: '张三', 地址: '北京' }, ['姓名', '地址'])
      expect(result.valid).toBe(true)
      expect(result.errors.length).toBe(0)
    })

    it('should return errors for missing form fields', () => {
      const result = validateForm('13800138000', { 姓名: '' }, ['姓名', '地址'])
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('请填写姓名')
      expect(result.errors).toContain('请填写地址')
    })

    it('should handle empty formFields array', () => {
      const result = validateForm('13800138000', {}, [])
      expect(result.valid).toBe(true)
    })

    it('should handle null formFields', () => {
      const result = validateForm('13800138000', {}, null)
      expect(result.valid).toBe(true)
    })
  })

  describe('loadSourceFromStorage', () => {
    let localStorageMock

    beforeEach(() => {
      localStorageMock = {
        store: {},
        getItem: vi.fn((key) => localStorageMock.store[key] || null),
        setItem: vi.fn((key, value) => { localStorageMock.store[key] = value }),
        clear: vi.fn(() => { localStorageMock.store = {} })
      }
      vi.stubGlobal('localStorage', localStorageMock)
    })

    afterEach(() => {
      vi.unstubAllGlobals()
    })

    it('should return empty source when nothing stored', () => {
      const result = loadSourceFromStorage()
      expect(result).toEqual({ c_id: '', u_id: '' })
    })

    it('should return stored source', () => {
      localStorageMock.store['dmh_source'] = JSON.stringify({ c_id: '123', u_id: '456' })
      const result = loadSourceFromStorage()
      expect(result).toEqual({ c_id: '123', u_id: '456' })
    })

    it('should handle invalid JSON', () => {
      localStorageMock.store['dmh_source'] = 'invalid json'
      const result = loadSourceFromStorage()
      expect(result).toEqual({ c_id: '', u_id: '' })
    })
  })

  describe('buildOrderPayload', () => {
    it('should build correct order payload', () => {
      const result = buildOrderPayload('123', '13800138000', { name: '张三' }, { u_id: '456' })
      expect(result).toEqual({
        campaignId: 123,
        phone: '13800138000',
        formData: { name: '张三' },
        referrerId: 456
      })
    })

    it('should handle empty u_id', () => {
      const result = buildOrderPayload('123', '13800138000', {}, { u_id: '' })
      expect(result.referrerId).toBe(0)
    })

    it('should handle missing u_id', () => {
      const result = buildOrderPayload('123', '13800138000', {}, {})
      expect(result.referrerId).toBe(0)
    })
  })

  describe('initializeFormData', () => {
    it('should initialize empty form data for fields', () => {
      const result = initializeFormData(['姓名', '地址'])
      expect(result).toEqual({ 姓名: '', 地址: '' })
    })

    it('should return empty object for empty array', () => {
      const result = initializeFormData([])
      expect(result).toEqual({})
    })

    it('should return empty object for null', () => {
      const result = initializeFormData(null)
      expect(result).toEqual({})
    })
  })

  describe('getSubmitButtonText', () => {
    it('should return submitting text when submitting', () => {
      expect(getSubmitButtonText(true)).toBe('提交中...')
    })

    it('should return normal text when not submitting', () => {
      expect(getSubmitButtonText(false)).toBe('提交报名')
    })
  })
})
