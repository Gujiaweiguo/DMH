import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockState = vi.hoisted(() => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

vi.mock('../../src/services/api.js', () => ({
  api: mockState.api,
}))

import { orderApi } from '../../src/services/brandApi.js'

describe('brandApi orderApi verify/unverify', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockState.api.post.mockResolvedValue({ code: 0 })
  })

  it('verifyOrder sends code-only payload', async () => {
    await orderApi.verifyOrder('code_123')

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/verify', { code: 'code_123' })
  })

  it('scanOrderCode passes flat query params', async () => {
    await orderApi.scanOrderCode('scan_123')

    expect(mockState.api.get).toHaveBeenCalledWith('/orders/scan', { code: 'scan_123' })
  })

  it('verifyOrder merges code with extra data', async () => {
    await orderApi.verifyOrder('code_123', { notes: 'ok' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/verify', {
      code: 'code_123',
      notes: 'ok',
    })
  })

  it('verifyOrder accepts object payload directly', async () => {
    await orderApi.verifyOrder({ code: 'code_456', notes: 'from-object' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/verify', {
      code: 'code_456',
      notes: 'from-object',
    })
  })

  it('unverifyOrder omits blank code and keeps data', async () => {
    await orderApi.unverifyOrder('', { orderId: 99, reason: 'manual' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/unverify', {
      orderId: 99,
      reason: 'manual',
    })
  })

  it('unverifyOrder sends code with extra data', async () => {
    await orderApi.unverifyOrder('code_789', { reason: 'mistake' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/unverify', {
      code: 'code_789',
      reason: 'mistake',
    })
  })

  it('verifyOrder accepts non-string code values', async () => {
    await orderApi.verifyOrder(1001, { notes: 'numeric-code' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/verify', {
      code: 1001,
      notes: 'numeric-code',
    })
  })

  it('unverifyOrder accepts non-string code values', async () => {
    await orderApi.unverifyOrder(1002, { reason: 'numeric-code' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/unverify', {
      code: 1002,
      reason: 'numeric-code',
    })
  })

  it('verifyOrder accepts object with no extra data', async () => {
    await orderApi.verifyOrder({ code: 'obj_code', notes: 'obj_notes' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/verify', {
      code: 'obj_code',
      notes: 'obj_notes',
    })
  })

  it('unverifyOrder accepts object payload directly', async () => {
    await orderApi.unverifyOrder({ code: 'obj_code', reason: 'obj_reason' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/unverify', {
      code: 'obj_code',
      reason: 'obj_reason',
    })
  })

  it('verifyOrder handles whitespace-only code', async () => {
    await orderApi.verifyOrder('   ', { notes: 'whitespace' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/verify', {
      notes: 'whitespace',
    })
  })

  it('unverifyOrder handles whitespace-only code', async () => {
    await orderApi.unverifyOrder('   ', { reason: 'whitespace' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/unverify', {
      reason: 'whitespace',
    })
  })

  it('verifyOrder handles null code', async () => {
    await orderApi.verifyOrder(null, { notes: 'null-code' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/verify', {
      notes: 'null-code',
    })
  })

  it('unverifyOrder handles null code', async () => {
    await orderApi.unverifyOrder(null, { reason: 'null-code' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/unverify', {
      reason: 'null-code',
    })
  })

  it('verifyOrder handles undefined code', async () => {
    await orderApi.verifyOrder(undefined, { notes: 'undefined-code' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/verify', {
      notes: 'undefined-code',
    })
  })

  it('unverifyOrder handles undefined code', async () => {
    await orderApi.unverifyOrder(undefined, { reason: 'undefined-code' })

    expect(mockState.api.post).toHaveBeenCalledWith('/orders/unverify', {
      reason: 'undefined-code',
    })
  })
})
