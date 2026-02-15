import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  TEST_USERNAME,
  TEST_PASSWORD,
  getDefaultForm,
  quickFillTestAccount,
  validateLoginForm,
  hasBrandAccess,
  getFirstBrandId,
  saveLoginInfo,
  clearLoginInfo,
  getLoginButtonText,
  buildLoginError
} from '../../src/views/brand/login.logic.js'

describe('brand/login.logic', () => {
  let localStorageMock

  beforeEach(() => {
    localStorageMock = {
      store: {},
      getItem: vi.fn((key) => localStorageMock.store[key] || null),
      setItem: vi.fn((key, value) => { localStorageMock.store[key] = value }),
      removeItem: vi.fn((key) => { delete localStorageMock.store[key] }),
      clear: vi.fn(() => { localStorageMock.store = {} })
    }
    vi.stubGlobal('localStorage', localStorageMock)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  describe('TEST_USERNAME and TEST_PASSWORD', () => {
    it('should have test credentials', () => {
      expect(TEST_USERNAME).toBe('brand_manager')
      expect(TEST_PASSWORD).toBe('123456')
    })
  })

  describe('getDefaultForm', () => {
    it('should return empty form', () => {
      const result = getDefaultForm()
      expect(result).toEqual({ username: '', password: '' })
    })
  })

  describe('quickFillTestAccount', () => {
    it('should fill form with test credentials', () => {
      const form = { username: '', password: '' }
      quickFillTestAccount(form)
      expect(form.username).toBe(TEST_USERNAME)
      expect(form.password).toBe(TEST_PASSWORD)
    })

    it('should return the form object', () => {
      const form = { username: '', password: '' }
      const result = quickFillTestAccount(form)
      expect(result).toBe(form)
    })
  })

  describe('validateLoginForm', () => {
    it('should return error for empty username', () => {
      const result = validateLoginForm({ username: '', password: '123456' })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('请输入用户名')
    })

    it('should return error for empty password', () => {
      const result = validateLoginForm({ username: 'test', password: '' })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('请输入密码')
    })

    it('should return errors for both empty fields', () => {
      const result = validateLoginForm({ username: '', password: '' })
      expect(result.valid).toBe(false)
      expect(result.errors.length).toBe(2)
    })

    it('should return valid for filled form', () => {
      const result = validateLoginForm({ username: 'test', password: '123456' })
      expect(result.valid).toBe(true)
      expect(result.errors.length).toBe(0)
    })

    it('should handle null form', () => {
      const result = validateLoginForm(null)
      expect(result.valid).toBe(false)
    })

    it('should handle whitespace only inputs', () => {
      const result = validateLoginForm({ username: '   ', password: '   ' })
      expect(result.valid).toBe(false)
      expect(result.errors.length).toBe(2)
    })
  })

  describe('hasBrandAccess', () => {
    it('should return false for null data', () => {
      expect(hasBrandAccess(null)).toBe(false)
    })

    it('should return false for data without brandIds', () => {
      expect(hasBrandAccess({})).toBe(false)
    })

    it('should return false for empty brandIds', () => {
      expect(hasBrandAccess({ brandIds: [] })).toBe(false)
    })

    it('should return false for non-array brandIds', () => {
      expect(hasBrandAccess({ brandIds: '123' })).toBe(false)
    })

    it('should return true for valid brandIds', () => {
      expect(hasBrandAccess({ brandIds: [1, 2, 3] })).toBe(true)
    })
  })

  describe('getFirstBrandId', () => {
    it('should return null for null', () => {
      expect(getFirstBrandId(null)).toBeNull()
    })

    it('should return null for empty array', () => {
      expect(getFirstBrandId([])).toBeNull()
    })

    it('should return first element', () => {
      expect(getFirstBrandId([1, 2, 3])).toBe(1)
    })

    it('should return null for non-array', () => {
      expect(getFirstBrandId('123')).toBeNull()
    })
  })

  describe('saveLoginInfo', () => {
    it('should save login info to localStorage', () => {
      const data = { token: 'test-token', brandIds: [1, 2] }
      const result = saveLoginInfo(data, 1)
      expect(result).toBe(true)
      expect(localStorage.getItem('dmh_token')).toBe('test-token')
      expect(localStorage.getItem('dmh_user_role')).toBe('participant')
      expect(localStorage.getItem('dmh_current_brand_id')).toBe('1')
    })

    it('should save user info as JSON', () => {
      const data = { token: 'test-token', brandIds: [1] }
      saveLoginInfo(data, 1)
      const userInfo = JSON.parse(localStorage.getItem('dmh_user_info'))
      expect(userInfo.token).toBe('test-token')
    })

    it('should return false on localStorage error', () => {
      localStorageMock.setItem.mockImplementation(() => { throw new Error('QuotaExceededError') })
      const data = { token: 'test-token', brandIds: [1] }
      const result = saveLoginInfo(data, 1)
      expect(result).toBe(false)
    })
  })

  describe('clearLoginInfo', () => {
    it('should clear all login info', () => {
      localStorage.setItem('dmh_token', 'test')
      localStorage.setItem('dmh_user_role', 'test')
      localStorage.setItem('dmh_user_info', '{}')
      localStorage.setItem('dmh_current_brand_id', '1')
      
      const result = clearLoginInfo()
      expect(result).toBe(true)
      expect(localStorage.getItem('dmh_token')).toBeNull()
      expect(localStorage.getItem('dmh_user_role')).toBeNull()
      expect(localStorage.getItem('dmh_user_info')).toBeNull()
      expect(localStorage.getItem('dmh_current_brand_id')).toBeNull()
    })

    it('should return false on localStorage error', () => {
      localStorageMock.removeItem.mockImplementation(() => { throw new Error('Storage error') })
      const result = clearLoginInfo()
      expect(result).toBe(false)
    })
  })

  describe('getLoginButtonText', () => {
    it('should return loading text when loading', () => {
      expect(getLoginButtonText(true)).toBe('登录中...')
    })

    it('should return normal text when not loading', () => {
      expect(getLoginButtonText(false)).toBe('登录')
    })
  })

  describe('buildLoginError', () => {
    it('should return default message for null', () => {
      expect(buildLoginError(null)).toBe('登录失败，请重试')
    })

    it('should return default message for undefined', () => {
      expect(buildLoginError(undefined)).toBe('登录失败，请重试')
    })

    it('should return string error as is', () => {
      expect(buildLoginError('用户名错误')).toBe('用户名错误')
    })

    it('should extract message from error object', () => {
      const error = { message: '网络错误' }
      expect(buildLoginError(error)).toBe('网络错误')
    })

    it('should return default for error without message', () => {
      const error = { code: 500 }
      expect(buildLoginError(error)).toBe('登录失败，请重试')
    })
  })
})
