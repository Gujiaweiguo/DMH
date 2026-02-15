import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
  enrichVerificationRecords,
  loadVerificationRecordsData,
  unverifyAndReloadVerificationRecords,
} from '../../src/views/brand/verificationRecords.actions.js'

describe('verificationRecords actions', () => {
  let orderApi

  beforeEach(() => {
    orderApi = {
      getVerificationRecords: vi.fn(),
      getOrders: vi.fn(),
      unverifyOrder: vi.fn(),
    }
  })

  it('enrichVerificationRecords returns empty list for empty input', async () => {
    const result = await enrichVerificationRecords(orderApi, [])
    expect(result).toEqual([])
    expect(orderApi.getOrders).not.toHaveBeenCalled()
  })

  it('enrichVerificationRecords merges records with orders', async () => {
    orderApi.getOrders.mockResolvedValue({
      orders: [{ id: 1, status: 'paid', amount: 99, phone: '13800001111' }],
    })

    const result = await enrichVerificationRecords(orderApi, [
      { orderId: 1, verifiedBy: 9, verificationStatus: 'verified', verifiedAt: '2026-02-13T08:00:00.000Z' },
    ])

    expect(result[0]).toMatchObject({
      orderStatus: 'paid',
      orderAmount: 99,
      userPhone: '13800001111',
      verifiedByName: '用户9',
    })
  })

  it('enrichVerificationRecords falls back when getOrders fails', async () => {
    orderApi.getOrders.mockRejectedValue(new Error('db down'))
    const base = [{ orderId: 2, verificationStatus: 'unverified' }]

    const result = await enrichVerificationRecords(orderApi, base)
    expect(result).toEqual(base)
  })

  it('loadVerificationRecordsData returns merged records and stats', async () => {
    orderApi.getVerificationRecords.mockResolvedValue({
      records: [
        { orderId: 1, verificationStatus: 'verified', verifiedAt: '2026-02-13T08:00:00.000Z', orderAmount: 0, verifiedBy: 1 },
      ],
    })
    orderApi.getOrders.mockResolvedValue({
      orders: [{ id: 1, status: 'paid', amount: 50, phone: '13800001111' }],
    })

    const result = await loadVerificationRecordsData(orderApi)

    expect(result.records).toHaveLength(1)
    expect(result.records[0]).toMatchObject({ orderStatus: 'paid', orderAmount: 50 })
    expect(result.stats.total).toBe(1)
    expect(result.stats.totalAmount).toBe(50)
  })

  it('loadVerificationRecordsData throws when getVerificationRecords fails', async () => {
    orderApi.getVerificationRecords.mockRejectedValue(new Error('network'))

    await expect(loadVerificationRecordsData(orderApi)).rejects.toThrow('network')
  })

  it('unverifyAndReloadVerificationRecords calls unverify then reload', async () => {
    orderApi.unverifyOrder.mockResolvedValue({ code: 0 })
    orderApi.getVerificationRecords.mockResolvedValue({ records: [] })

    const result = await unverifyAndReloadVerificationRecords(orderApi, 88)

    expect(orderApi.unverifyOrder).toHaveBeenCalledWith('', { orderId: 88 })
    expect(orderApi.getVerificationRecords).toHaveBeenCalledTimes(1)
    expect(result).toEqual({
      records: [],
      stats: { total: 0, today: 0, totalAmount: 0 },
    })
  })
})
