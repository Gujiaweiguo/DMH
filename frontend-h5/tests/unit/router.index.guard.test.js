import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockState = vi.hoisted(() => ({
  guard: null,
  resolveGuardNavigation: vi.fn(),
  createRouter: vi.fn(),
  createWebHistory: vi.fn(),
}))

vi.mock('vue-router', () => ({
  createWebHistory: (...args) => {
    mockState.createWebHistory(...args)
    return { type: 'web-history' }
  },
  createRouter: (options) => {
    mockState.createRouter(options)
    return {
      beforeEach: (guardFn) => {
        mockState.guard = guardFn
      },
    }
  },
}))

vi.mock('../../src/router/guard.js', () => ({
  resolveGuardNavigation: mockState.resolveGuardNavigation,
}))

vi.mock('../../src/views/brand/Analytics.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/CampaignEditorVant.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/CampaignPageDesigner.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/Campaigns.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/Dashboard.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/DistributorApproval.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/DistributorLevelRewards.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/Distributors.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/Login.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/Materials.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/MemberDetail.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/Members.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/Orders.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/PosterRecords.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/Promoters.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/Settings.vue', () => ({ default: {} }))
vi.mock('../../src/views/brand/VerificationRecords.vue', () => ({ default: {} }))
vi.mock('../../src/views/CampaignDetail.vue', () => ({ default: {} }))
vi.mock('../../src/views/CampaignForm.vue', () => ({ default: {} }))
vi.mock('../../src/views/CampaignList.vue', () => ({ default: {} }))
vi.mock('../../src/views/FeedbackCenter.vue', () => ({ default: {} }))
vi.mock('../../src/views/distributor/DistributorApply.vue', () => ({ default: {} }))
vi.mock('../../src/views/distributor/DistributorCenter.vue', () => ({ default: {} }))
vi.mock('../../src/views/distributor/DistributorPromotion.vue', () => ({ default: {} }))
vi.mock('../../src/views/distributor/DistributorRewards.vue', () => ({ default: {} }))
vi.mock('../../src/views/distributor/DistributorSubordinates.vue', () => ({ default: {} }))
vi.mock('../../src/views/distributor/Login.vue', () => ({ default: {} }))
vi.mock('../../src/views/MyOrders.vue', () => ({ default: {} }))
vi.mock('../../src/views/order/OrderVerification.vue', () => ({ default: {} }))
vi.mock('../../src/views/PaymentQrcode.vue', () => ({ default: {} }))
vi.mock('../../src/views/poster/PosterGenerator.vue', () => ({ default: {} }))
vi.mock('../../src/views/Success.vue', () => ({ default: {} }))

import router from '../../src/router/index.js'

const createLocalStorageMock = (initial = {}) => {
  const store = { ...initial }
  return {
    getItem: vi.fn((key) => (key in store ? store[key] : null)),
    setItem: vi.fn((key, value) => {
      store[key] = String(value)
    }),
    removeItem: vi.fn((key) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      Object.keys(store).forEach((key) => {
        delete store[key]
      })
    }),
  }
}

describe('router index guard integration', () => {
  beforeEach(() => {
    mockState.resolveGuardNavigation.mockReset()
    global.localStorage = createLocalStorageMock()
  })

  it('registers beforeEach guard during router creation', () => {
    expect(router).toBeTruthy()
    expect(mockState.createRouter).toHaveBeenCalledTimes(1)
    expect(typeof mockState.guard).toBe('function')
  })

  it('passes localStorage context to resolveGuardNavigation', () => {
    global.localStorage = createLocalStorageMock({
      dmh_token: 'token-1',
      dmh_user_role: 'brand_admin',
      dmh_user_info: '{"brandIds":[1]}',
    })
    mockState.resolveGuardNavigation.mockReturnValue(null)
    const next = vi.fn()
    const to = { path: '/brand/dashboard', meta: { requiresAuth: true, hasBrand: true } }

    mockState.guard(to, {}, next)

    expect(mockState.resolveGuardNavigation).toHaveBeenCalledWith(
      to,
      'token-1',
      'brand_admin',
      '{"brandIds":[1]}',
    )
    expect(next).toHaveBeenCalledWith()
  })

  it('redirects when resolver returns a path', () => {
    mockState.resolveGuardNavigation.mockReturnValue('/brand/login')
    const next = vi.fn()

    mockState.guard({ path: '/brand/orders', meta: { requiresAuth: true } }, {}, next)

    expect(next).toHaveBeenCalledWith('/brand/login')
  })
})
