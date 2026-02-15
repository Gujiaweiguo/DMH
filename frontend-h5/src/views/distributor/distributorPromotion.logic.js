/**
 * DistributorPromotion business logic
 */

export const PROMOTION_GUIDE = [
  '1. 选择要推广的活动，生成专属海报或链接',
  '2. 将海报或二维码分享给朋友，引导他们下单',
  '3. 用户通过您的推广下单，您即可获得佣金奖励',
  '4. 最多支持三级分销，您的下级推广也能为您带来收益'
]

export const buildCampaignOptions = (campaigns) => {
  if (!Array.isArray(campaigns)) return []
  return campaigns.map(c => ({
    text: c.name,
    value: c.id
  }))
}

export const buildPosterDownloadName = (baseName) => `${baseName}_${Date.now()}.png`

export const buildQrcodeDownloadName = (linkCode) => `推广二维码_${linkCode}.png`

export const buildQrcodeUrl = (link, qrcodeUrl) => {
  if (qrcodeUrl) return qrcodeUrl
  return `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(link)}`
}

export const canGenerateCampaignPoster = (campaignId) => !!campaignId

export const buildPosterPayload = (type, campaignId = null) => {
  const payload = { type }
  if (campaignId) {
    payload.campaignId = campaignId
  }
  return payload
}

export const buildLinkPayload = (campaignId) => ({ campaignId })

export const formatClickCount = (count) => `点击 ${count} 次`

export const getCampaignDisplayName = (link, campaigns) => {
  if (link.campaignName) return link.campaignName
  const campaign = campaigns.find(c => c.id === link.campaignId)
  return campaign ? campaign.name : `活动 ${link.campaignId}`
}

export const DEFAULT_POSTER_DATA = {
  campaignName: '',
  campaignDesc: '',
  distributorName: '',
  campaignCount: 0
}

export const DEFAULT_LINK_DATA = {
  link: '',
  linkCode: '',
  qrcodeUrl: ''
}
