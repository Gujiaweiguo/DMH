import { describe, expect, it } from 'vitest'
import {
  buildExportRequestBody,
  buildMembersQueryParams,
  formatMemberDate,
  getBrandIdFromStorage,
  mergeMembersPage,
} from '../../src/views/brand/members.logic.js'

describe('members logic', () => {
  it('resolves brand id from direct storage first', () => {
    expect(getBrandIdFromStorage('12', '{"brandIds":[99]}')).toBe(12)
  })

  it('resolves brand id from user info fallback', () => {
    expect(getBrandIdFromStorage('', '{"brandIds":[7,8]}')).toBe(7)
    expect(getBrandIdFromStorage(null, '{"brandIds":[]}')).toBeNull()
    expect(getBrandIdFromStorage(null, '{bad-json')).toBeNull()
  })

  it('builds members query params', () => {
    const params = buildMembersQueryParams({
      page: 2,
      pageSize: 20,
      brandId: 9,
      keyword: 'alice',
      source: 'wechat',
      status: 'active',
    })

    expect(params.get('page')).toBe('2')
    expect(params.get('pageSize')).toBe('20')
    expect(params.get('brandId')).toBe('9')
    expect(params.get('keyword')).toBe('alice')
    expect(params.get('source')).toBe('wechat')
    expect(params.get('status')).toBe('active')
  })

  it('merges first page and sets pagination state', () => {
    const merged = mergeMembersPage({
      currentMembers: [{ id: 1 }],
      incomingMembers: [{ id: 2 }, { id: 3 }],
      page: 1,
      pageSize: 2,
    })

    expect(merged).toEqual({
      members: [{ id: 2 }, { id: 3 }],
      finished: false,
      nextPage: 2,
    })
  })

  it('appends later pages and marks finished when short page', () => {
    const merged = mergeMembersPage({
      currentMembers: [{ id: 1 }],
      incomingMembers: [{ id: 2 }],
      page: 2,
      pageSize: 2,
    })

    expect(merged).toEqual({
      members: [{ id: 1 }, { id: 2 }],
      finished: true,
      nextPage: 2,
    })
  })

  it('builds export body with serialized filters', () => {
    const body = buildExportRequestBody({
      brandId: 10,
      reason: 'analysis',
      keyword: 'kw',
      source: 'wechat',
      status: 'active',
    })

    expect(body.brandId).toBe(10)
    expect(body.reason).toBe('analysis')
    expect(JSON.parse(body.filters)).toEqual({
      keyword: 'kw',
      source: 'wechat',
      status: 'active',
    })
  })

  it('formats member date and handles empty input', () => {
    expect(formatMemberDate('2026-02-13T10:00:00.000Z')).toBe('2026-02-13')
    expect(formatMemberDate('')).toBe('-')
  })
})
