import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockState = vi.hoisted(() => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
  appendSpy: vi.fn(),
}))

vi.mock('../../src/services/api.js', () => ({
  api: mockState.api,
}))

import {
  analyticsApi,
  authApi,
  campaignApi,
  materialApi,
  orderApi,
  posterApi,
  promoterApi,
  settingsApi,
} from '../../src/services/brandApi.js'

describe('brandApi wrappers', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    global.FormData = class {
      append(key, value) {
        mockState.appendSpy(key, value)
      }
    }
  })

  it('auth api calls expected endpoints', async () => {
    await authApi.login('admin', '123456')
    expect(mockState.api.post).toHaveBeenCalledWith('/auth/login', { username: 'admin', password: '123456' })

    await authApi.getUserInfo()
    expect(mockState.api.get).toHaveBeenCalledWith('/auth/userinfo')

    await authApi.logout()
    expect(mockState.api.post).toHaveBeenCalledWith('/auth/logout')
  })

  it('campaign api uses correct routes', async () => {
    await campaignApi.getCampaigns({ page: 1 })
    expect(mockState.api.get).toHaveBeenCalledWith('/campaigns', { page: 1 })

    await campaignApi.getCampaign(10)
    expect(mockState.api.get).toHaveBeenCalledWith('/campaign/detail/10')

    await campaignApi.createCampaign({ name: 'c1' })
    expect(mockState.api.post).toHaveBeenCalledWith('/campaign/create', { name: 'c1' })

    await campaignApi.updateCampaign(10, { name: 'c2' })
    expect(mockState.api.put).toHaveBeenCalledWith('/campaign/update/10', { name: 'c2' })

    await campaignApi.deleteCampaign(10)
    expect(mockState.api.delete).toHaveBeenCalledWith('/campaign/delete/10')

    await campaignApi.savePageConfig(10, { blocks: [] })
    expect(mockState.api.post).toHaveBeenCalledWith('/campaign/page-config/10', { blocks: [] })

    await campaignApi.getPageConfig(10)
    expect(mockState.api.get).toHaveBeenCalledWith('/campaign/page-config/10')

    await campaignApi.getPaymentQrcode(10)
    expect(mockState.api.get).toHaveBeenCalledWith('/campaigns/10/payment-qrcode')
  })

  it('order list/detail/status wrappers use expected payloads', async () => {
    await orderApi.getOrders({ page: 2 })
    expect(mockState.api.get).toHaveBeenCalledWith('/order/list', { page: 2 })

    await orderApi.getOrder(11)
    expect(mockState.api.get).toHaveBeenCalledWith('/order/detail/11')

    await orderApi.updateOrderStatus(11, 'paid')
    expect(mockState.api.put).toHaveBeenCalledWith('/order/status/11', { status: 'paid' })

    await orderApi.getVerificationRecords()
    expect(mockState.api.get).toHaveBeenCalledWith('/order/verification-records')
  })

  it('promoter api uses expected routes', async () => {
    await promoterApi.getPromoters({ page: 1 })
    expect(mockState.api.get).toHaveBeenCalledWith('/promoter/list', { page: 1 })

    await promoterApi.getPromoter(8)
    expect(mockState.api.get).toHaveBeenCalledWith('/promoter/detail/8')

    await promoterApi.generateLink(2, 3)
    expect(mockState.api.post).toHaveBeenCalledWith('/promoter/generate-link', { promoterId: 2, campaignId: 3 })

    await promoterApi.getRewards(2, { page: 1 })
    expect(mockState.api.get).toHaveBeenCalledWith('/promoter/rewards/2', { page: 1 })
  })

  it('analytics api sends period query object', async () => {
    await analyticsApi.getMetrics('week')
    expect(mockState.api.get).toHaveBeenCalledWith('/analytics/metrics', { period: 'week' })

    await analyticsApi.getTrends('month')
    expect(mockState.api.get).toHaveBeenCalledWith('/analytics/trends', { period: 'month' })

    await analyticsApi.getCampaignRanking('day')
    expect(mockState.api.get).toHaveBeenCalledWith('/analytics/campaign-ranking', { period: 'day' })

    await analyticsApi.getPromoterRanking('year')
    expect(mockState.api.get).toHaveBeenCalledWith('/analytics/promoter-ranking', { period: 'year' })
  })

  it('poster api wrappers use expected routes', async () => {
    await posterApi.getPosterTemplates()
    expect(mockState.api.get).toHaveBeenCalledWith('/poster/templates')

    await posterApi.generatePoster(6, 9)
    expect(mockState.api.post).toHaveBeenCalledWith('/campaigns/6/poster', { distributorId: 9 })

    await posterApi.getPosterFile('f.png')
    expect(mockState.api.get).toHaveBeenCalledWith('/poster/f.png')

    await posterApi.getPosterRecords()
    expect(mockState.api.get).toHaveBeenCalledWith('/poster/records')
  })

  it('material api upload uses FormData and multipart header', async () => {
    const file = { name: 'a.png' }
    await materialApi.uploadMaterial(file, 'image')

    expect(mockState.appendSpy).toHaveBeenCalledWith('file', file)
    expect(mockState.appendSpy).toHaveBeenCalledWith('type', 'image')
    expect(mockState.api.post).toHaveBeenCalledWith(
      '/material/upload',
      expect.any(FormData),
      { headers: { 'Content-Type': 'multipart/form-data' } },
    )

    await materialApi.getMaterials({ page: 1 })
    expect(mockState.api.get).toHaveBeenCalledWith('/material/list', { page: 1 })

    await materialApi.deleteMaterial(12)
    expect(mockState.api.delete).toHaveBeenCalledWith('/material/delete/12')
  })

  it('settings api wrappers use expected routes', async () => {
    await settingsApi.getBrandInfo()
    expect(mockState.api.get).toHaveBeenCalledWith('/brand/info')

    await settingsApi.updateBrandInfo({ name: 'brand' })
    expect(mockState.api.put).toHaveBeenCalledWith('/brand/info', { name: 'brand' })

    await settingsApi.getRewardSettings()
    expect(mockState.api.get).toHaveBeenCalledWith('/brand/reward-settings')

    await settingsApi.updateRewardSettings({ level: 2 })
    expect(mockState.api.put).toHaveBeenCalledWith('/brand/reward-settings', { level: 2 })

    await settingsApi.getNotificationSettings()
    expect(mockState.api.get).toHaveBeenCalledWith('/brand/notification-settings')

    await settingsApi.updateNotificationSettings({ sms: true })
    expect(mockState.api.put).toHaveBeenCalledWith('/brand/notification-settings', { sms: true })
  })
})
