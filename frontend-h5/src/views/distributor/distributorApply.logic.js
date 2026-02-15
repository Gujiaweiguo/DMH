/**
 * DistributorApply business logic
 */

export const DISTRIBUTOR_DESCRIPTION = '成为分销商，分享推广链接，获得订单佣金奖励'

export const getDefaultForm = () => ({
  brandId: null,
  reason: ''
})

export const buildBrandOptions = (brands) => {
  if (!Array.isArray(brands)) return []
  return brands.map(b => ({
    text: b.name,
    value: b.id
  }))
}

export const canSubmit = (form, agreeTerms) => {
  return !!(
    form &&
    form.brandId &&
    form.reason &&
    form.reason.trim() &&
    agreeTerms
  )
}

export const buildApplyPayload = (form) => ({
  brandId: form?.brandId,
  reason: form?.reason || ''
})

export const validateForm = (form, agreeTerms) => {
  const errors = []
  
  if (!form?.brandId) {
    errors.push('请选择品牌')
  }
  
  if (!form?.reason || form.reason.trim() === '') {
    errors.push('请填写申请理由')
  }
  
  if (!agreeTerms) {
    errors.push('请先同意分销商协议')
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

export const TERMS_CONTENT = {
  title: '分销商协议',
  sections: [
    {
      title: '一、分销商资格',
      items: [
        '1.1 分销商需经品牌管理员审批后方可获得分销资格。',
        '1.2 分销商需遵守平台规则和品牌相关规定。'
      ]
    },
    {
      title: '二、推广规则',
      items: [
        '2.1 分销商可通过专属推广链接和二维码进行推广。',
        '2.2 推广过程中不得进行虚假宣传或误导消费者。'
      ]
    },
    {
      title: '三、奖励说明',
      items: [
        '3.1 分销商可获得推广订单的佣金奖励。',
        '3.2 奖励比例根据分销商级别和品牌设置而定。',
        '3.3 最多支持三级分销，符合相关法规要求。'
      ]
    },
    {
      title: '四、违规处理',
      items: [
        '4.1 违反协议规定的分销商将被暂停或取消分销资格。',
        '4.2 涉嫌违法行为的将移交司法机关处理。'
      ]
    }
  ]
}
