/**
 * CampaignPageDesigner business logic
 */

export const COMPONENT_TYPES = {
  banner: 'BannerComponent',
  text: 'TitleComponent',
  video: 'VideoComponent',
  countdown: 'TimeComponent',
  faq: 'FaqComponent',
  contact: 'ContactComponent',
  social: 'SocialComponent'
}

export const COMPONENT_LIBRARY = [
  { name: '横幅图片', type: 'banner' },
  { name: '文本标题', type: 'text' },
  { name: '倒计时', type: 'countdown' },
  { name: 'FAQ', type: 'faq' },
  { name: '联系方式', type: 'contact' }
]

export const getComponentType = (type) => COMPONENT_TYPES[type] || type

export const canMoveUp = (index) => index > 0

export const canMoveDown = (index, length) => index < length - 1

export const moveUp = (components, index) => {
  if (!canMoveUp(index)) return components
  const result = [...components]
  ;[result[index - 1], result[index]] = [result[index], result[index - 1]]
  return result
}

export const moveDown = (components, index) => {
  if (!canMoveDown(index, components.length)) return components
  const result = [...components]
  ;[result[index], result[index + 1]] = [result[index + 1], result[index]]
  return result
}

export const removeComponent = (components, id) => {
  return components.filter(c => c.id !== id)
}

export const addComponent = (components, component) => {
  return [...components, component]
}

export const generateComponentId = () => `comp_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`

export const createComponent = (type, data = {}) => ({
  id: generateComponentId(),
  type,
  data,
  style: {}
})

export const hasComponents = (components) => !!(components && components.length > 0)

export const getComponentCount = (components) => components ? components.length : 0
