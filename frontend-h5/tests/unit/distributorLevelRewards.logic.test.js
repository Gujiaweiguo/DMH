import { describe, it, expect } from 'vitest'
import {
  getDefaultRewards,
  clampRewardPercentage,
  validateRewards,
  buildSavePayload,
  mergeRewardsFromServer,
  calculateTotalReward,
  getLevelLabel,
  LEVEL_LABELS
} from '../../src/views/brand/distributorLevelRewards.logic.js'

describe('distributorLevelRewards.logic', () => {
  describe('getDefaultRewards', () => {
    it('returns 3 levels', () => {
      const rewards = getDefaultRewards()
      expect(rewards).toHaveLength(3)
    })

    it('each level has rewardPercentage 0', () => {
      const rewards = getDefaultRewards()
      rewards.forEach(r => {
        expect(r.rewardPercentage).toBe(0)
      })
    })

    it('levels are 1, 2, 3', () => {
      const rewards = getDefaultRewards()
      expect(rewards.map(r => r.level)).toEqual([1, 2, 3])
    })
  })

  describe('clampRewardPercentage', () => {
    it('returns 0 for non-numeric value', () => {
      expect(clampRewardPercentage('abc')).toBe(0)
    })

    it('returns 0 for null', () => {
      expect(clampRewardPercentage(null)).toBe(0)
    })

    it('returns 0 for undefined', () => {
      expect(clampRewardPercentage(undefined)).toBe(0)
    })

    it('returns 0 for negative values', () => {
      expect(clampRewardPercentage(-10)).toBe(0)
    })

    it('returns 100 for values over 100', () => {
      expect(clampRewardPercentage(150)).toBe(100)
    })

    it('returns the value for valid range', () => {
      expect(clampRewardPercentage(50)).toBe(50)
    })

    it('handles decimal values', () => {
      expect(clampRewardPercentage(5.5)).toBe(5.5)
    })
  })

  describe('validateRewards', () => {
    it('fails for empty array', () => {
      const result = validateRewards([])
      expect(result.valid).toBe(false)
      expect(result.errors).toContain('奖励配置不能为空')
    })

    it('fails for null', () => {
      const result = validateRewards(null)
      expect(result.valid).toBe(false)
    })

    it('fails for non-array', () => {
      const result = validateRewards('not array')
      expect(result.valid).toBe(false)
    })

    it('fails for invalid level (less than 1)', () => {
      const result = validateRewards([
        { level: 0, rewardPercentage: 5 }
      ])
      expect(result.valid).toBe(false)
      expect(result.errors.some(e => e.includes('级别必须大于0'))).toBe(true)
    })

    it('fails for invalid level (non-numeric)', () => {
      const result = validateRewards([
        { level: 'invalid', rewardPercentage: 5 }
      ])
      expect(result.valid).toBe(false)
    })

    it('fails for non-numeric percentage', () => {
      const result = validateRewards([
        { level: 1, rewardPercentage: 'not-a-number' }
      ])
      expect(result.valid).toBe(false)
      expect(result.errors.some(e => e.includes('奖励比例无效'))).toBe(true)
    })

    it('fails for invalid percentage (negative)', () => {
      const result = validateRewards([
        { level: 1, rewardPercentage: -5 }
      ])
      expect(result.valid).toBe(false)
    })

    it('fails for invalid percentage (over 100)', () => {
      const result = validateRewards([
        { level: 1, rewardPercentage: 150 }
      ])
      expect(result.valid).toBe(false)
    })

    it('passes for valid rewards', () => {
      const result = validateRewards([
        { level: 1, rewardPercentage: 5 },
        { level: 2, rewardPercentage: 3 }
      ])
      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    it('passes for zero percentage', () => {
      const result = validateRewards([
        { level: 1, rewardPercentage: 0 }
      ])
      expect(result.valid).toBe(true)
    })
  })

  describe('buildSavePayload', () => {
    it('builds payload with rewards', () => {
      const rewards = [
        { level: 1, rewardPercentage: 5 },
        { level: 2, rewardPercentage: 3 }
      ]
      const payload = buildSavePayload(rewards)
      expect(payload).toEqual({
        rewards: [
          { level: 1, rewardPercentage: 5 },
          { level: 2, rewardPercentage: 3 }
        ]
      })
    })

    it('handles null input', () => {
      const payload = buildSavePayload(null)
      expect(payload.rewards).toEqual([])
    })

    it('converts string numbers to numbers', () => {
      const rewards = [
        { level: '1', rewardPercentage: '5.5' }
      ]
      const payload = buildSavePayload(rewards)
      expect(payload.rewards[0].level).toBe(1)
      expect(payload.rewards[0].rewardPercentage).toBe(5.5)
    })

    it('defaults invalid percentage to 0', () => {
      const rewards = [
        { level: 1, rewardPercentage: 'invalid' }
      ]
      const payload = buildSavePayload(rewards)
      expect(payload.rewards[0].rewardPercentage).toBe(0)
    })
  })

  describe('mergeRewardsFromServer', () => {
    it('merges server values into form structure', () => {
      const formRewards = [
        { level: 1, rewardPercentage: 0 },
        { level: 2, rewardPercentage: 0 }
      ]
      const serverRewards = [
        { level: 1, rewardPercentage: 5 },
        { level: 2, rewardPercentage: 3 }
      ]
      const merged = mergeRewardsFromServer(formRewards, serverRewards)
      expect(merged[0].rewardPercentage).toBe(5)
      expect(merged[1].rewardPercentage).toBe(3)
    })

    it('keeps 0 for missing server levels', () => {
      const formRewards = [
        { level: 1, rewardPercentage: 0 },
        { level: 2, rewardPercentage: 0 }
      ]
      const serverRewards = [
        { level: 1, rewardPercentage: 5 }
      ]
      const merged = mergeRewardsFromServer(formRewards, serverRewards)
      expect(merged[0].rewardPercentage).toBe(5)
      expect(merged[1].rewardPercentage).toBe(0)
    })

    it('handles null server rewards', () => {
      const formRewards = [
        { level: 1, rewardPercentage: 0 }
      ]
      const merged = mergeRewardsFromServer(formRewards, null)
      expect(merged[0].rewardPercentage).toBe(0)
    })

    it('handles non-array server rewards', () => {
      const formRewards = [
        { level: 1, rewardPercentage: 0 }
      ]
      const merged = mergeRewardsFromServer(formRewards, 'invalid')
      expect(merged[0].rewardPercentage).toBe(0)
    })
  })

  describe('calculateTotalReward', () => {
    it('sums all reward percentages', () => {
      const rewards = [
        { level: 1, rewardPercentage: 5 },
        { level: 2, rewardPercentage: 3 },
        { level: 3, rewardPercentage: 2 }
      ]
      expect(calculateTotalReward(rewards)).toBe(10)
    })

    it('returns 0 for empty array', () => {
      expect(calculateTotalReward([])).toBe(0)
    })

    it('returns 0 for null', () => {
      expect(calculateTotalReward(null)).toBe(0)
    })

    it('handles missing percentage (defaults to 0)', () => {
      const rewards = [
        { level: 1 },
        { level: 2, rewardPercentage: 5 }
      ]
      expect(calculateTotalReward(rewards)).toBe(5)
    })
  })

  describe('getLevelLabel', () => {
    it('returns "一级分销商" for level 1', () => {
      expect(getLevelLabel(1)).toBe('一级分销商')
    })

    it('returns "二级分销商" for level 2', () => {
      expect(getLevelLabel(2)).toBe('二级分销商')
    })

    it('returns "三级分销商" for level 3', () => {
      expect(getLevelLabel(3)).toBe('三级分销商')
    })

    it('returns generic label for unknown levels', () => {
      expect(getLevelLabel(4)).toBe('4级分销商')
    })
  })

  describe('LEVEL_LABELS constant', () => {
    it('has 3 levels defined', () => {
      expect(Object.keys(LEVEL_LABELS)).toHaveLength(3)
    })
  })
})
