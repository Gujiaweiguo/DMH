import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  getDefaultForm,
  hasDistributorRole,
  saveDistributorAuth,
  isAlreadyLoggedIn,
  getLoginErrorMessage,
  validateLoginForm,
  getButtonText
} from '../../src/views/distributor/distributorLogin.logic.js'

describe('distributorLogin.logic', () => {
  describe('getDefaultForm', () => {
    it('returns empty form', () => {
      const form = getDefaultForm()
      expect(form).toEqual({ username: '', password: '' })
    })
  })

  describe('hasDistributorRole', () => {
    it('returns true when roles include distributor', () => {
      expect(hasDistributorRole({ roles: ['distributor', 'user'] })).toBe(true)
    })

    it('returns false when roles do not include distributor', () => {
      expect(hasDistributorRole({ roles: ['user', 'admin'] })).toBe(false)
    })

    it('returns false for null data', () => {
      expect(hasDistributorRole(null)).toBe(false)
    })

    it('returns false when roles is missing', () => {
      expect(hasDistributorRole({ token: 'abc' })).toBe(false)
    })
  })

  describe('saveDistributorAuth', () => {
    beforeEach(() => {
      vi.stubGlobal('localStorage', {
        setItem: vi.fn()
      })
    })

    it('saves auth data to localStorage', () => {
      const data = { token: 'test-token', roles: ['distributor'] }
      saveDistributorAuth(data)
      expect(localStorage.setItem).toHaveBeenCalledWith('dmh_token', 'test-token')
      expect(localStorage.setItem).toHaveBeenCalledWith('dmh_user_role', 'distributor')
    })

    it('handles null data', () => {
      saveDistributorAuth(null)
      expect(localStorage.setItem).not.toHaveBeenCalled()
    })
  })

  describe('isAlreadyLoggedIn', () => {
    it('returns true when token and distributor role exist', () => {
      vi.stubGlobal('localStorage', {
        getItem: vi.fn((key) => {
          if (key === 'dmh_token') return 'token'
          if (key === 'dmh_user_role') return 'distributor'
          return null
        })
      })
      expect(isAlreadyLoggedIn()).toBe(true)
    })

    it('returns false when no token', () => {
      vi.stubGlobal('localStorage', {
        getItem: vi.fn(() => null)
      })
      expect(isAlreadyLoggedIn()).toBe(false)
    })

    it('returns false when role is not distributor', () => {
      vi.stubGlobal('localStorage', {
        getItem: vi.fn((key) => {
          if (key === 'dmh_token') return 'token'
          if (key === 'dmh_user_role') return 'admin'
          return null
        })
      })
      expect(isAlreadyLoggedIn()).toBe(false)
    })
  })

  describe('getLoginErrorMessage', () => {
    it('returns error message', () => {
      expect(getLoginErrorMessage({ message: 'Test error' })).toBe('Test error')
    })

    it('returns default for null error', () => {
      expect(getLoginErrorMessage(null)).toBe('登录失败，请重试')
    })

    it('returns default for error without message', () => {
      expect(getLoginErrorMessage({})).toBe('登录失败，请重试')
    })
  })

  describe('validateLoginForm', () => {
    it('returns valid for complete form', () => {
      const result = validateLoginForm({ username: 'user', password: 'pass' })
      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    it('returns error for empty username', () => {
      const result = validateLoginForm({ username: '', password: 'pass' })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('请输入用户名')
    })

    it('returns error for empty password', () => {
      const result = validateLoginForm({ username: 'user', password: '' })
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('请输入密码')
    })

    it('returns error for null form', () => {
      const result = validateLoginForm(null)
      expect(result.valid).toBe(false)
    })
  })

  describe('getButtonText', () => {
    it('returns "登录" when not loading', () => {
      expect(getButtonText(false)).toBe('登录')
    })

    it('returns "登录中..." when loading', () => {
      expect(getButtonText(true)).toBe('登录中...')
    })
  })
})
