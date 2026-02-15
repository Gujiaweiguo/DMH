import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  safeGetItem,
  safeSetItem,
  safeRemoveItem,
  safeGetJson,
  safeSetJson,
  clearAll,
  getStorageKeys,
  getStorageSize,
  hasKey,
  STORAGE_KEYS,
  getToken,
  setToken,
  removeToken,
  getUserRole,
  setUserRole,
  getUserInfo,
  setUserInfo,
  getBrandId,
  setBrandId,
  isLoggedIn
} from '../../src/utils/storage.logic.js'

describe('storage.logic', () => {
  let localStorageMock

  beforeEach(() => {
    localStorageMock = {
      store: {},
      getItem: vi.fn((key) => localStorageMock.store[key] || null),
      setItem: vi.fn((key, value) => { localStorageMock.store[key] = value }),
      removeItem: vi.fn((key) => { delete localStorageMock.store[key] }),
      clear: vi.fn(() => { localStorageMock.store = {} }),
      key: vi.fn((index) => Object.keys(localStorageMock.store)[index] || null)
    }
    Object.defineProperty(localStorageMock, 'length', {
      get: () => Object.keys(localStorageMock.store).length
    })
    vi.stubGlobal('localStorage', localStorageMock)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  describe('safeGetItem', () => {
    it('should get item', () => {
      localStorageMock.store.test = 'value'
      expect(safeGetItem('test')).toBe('value')
    })

    it('should return null for missing key', () => {
      expect(safeGetItem('missing')).toBeNull()
    })

    it('should return null on error', () => {
      localStorageMock.getItem.mockImplementation(() => { throw new Error('fail') })
      expect(safeGetItem('test')).toBeNull()
    })
  })

  describe('safeSetItem', () => {
    it('should set item', () => {
      expect(safeSetItem('test', 'value')).toBe(true)
      expect(localStorageMock.store.test).toBe('value')
    })

    it('should return false on error', () => {
      localStorageMock.setItem.mockImplementation(() => { throw new Error('fail') })
      expect(safeSetItem('test', 'value')).toBe(false)
    })
  })

  describe('safeRemoveItem', () => {
    it('should remove item', () => {
      localStorageMock.store.test = 'value'
      expect(safeRemoveItem('test')).toBe(true)
      expect(localStorageMock.store.test).toBeUndefined()
    })

    it('should return false on error', () => {
      localStorageMock.removeItem.mockImplementation(() => { throw new Error('fail') })
      expect(safeRemoveItem('test')).toBe(false)
    })
  })

  describe('safeGetJson', () => {
    it('should parse JSON', () => {
      localStorageMock.store.test = '{"a":1}'
      expect(safeGetJson('test')).toEqual({ a: 1 })
    })

    it('should return default for invalid JSON', () => {
      localStorageMock.store.test = 'invalid'
      expect(safeGetJson('test', {})).toEqual({})
    })

    it('should return default on error', () => {
      localStorageMock.getItem.mockImplementation(() => { throw new Error('fail') })
      expect(safeGetJson('test', { default: true })).toEqual({ default: true })
    })
  })

  describe('safeSetJson', () => {
    it('should stringify and set', () => {
      expect(safeSetJson('test', { a: 1 })).toBe(true)
      expect(localStorageMock.store.test).toBe('{"a":1}')
    })

    it('should return false on error', () => {
      localStorageMock.setItem.mockImplementation(() => { throw new Error('fail') })
      expect(safeSetJson('test', { a: 1 })).toBe(false)
    })
  })

  describe('clearAll', () => {
    it('should clear all items', () => {
      localStorageMock.store.test = 'value'
      expect(clearAll()).toBe(true)
      expect(Object.keys(localStorageMock.store).length).toBe(0)
    })

    it('should return false on error', () => {
      localStorageMock.clear.mockImplementation(() => { throw new Error('fail') })
      expect(clearAll()).toBe(false)
    })
  })

  describe('getStorageKeys', () => {
    it('should return keys', () => {
      localStorageMock.store.a = '1'
      localStorageMock.store.b = '2'
      const keys = getStorageKeys()
      expect(keys).toContain('a')
      expect(keys).toContain('b')
    })

    it('should return empty array on key error', () => {
      localStorageMock.store.a = '1'
      localStorageMock.key.mockImplementation(() => { throw new Error('fail') })
      expect(getStorageKeys()).toEqual([])
    })
  })

  describe('getStorageSize', () => {
    it('should return size', () => {
      localStorageMock.store.test = '12345'
      const size = getStorageSize()
      expect(size).toBeGreaterThan(0)
    })

    it('should return 0 on error', () => {
      localStorageMock.getItem.mockImplementation(() => { throw new Error('fail') })
      localStorageMock.store.test = '123'
      expect(getStorageSize()).toBe(0)
    })
  })

  describe('hasKey', () => {
    it('should return true for existing key', () => {
      localStorageMock.store.test = 'value'
      expect(hasKey('test')).toBe(true)
    })

    it('should return false for missing key', () => {
      expect(hasKey('missing')).toBe(false)
    })

    it('should return false on error', () => {
      localStorageMock.getItem.mockImplementation(() => { throw new Error('fail') })
      expect(hasKey('test')).toBe(false)
    })
  })

  describe('STORAGE_KEYS', () => {
    it('should have required keys', () => {
      expect(STORAGE_KEYS.TOKEN).toBe('dmh_token')
      expect(STORAGE_KEYS.USER_ROLE).toBe('dmh_user_role')
    })
  })

  describe('token helpers', () => {
    it('should set and get token', () => {
      setToken('abc123')
      expect(getToken()).toBe('abc123')
    })

    it('should remove token', () => {
      setToken('abc123')
      removeToken()
      expect(getToken()).toBeNull()
    })
  })

  describe('user role helpers', () => {
    it('should set and get user role', () => {
      setUserRole('admin')
      expect(getUserRole()).toBe('admin')
    })
  })

  describe('user info helpers', () => {
    it('should set and get user info', () => {
      setUserInfo({ name: 'test' })
      expect(getUserInfo()).toEqual({ name: 'test' })
    })
  })

  describe('brand id helpers', () => {
    it('should set and get brand id', () => {
      setBrandId(123)
      expect(getBrandId()).toBe('123')
    })
  })

  describe('isLoggedIn', () => {
    it('should return true when token exists', () => {
      setToken('abc')
      expect(isLoggedIn()).toBe(true)
    })

    it('should return false when no token', () => {
      expect(isLoggedIn()).toBe(false)
    })
  })
})
