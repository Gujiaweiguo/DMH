import { describe, it, expect } from 'vitest'
import {
  DEFAULT_H5_URL,
  DEFAULT_API_URL,
  STATUS_CHECKING,
  STATUS_OK,
  STATUS_FAIL,
  STATUS_RUNNING,
  LOGIN_TEST_NOT_TESTED,
  LOGIN_TEST_TESTING,
  LOGIN_TEST_SUCCESS,
  LOGIN_TEST_FAIL,
  getConnectionStatusText,
  getLoginTestResultText,
  buildLoginTestPayload,
  isConnectionOk,
  getDefaultState,
  LOGIN_ROUTE
} from '../../src/views/apiTest.logic.js'

describe('apiTest.logic', () => {
  describe('Constants', () => {
    it('should have correct default URLs', () => {
      expect(DEFAULT_H5_URL).toBe('http://localhost:3100')
      expect(DEFAULT_API_URL).toBe('/api/v1')
    })

    it('should have correct status constants', () => {
      expect(STATUS_CHECKING).toBe('检测中...')
      expect(STATUS_OK).toBe('✅ 连接正常')
      expect(STATUS_FAIL).toBe('❌ 连接失败')
      expect(STATUS_RUNNING).toBe('✅ 运行正常')
    })

    it('should have correct login test constants', () => {
      expect(LOGIN_TEST_NOT_TESTED).toBe('未测试')
      expect(LOGIN_TEST_TESTING).toBe('测试中...')
      expect(LOGIN_TEST_SUCCESS).toBe('✅ 登录成功')
      expect(LOGIN_TEST_FAIL).toBe('❌ 登录失败')
    })
  })

  describe('getConnectionStatusText', () => {
    it('should return OK status when ok is true', () => {
      expect(getConnectionStatusText(true, 200)).toBe(STATUS_OK)
    })

    it('should return error with status when ok is false', () => {
      expect(getConnectionStatusText(false, 500)).toBe('❌ 连接异常 (500)')
    })

    it('should return fail status when no status', () => {
      expect(getConnectionStatusText(false, null)).toBe(STATUS_FAIL)
    })
  })

  describe('getLoginTestResultText', () => {
    it('should return empty string for null', () => {
      expect(getLoginTestResultText(null)).toBe('')
    })

    it('should format result correctly', () => {
      const result = {
        username: 'admin',
        roles: ['admin', 'user']
      }
      expect(getLoginTestResultText(result)).toBe('用户: admin, 角色: admin, user')
    })

    it('should handle empty roles', () => {
      const result = {
        username: 'admin',
        roles: []
      }
      expect(getLoginTestResultText(result)).toBe('用户: admin, 角色: ')
    })
  })

  describe('buildLoginTestPayload', () => {
    it('should return test credentials', () => {
      const payload = buildLoginTestPayload()
      expect(payload.username).toBe('test')
      expect(payload.password).toBe('test')
    })
  })

  describe('isConnectionOk', () => {
    it('should return true for ok response', () => {
      expect(isConnectionOk({ ok: true, status: 200 })).toBe(true)
    })

    it('should return true for 401 status', () => {
      expect(isConnectionOk({ ok: false, status: 401 })).toBe(true)
    })

    it('should return false for other status', () => {
      expect(isConnectionOk({ ok: false, status: 500 })).toBe(false)
    })
  })

  describe('getDefaultState', () => {
    it('should return default state', () => {
      const state = getDefaultState()
      expect(state.h5Status).toBe(STATUS_CHECKING)
      expect(state.h5Url).toBe(DEFAULT_H5_URL)
      expect(state.apiStatus).toBe(STATUS_CHECKING)
      expect(state.apiUrl).toBe(DEFAULT_API_URL)
      expect(state.loginTestStatus).toBe(LOGIN_TEST_NOT_TESTED)
      expect(state.loginTestResult).toBe('')
      expect(state.errorMessage).toBe('')
    })
  })

  describe('LOGIN_ROUTE', () => {
    it('should be correct route', () => {
      expect(LOGIN_ROUTE).toBe('/brand/login')
    })
  })
})
