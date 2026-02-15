import { describe, expect, it } from 'vitest'
import {
  calculateVerificationStats,
  filterVerificationRecords,
  formatDateTime,
  getVerificationStatusText,
  mergeVerificationRecordsWithOrders,
} from '../../src/views/brand/verificationRecords.logic.js'

describe('verificationRecords logic', () => {
  const records = [
    { orderId: 1001, verificationStatus: 'verified', userPhone: '13800001111', verifiedBy: 1, verifiedAt: '2026-02-13T08:00:00.000Z', orderAmount: 12 },
    { orderId: 1002, verificationStatus: 'unverified', userPhone: '13900002222', verifiedBy: null, verifiedAt: null, orderAmount: 5 },
    { orderId: 2001, verificationStatus: 'verified', userPhone: '13700003333', verifiedBy: 2, verifiedAt: '2026-02-12T08:00:00.000Z', orderAmount: 20 },
  ]

  it('maps verification status text', () => {
    expect(getVerificationStatusText('verified')).toBe('已核销')
    expect(getVerificationStatusText('unverified')).toBe('未核销')
    expect(getVerificationStatusText('cancelled')).toBe('已取消')
    expect(getVerificationStatusText('custom')).toBe('custom')
  })

  it('formats datetime and handles empty', () => {
    expect(formatDateTime('2026-02-13T08:00:00.000Z')).toContain('2026')
    expect(formatDateTime('')).toBe('')
  })

  it('filters by status and keyword', () => {
    const byStatus = filterVerificationRecords(records, 'verified', '')
    expect(byStatus).toHaveLength(2)

    const byKeywordOrder = filterVerificationRecords(records, 'all', '1002')
    expect(byKeywordOrder).toHaveLength(1)
    expect(byKeywordOrder[0].orderId).toBe(1002)

    const byKeywordPhone = filterVerificationRecords(records, 'all', '1370000')
    expect(byKeywordPhone).toHaveLength(1)
    expect(byKeywordPhone[0].orderId).toBe(2001)
  })

  it('merges order fields into verification records', () => {
    const merged = mergeVerificationRecordsWithOrders(
      [{ orderId: 1001, verifiedBy: 3 }, { orderId: 9999, verifiedBy: null }],
      [{ id: 1001, status: 'paid', amount: 88, phone: '13800138000' }],
    )

    expect(merged[0]).toMatchObject({
      orderStatus: 'paid',
      orderAmount: 88,
      userPhone: '13800138000',
      verifiedByName: '用户3',
    })
    expect(merged[1]).toMatchObject({
      orderStatus: '',
      orderAmount: 0,
      userPhone: '',
      verifiedByName: '用户-',
    })
  })

  it('calculates verified stats with fixed today string', () => {
    const stats = calculateVerificationStats(records, '2026/2/13')
    expect(stats).toEqual({
      total: 2,
      today: 1,
      totalAmount: 32,
    })
  })
})
