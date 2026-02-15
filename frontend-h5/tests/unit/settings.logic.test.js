import { describe, expect, it } from 'vitest'
import {
  getDefaultBrandInfo,
  getDefaultNotificationSettings,
  getDefaultPasswordForm,
  getDefaultRewardSettings,
  getDefaultSyncSettings,
  getSyncStatusText,
  validatePasswordForm,
} from '../../src/views/brand/settings.logic.js'

describe('settings logic', () => {
  it('provides default settings objects', () => {
    expect(getDefaultBrandInfo()).toMatchObject({
      name: '示例品牌',
      phone: '400-123-4567',
    })
    expect(getDefaultRewardSettings()).toEqual({
      defaultRate: 20,
      minWithdraw: 100,
      settlementType: 'instant',
    })
    expect(getDefaultNotificationSettings()).toEqual({
      newOrder: true,
      newPromoter: true,
      dailyReport: false,
      email: 'admin@brand.com',
    })
    expect(getDefaultSyncSettings()).toEqual({
      status: 'connected',
      frequency: 'realtime',
      dataTypes: ['orders', 'users'],
    })
    expect(getDefaultPasswordForm()).toEqual({
      oldPassword: '',
      newPassword: '',
      confirmPassword: '',
    })
  })

  it('maps sync status text', () => {
    expect(getSyncStatusText('connected')).toBe('已连接')
    expect(getSyncStatusText('disconnected')).toBe('未连接')
    expect(getSyncStatusText('error')).toBe('连接错误')
    expect(getSyncStatusText('custom')).toBe('custom')
  })

  it('validates password form fields', () => {
    expect(
      validatePasswordForm({ oldPassword: '', newPassword: '123456', confirmPassword: '123456' }),
    ).toBe('请填写完整的密码信息')

    expect(
      validatePasswordForm({ oldPassword: 'old', newPassword: '123456', confirmPassword: '654321' }),
    ).toBe('两次输入的新密码不一致')

    expect(
      validatePasswordForm({ oldPassword: 'old', newPassword: '123456', confirmPassword: '123456' }),
    ).toBe('')
  })
})
