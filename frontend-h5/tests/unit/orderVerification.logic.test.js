import { describe, expect, it } from 'vitest'
import {
  buildUnverifiedOrderState,
  buildVerifiedOrderState,
  extractCodeFromText,
  getOrderStatusText,
  getPaymentStatusText,
  getVerifyStatusText,
  isOperationSuccess,
} from '../../src/views/order/orderVerification.logic.js'

describe('orderVerification logic', () => {
  it('extracts code from url query', () => {
    const code = extractCodeFromText('https://dmh.example.com/verify?code=abc123')
    expect(code).toBe('abc123')
  })

  it('extracts orderCode from url query', () => {
    const code = extractCodeFromText('https://dmh.example.com/verify?orderCode=ord001')
    expect(code).toBe('ord001')
  })

  it('extracts verification_code from url query', () => {
    const code = extractCodeFromText('https://dmh.example.com/verify?verification_code=vc123')
    expect(code).toBe('vc123')
  })

  it('extracts verificationCode from url query', () => {
    const code = extractCodeFromText('https://dmh.example.com/verify?verificationCode=vc456')
    expect(code).toBe('vc456')
  })

  it('returns original text for url without code param', () => {
    const text = 'https://dmh.example.com/verify?other=value'
    expect(extractCodeFromText(text)).toBe(text)
  })

  it('returns original text for non-url input', () => {
    expect(extractCodeFromText('raw-code')).toBe('raw-code')
  })

  it('returns empty string for empty input', () => {
    expect(extractCodeFromText('')).toBe('')
    expect(extractCodeFromText(null)).toBe('')
  })

  it('maps payment status text', () => {
    expect(getPaymentStatusText('paid')).toBe('已支付')
    expect(getPaymentStatusText('unpaid')).toBe('未支付')
    expect(getPaymentStatusText('refunded')).toBe('已退款')
    expect(getPaymentStatusText('custom')).toBe('custom')
  })

  it('maps verify status and order status text', () => {
    expect(getVerifyStatusText('verified')).toBe('已核销')
    expect(getVerifyStatusText('unverified')).toBe('未核销')
    expect(getVerifyStatusText('other')).toBe('other')
    expect(getOrderStatusText('verified')).toBe('已核销')
  })

  it('judges operation success response correctly', () => {
    expect(isOperationSuccess({ code: 0, data: { id: 1 } })).toBe(true)
    expect(isOperationSuccess({ code: 0, data: null })).toBe(false)
    expect(isOperationSuccess({ code: 1, data: { id: 1 } })).toBe(false)
    expect(isOperationSuccess(null)).toBe(false)
  })

  it('builds verified order state without mutating original', () => {
    const original = { id: 1, verifyStatus: 'unverified' }
    const next = buildVerifiedOrderState(original, 'ok', '2026-02-13T15:00:00.000Z')

    expect(next).toEqual({
      id: 1,
      verifyStatus: 'verified',
      verifiedBy: '品牌管理员',
      verifiedAt: '2026-02-13T15:00:00.000Z',
      notes: 'ok',
    })
    expect(original).toEqual({ id: 1, verifyStatus: 'unverified' })
  })

  it('builds unverified order state without mutating original', () => {
    const original = { id: 2, verifyStatus: 'verified' }
    const next = buildUnverifiedOrderState(original, 'manual', '2026-02-13T15:10:00.000Z')

    expect(next).toEqual({
      id: 2,
      verifyStatus: 'unverified',
      cancelledBy: '品牌管理员',
      cancelledAt: '2026-02-13T15:10:00.000Z',
      cancelReason: 'manual',
    })
    expect(original).toEqual({ id: 2, verifyStatus: 'verified' })
  })
})
