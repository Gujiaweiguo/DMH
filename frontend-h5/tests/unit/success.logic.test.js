import { describe, it, expect, vi } from 'vitest'
import { HOME_ROUTE, goHome } from '../../src/views/success.logic.js'

describe('success.logic', () => {
  describe('HOME_ROUTE', () => {
    it('should be root path', () => {
      expect(HOME_ROUTE).toBe('/')
    })
  })

  describe('goHome', () => {
    it('should call router.push with HOME_ROUTE', () => {
      const router = { push: vi.fn() }
      goHome(router)
      expect(router.push).toHaveBeenCalledWith(HOME_ROUTE)
    })
  })
})
