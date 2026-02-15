import { describe, expect, it } from 'vitest'
import {
  buildDistributorListQuery,
  countDistributorsByStatus,
  formatDistributorTime,
  getDistributorStatusText,
  mergeDistributorList,
  parseCurrentBrandId,
} from '../../src/views/brand/distributors.logic.js'

describe('distributors logic', () => {
  it('parses current brand id from storage first', () => {
    expect(parseCurrentBrandId('11', '{"brandIds":[9]}')).toBe(11)
  })

  it('falls back to first brand id from user info', () => {
    expect(parseCurrentBrandId('', '{"brandIds":[7,8]}')).toBe(7)
    expect(parseCurrentBrandId(null, '{"brandIds":[]}')).toBe(0)
    expect(parseCurrentBrandId(null, '{bad-json')).toBe(0)
  })

  it('formats distributor time and handles empty', () => {
    expect(formatDistributorTime('2026-02-13 10:00:00')).toContain('2026')
    expect(formatDistributorTime('')).toBe('-')
  })

  it('maps distributor status text', () => {
    expect(getDistributorStatusText('active')).toBe('启用')
    expect(getDistributorStatusText('suspended')).toBe('停用')
    expect(getDistributorStatusText('custom')).toBe('custom')
    expect(getDistributorStatusText('')).toBe('-')
  })

  it('counts distributors by status', () => {
    const list = [
      { id: 1, status: 'active' },
      { id: 2, status: 'suspended' },
      { id: 3, status: 'active' },
    ]
    expect(countDistributorsByStatus(list, 'active')).toBe(2)
    expect(countDistributorsByStatus(list, 'suspended')).toBe(1)
  })

  it('builds distributor query params', () => {
    const query = buildDistributorListQuery({
      page: 2,
      pageSize: 20,
      brandId: 9,
      keyword: 'alice ',
      status: 'active',
      level: 3,
    })
    expect(query.get('page')).toBe('2')
    expect(query.get('pageSize')).toBe('20')
    expect(query.get('brandId')).toBe('9')
    expect(query.get('keyword')).toBe('alice')
    expect(query.get('status')).toBe('active')
    expect(query.get('level')).toBe('3')
  })

  it('merges distributor list and computes hasMore', () => {
    const replaced = mergeDistributorList({
      currentList: [{ id: 1 }],
      incomingList: [{ id: 2 }, { id: 3 }],
      replace: true,
      total: 5,
    })
    expect(replaced).toEqual({
      distributors: [{ id: 2 }, { id: 3 }],
      hasMore: true,
    })

    const appended = mergeDistributorList({
      currentList: [{ id: 1 }, { id: 2 }],
      incomingList: [{ id: 3 }],
      replace: false,
      total: 3,
    })
    expect(appended).toEqual({
      distributors: [{ id: 1 }, { id: 2 }, { id: 3 }],
      hasMore: false,
    })
  })
})
