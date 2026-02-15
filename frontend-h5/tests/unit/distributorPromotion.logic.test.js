import { describe, it, expect } from 'vitest'
import {
  PROMOTION_GUIDE,
  buildCampaignOptions,
  buildPosterDownloadName,
  buildQrcodeDownloadName,
  buildQrcodeUrl,
  canGenerateCampaignPoster,
  buildPosterPayload,
  buildLinkPayload,
  formatClickCount,
  getCampaignDisplayName,
  DEFAULT_POSTER_DATA,
  DEFAULT_LINK_DATA
} from '../../src/views/distributor/distributorPromotion.logic.js'

describe('distributorPromotion.logic', () => {
  describe('PROMOTION_GUIDE', () => {
    it('has 4 guide items', () => {
      expect(PROMOTION_GUIDE).toHaveLength(4)
    })
  })

  describe('buildCampaignOptions', () => {
    it('converts campaigns to options', () => {
      const campaigns = [
        { id: 1, name: 'Campaign A' },
        { id: 2, name: 'Campaign B' }
      ]
      const options = buildCampaignOptions(campaigns)
      expect(options).toEqual([
        { text: 'Campaign A', value: 1 },
        { text: 'Campaign B', value: 2 }
      ])
    })

    it('returns empty array for null', () => {
      expect(buildCampaignOptions(null)).toEqual([])
    })

    it('returns empty array for non-array', () => {
      expect(buildCampaignOptions('not array')).toEqual([])
    })
  })

  describe('buildPosterDownloadName', () => {
    it('creates filename with timestamp', () => {
      const name = buildPosterDownloadName('活动海报')
      expect(name).toMatch(/^活动海报_\d+\.png$/)
    })
  })

  describe('buildQrcodeDownloadName', () => {
    it('creates filename with link code', () => {
      const name = buildQrcodeDownloadName('ABC123')
      expect(name).toBe('推广二维码_ABC123.png')
    })
  })

  describe('buildQrcodeUrl', () => {
    it('returns qrcodeUrl if provided', () => {
      const url = buildQrcodeUrl('https://example.com', 'https://qrcode.url')
      expect(url).toBe('https://qrcode.url')
    })

    it('generates qrserver URL if no qrcodeUrl', () => {
      const url = buildQrcodeUrl('https://example.com', null)
      expect(url).toContain('qrserver.com')
      expect(url).toContain(encodeURIComponent('https://example.com'))
    })
  })

  describe('canGenerateCampaignPoster', () => {
    it('returns true for valid campaignId', () => {
      expect(canGenerateCampaignPoster(123)).toBe(true)
    })

    it('returns false for null', () => {
      expect(canGenerateCampaignPoster(null)).toBe(false)
    })

    it('returns false for undefined', () => {
      expect(canGenerateCampaignPoster(undefined)).toBe(false)
    })

    it('returns false for 0', () => {
      expect(canGenerateCampaignPoster(0)).toBe(false)
    })
  })

  describe('buildPosterPayload', () => {
    it('builds campaign payload', () => {
      const payload = buildPosterPayload('campaign', 123)
      expect(payload).toEqual({ type: 'campaign', campaignId: 123 })
    })

    it('builds distributor payload without campaignId', () => {
      const payload = buildPosterPayload('distributor')
      expect(payload).toEqual({ type: 'distributor' })
    })
  })

  describe('buildLinkPayload', () => {
    it('builds link payload', () => {
      expect(buildLinkPayload(456)).toEqual({ campaignId: 456 })
    })
  })

  describe('formatClickCount', () => {
    it('formats click count', () => {
      expect(formatClickCount(10)).toBe('点击 10 次')
    })
  })

  describe('getCampaignDisplayName', () => {
    it('returns campaignName if present', () => {
      const link = { campaignId: 1, campaignName: 'Test Campaign' }
      expect(getCampaignDisplayName(link, [])).toBe('Test Campaign')
    })

    it('finds campaign name from campaigns list', () => {
      const link = { campaignId: 2 }
      const campaigns = [{ id: 2, name: 'Found Campaign' }]
      expect(getCampaignDisplayName(link, campaigns)).toBe('Found Campaign')
    })

    it('returns fallback name if not found', () => {
      const link = { campaignId: 999 }
      expect(getCampaignDisplayName(link, [])).toBe('活动 999')
    })
  })

  describe('DEFAULT_POSTER_DATA', () => {
    it('has expected structure', () => {
      expect(DEFAULT_POSTER_DATA).toHaveProperty('campaignName')
      expect(DEFAULT_POSTER_DATA).toHaveProperty('campaignDesc')
      expect(DEFAULT_POSTER_DATA).toHaveProperty('distributorName')
      expect(DEFAULT_POSTER_DATA).toHaveProperty('campaignCount')
    })
  })

  describe('DEFAULT_LINK_DATA', () => {
    it('has expected structure', () => {
      expect(DEFAULT_LINK_DATA).toHaveProperty('link')
      expect(DEFAULT_LINK_DATA).toHaveProperty('linkCode')
      expect(DEFAULT_LINK_DATA).toHaveProperty('qrcodeUrl')
    })
  })
})
