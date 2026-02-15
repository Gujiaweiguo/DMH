import { describe, it, expect } from 'vitest'
import {
  PAY_TYPES,
  STATUS_TYPES,
  STATUS_TEXTS,
  PAY_TYPE_TEXTS,
  PAY_ACCOUNT_LABELS,
  PAY_REAL_NAME_LABELS,
  getStatusType,
  getStatusText,
  getPayTypeText,
  getPayAccountLabel,
  getPayRealNameLabel,
  getDefaultWithdrawalForm,
  validateAmount,
  buildWithdrawalPayload,
  formatBalance,
  formatDate,
  calculateNewBalance
} from '../../src/views/distributor/distributorWithdrawals.logic.js'

describe('distributorWithdrawals.logic', () => {
  describe('PAY_TYPES', () => {
    it('has 3 pay types', () => {
      expect(PAY_TYPES).toHaveLength(3)
    })
  })

  describe('getStatusType', () => {
    it('returns correct types', () => {
      expect(getStatusType('pending')).toBe('warning')
      expect(getStatusType('approved')).toBe('primary')
      expect(getStatusType('rejected')).toBe('danger')
      expect(getStatusType('completed')).toBe('success')
    })

    it('returns default for unknown', () => {
      expect(getStatusType('unknown')).toBe('default')
    })
  })

  describe('getStatusText', () => {
    it('returns correct texts', () => {
      expect(getStatusText('pending')).toBe('待审核')
      expect(getStatusText('approved')).toBe('已批准')
      expect(getStatusText('rejected')).toBe('已拒绝')
      expect(getStatusText('completed')).toBe('已完成')
    })

    it('returns original status for unknown', () => {
      expect(getStatusText('unknown')).toBe('unknown')
    })
  })

  describe('getPayTypeText', () => {
    it('returns correct texts', () => {
      expect(getPayTypeText('wechat')).toBe('微信')
      expect(getPayTypeText('alipay')).toBe('支付宝')
      expect(getPayTypeText('bank')).toBe('银行卡')
    })

    it('returns original for unknown', () => {
      expect(getPayTypeText('unknown')).toBe('unknown')
    })
  })

  describe('getPayAccountLabel', () => {
    it('returns correct labels', () => {
      expect(getPayAccountLabel('wechat')).toBe('微信号')
      expect(getPayAccountLabel('alipay')).toBe('支付宝账号')
      expect(getPayAccountLabel('bank')).toBe('银行卡号')
    })

    it('returns default for unknown', () => {
      expect(getPayAccountLabel('unknown')).toBe('提现账号')
    })
  })

  describe('getPayRealNameLabel', () => {
    it('returns correct labels', () => {
      expect(getPayRealNameLabel('wechat')).toBe('微信昵称')
      expect(getPayRealNameLabel('alipay')).toBe('支付宝姓名')
      expect(getPayRealNameLabel('bank')).toBe('银行卡持卡人')
    })

    it('returns default for unknown', () => {
      expect(getPayRealNameLabel('unknown')).toBe('真实姓名')
    })
  })

  describe('getDefaultWithdrawalForm', () => {
    it('returns default form', () => {
      const form = getDefaultWithdrawalForm()
      expect(form).toEqual({
        amount: '',
        payType: 'wechat',
        payAccount: '',
        payRealName: ''
      })
    })
  })

  describe('validateAmount', () => {
    it('returns true for valid amount', () => {
      expect(validateAmount('100', 200)).toBe(true)
    })

    it('returns error for non-numeric', () => {
      expect(validateAmount('abc', 200)).toBe('提现金额必须大于0')
    })

    it('returns error for zero', () => {
      expect(validateAmount('0', 200)).toBe('提现金额必须大于0')
    })

    it('returns error for negative', () => {
      expect(validateAmount('-10', 200)).toBe('提现金额必须大于0')
    })

    it('returns error when exceeding balance', () => {
      expect(validateAmount('300', 200)).toBe('提现金额不能超过可用余额')
    })
  })

  describe('buildWithdrawalPayload', () => {
    it('builds payload from form', () => {
      const form = {
        amount: '100.50',
        payType: 'alipay',
        payAccount: 'test@example.com',
        payRealName: '张三'
      }
      const payload = buildWithdrawalPayload(form)
      expect(payload).toEqual({
        amount: 100.50,
        payType: 'alipay',
        payAccount: 'test@example.com',
        payRealName: '张三'
      })
    })
  })

  describe('formatBalance', () => {
    it('formats balance to 2 decimals', () => {
      expect(formatBalance(123.456)).toBe('123.46')
    })

    it('handles zero', () => {
      expect(formatBalance(0)).toBe('0.00')
    })

    it('handles null', () => {
      expect(formatBalance(null)).toBe('0.00')
    })
  })

  describe('formatDate', () => {
    it('formats date correctly', () => {
      const result = formatDate('2024-06-15T14:30:00')
      expect(result).toMatch(/2024-06-15/)
      expect(result).toMatch(/14:30/)
    })

    it('returns empty string for null', () => {
      expect(formatDate(null)).toBe('')
    })

    it('returns empty string for undefined', () => {
      expect(formatDate(undefined)).toBe('')
    })
  })

  describe('calculateNewBalance', () => {
    it('subtracts withdrawal amount', () => {
      expect(calculateNewBalance(100, 30)).toBe(70)
    })

    it('does not go below zero', () => {
      expect(calculateNewBalance(50, 100)).toBe(0)
    })
  })

  describe('constants', () => {
    it('STATUS_TYPES has all expected keys', () => {
      expect(Object.keys(STATUS_TYPES)).toContain('pending')
      expect(Object.keys(STATUS_TYPES)).toContain('completed')
    })

    it('PAY_ACCOUNT_LABELS has all pay types', () => {
      expect(Object.keys(PAY_ACCOUNT_LABELS)).toHaveLength(3)
    })
  })
})
