/**
 * DistributorLevelRewards business logic
 * Extracted from DistributorLevelRewards.vue for testability
 */

/**
 * Get default rewards form with 3 levels
 * @returns {Array} Default rewards array
 */
export const getDefaultRewards = () => [
  { level: 1, rewardPercentage: 0 },
  { level: 2, rewardPercentage: 0 },
  { level: 3, rewardPercentage: 0 }
]

/**
 * Clamp reward percentage to valid range
 * @param {number} value - Raw percentage value
 * @returns {number} Clamped value between 0 and 100
 */
export const clampRewardPercentage = (value) => {
  const num = Number(value)
  if (!Number.isFinite(num)) return 0
  return Math.max(0, Math.min(100, num))
}

/**
 * Validate rewards form data
 * @param {Array} rewards - Rewards array
 * @returns {Object} Validation result { valid: boolean, errors: string[] }
 */
export const validateRewards = (rewards) => {
  const errors = []
  
  if (!Array.isArray(rewards) || rewards.length === 0) {
    errors.push('奖励配置不能为空')
    return { valid: false, errors }
  }

  for (const reward of rewards) {
    if (!Number.isFinite(reward.level) || reward.level < 1) {
      errors.push(`级别必须大于0`)
    }
    
    const percentage = Number(reward.rewardPercentage)
    if (!Number.isFinite(percentage)) {
      errors.push(`级别 ${reward.level} 的奖励比例无效`)
    } else if (percentage < 0 || percentage > 100) {
      errors.push(`级别 ${reward.level} 的奖励比例必须在 0-100 之间`)
    }
  }

  return { valid: errors.length === 0, errors }
}

/**
 * Build save payload from form data
 * @param {Array} rewards - Rewards array from form
 * @returns {Object} API payload
 */
export const buildSavePayload = (rewards) => ({
  rewards: (rewards || []).map(r => ({
    level: Number(r.level),
    rewardPercentage: Number(r.rewardPercentage) || 0
  }))
})

/**
 * Merge server rewards into form structure
 * @param {Array} formRewards - Current form rewards (with reactive refs)
 * @param {Array} serverRewards - Rewards from server
 * @returns {Array} Merged rewards with updated values
 */
export const mergeRewardsFromServer = (formRewards, serverRewards) => {
  const rewards = Array.isArray(serverRewards) ? serverRewards : []
  
  return formRewards.map(row => {
    const hit = rewards.find(r => Number(r.level) === Number(row.level))
    return {
      ...row,
      rewardPercentage: hit ? Number(hit.rewardPercentage) || 0 : 0
    }
  })
}

/**
 * Calculate total reward percentage
 * @param {Array} rewards - Rewards array
 * @returns {number} Sum of all reward percentages
 */
export const calculateTotalReward = (rewards) => {
  if (!Array.isArray(rewards)) return 0
  return rewards.reduce((sum, r) => sum + (Number(r.rewardPercentage) || 0), 0)
}

/**
 * Level label mapping
 */
export const LEVEL_LABELS = {
  1: '一级分销商',
  2: '二级分销商',
  3: '三级分销商'
}

/**
 * Get level display label
 * @param {number} level - Level number
 * @returns {string} Level label
 */
export const getLevelLabel = (level) => {
  return LEVEL_LABELS[level] || `${level}级分销商`
}
