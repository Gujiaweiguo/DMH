import { api } from './api.js'

// 认证相关API
export const authApi = {
  // 品牌管理员登录
  login: (username, password) => {
    return api.post('/auth/login', { username, password })
  },

  // 获取当前用户信息
  getUserInfo: () => {
    return api.get('/auth/userinfo')
  },

  // 登出
  logout: () => {
    return api.post('/auth/logout')
  }
}

// 活动管理API
export const campaignApi = {
  // 获取活动列表
  getCampaigns: (params = {}) => {
    return api.get('/campaigns', params)
  },

  // 获取活动详情
  getCampaign: (id) => {
    return api.get(`/campaign/detail/${id}`)
  },

  // 创建活动
  createCampaign: (data) => {
    return api.post('/campaign/create', data)
  },

  // 更新活动
  updateCampaign: (id, data) => {
    return api.put(`/campaign/update/${id}`, data)
  },

  // 删除活动
  deleteCampaign: (id) => {
    return api.delete(`/campaign/delete/${id}`)
  },

  // 保存页面配置
  savePageConfig: (campaignId, config) => {
    return api.post(`/campaign/page-config/${campaignId}`, config)
  },

  // 获取页面配置
  getPageConfig: (campaignId) => {
    return api.get(`/campaign/page-config/${campaignId}`)
  },

  // 获取支付二维码
  getPaymentQrcode: (campaignId) => {
    return api.get(`/campaigns/${campaignId}/payment-qrcode`)
  }
}

// 订单管理API
export const orderApi = {
  // 获取订单列表
  getOrders: (params = {}) => {
    return api.get('/order/list', params)
  },

  // 获取订单详情
  getOrder: (id) => {
    return api.get(`/order/detail/${id}`)
  },

  // 更新订单状态
  updateOrderStatus: (id, status) => {
    return api.put(`/order/status/${id}`, { status })
  },

  scanOrderCode: (code) => {
    return api.get('/orders/scan', { code })
  },

  verifyOrder: (code, data) => {
    return api.post('/orders/verify', { code, ...data })
  },

  unverifyOrder: (code, data) => {
    return api.post('/orders/unverify', { code, ...data })
  },

  getVerificationRecords: () => {
    return api.get('/order/verification-records')
  }
}

// 推广员管理API
export const promoterApi = {
  // 获取推广员列表
  getPromoters: (params = {}) => {
    return api.get('/promoter/list', params)
  },

  // 获取推广员详情
  getPromoter: (id) => {
    return api.get(`/promoter/detail/${id}`)
  },

  // 生成推广链接
  generateLink: (promoterId, campaignId) => {
    return api.post('/promoter/generate-link', { promoterId, campaignId })
  },

  // 获取推广员奖励记录
  getRewards: (promoterId, params = {}) => {
    return api.get(`/promoter/rewards/${promoterId}`, params)
  }
}

// 数据分析API
export const analyticsApi = {
  // 获取核心指标
  getMetrics: (period = 'month') => {
    return api.get('/analytics/metrics', { period })
  },

  // 获取趋势数据
  getTrends: (period = 'month') => {
    return api.get('/analytics/trends', { period })
  },

  // 获取活动排行
  getCampaignRanking: (period = 'month') => {
    return api.get('/analytics/campaign-ranking', { period })
  },

  // 获取推广员排行
  getPromoterRanking: (period = 'month') => {
    return api.get('/analytics/promoter-ranking', { period })
  }
}

// 海报生成相关API
export const posterApi = {
  // 获取海报模板列表
  getPosterTemplates: () => {
    return api.get('/poster/templates')
  },

  // 生成海报（活动海报或分销商海报）
  generatePoster: (campaignId, distributorId) => {
    return api.post(`/campaigns/${campaignId}/poster`, {
      distributorId: distributorId
    })
  },

  // 获取海报文件
  getPosterFile: (filename) => {
    return api.get(`/poster/${filename}`)
  },

  // 获取海报生成记录
  getPosterRecords: () => {
    return api.get('/poster/records')
  }
}

// 素材管理API
export const materialApi = {
  // 上传素材
  uploadMaterial: (file, type = 'image') => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', type)
    return api.post('/material/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // 获取素材列表
  getMaterials: (params = {}) => {
    return api.get('/material/list', params)
  },

  // 删除素材
  deleteMaterial: (id) => {
    return api.delete(`/material/delete/${id}`)
  }
}

// 品牌设置API
export const settingsApi = {
  // 获取品牌信息
  getBrandInfo: () => {
    return api.get('/brand/info')
  },

  // 更新品牌信息
  updateBrandInfo: (data) => {
    return api.put('/brand/info', data)
  },

  // 获取奖励设置
  getRewardSettings: () => {
    return api.get('/brand/reward-settings')
  },

  // 更新奖励设置
  updateRewardSettings: (data) => {
    return api.put('/brand/reward-settings', data)
  },

  // 获取通知设置
  getNotificationSettings: () => {
    return api.get('/brand/notification-settings')
  },

  // 更新通知设置
  updateNotificationSettings: (data) => {
    return api.put('/brand/notification-settings', data)
  }
}
