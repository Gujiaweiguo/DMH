/**
 * FeedbackCenter business logic
 */

export const CATEGORY_OPTIONS = [
  { value: 'poster', label: '海报' },
  { value: 'payment', label: '支付' },
  { value: 'verification', label: '核销' },
  { value: 'other', label: '其他' }
]

export const PRIORITY_OPTIONS = [
  { value: 'medium', label: '中' },
  { value: 'high', label: '高' },
  { value: 'low', label: '低' }
]

export const RATING_OPTIONS = [
  { value: null, label: '不评分' },
  { value: 1, label: '1 分' },
  { value: 2, label: '2 分' },
  { value: 3, label: '3 分' },
  { value: 4, label: '4 分' },
  { value: 5, label: '5 分' }
]

export const STATUS_TEXT_MAP = {
  pending: '待处理',
  reviewing: '处理中',
  resolved: '已解决',
  closed: '已关闭'
}

export const getDefaultForm = () => ({
  category: 'poster',
  priority: 'medium',
  title: '',
  content: '',
  rating: null
})

export const statusText = (status) => STATUS_TEXT_MAP[status] || status

export const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

export const validateFeedbackForm = (form) => {
  const errors = []
  
  if (!form?.title || form.title.trim() === '') {
    errors.push('请填写标题')
  }
  
  if (!form?.content || form.content.trim() === '') {
    errors.push('请填写内容')
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

export const buildFeedbackPayload = (form) => ({
  category: form.category,
  priority: form.priority,
  title: form.title,
  content: form.content,
  rating: form.rating,
  featureUseCase: 'h5_feedback_center',
  deviceInfo: typeof navigator !== 'undefined' ? navigator.userAgent : '',
  browserInfo: typeof navigator !== 'undefined' ? navigator.userAgent : ''
})

export const resetForm = () => getDefaultForm()

export const getSubmitButtonText = (submitting) => {
  return submitting ? '提交中...' : '提交反馈'
}

export const getHelpfulCount = (faq, type) => {
  if (type === 'helpful') {
    return faq.helpfulCount || 0
  }
  return faq.notHelpfulCount || 0
}
