<template>
  <div class="container">
    <div class="header">
      <button class="back-btn" @click="goBack">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
      </button>
      <h1 class="title">我的报名</h1>
      <div style="width: 40px;"></div>
    </div>

    <!-- 手机号输入 -->
    <div class="search-box">
      <input 
        v-model="phone" 
        type="tel" 
        placeholder="请输入您的手机号" 
        class="phone-input"
        maxlength="11"
      />
      <button @click="searchOrders" class="search-btn" :disabled="!isValidPhone">
        查询
      </button>
    </div>

    <!-- 提示信息 -->
    <div v-if="!searched" class="tips">
      <p>💡 输入您报名时使用的手机号，即可查看所有报名记录</p>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="loading">加载中...</div>

    <!-- 订单列表 -->
    <div v-else-if="searched && orders.length > 0" class="orders-list">
      <div v-for="order in orders" :key="order.id" class="order-card">
        <div class="order-header">
          <span class="order-id">#{{ order.id }}</span>
          <span class="order-status" :class="order.status">
            {{ statusText(order.status) }}
          </span>
        </div>
        
        <h3 class="campaign-name">{{ getCampaignName(order.campaignId) }}</h3>
        
        <div class="order-info">
          <div class="info-row">
            <span class="label">手机号：</span>
            <span class="value">{{ order.phone }}</span>
          </div>
          <div class="info-row">
            <span class="label">报名时间：</span>
            <span class="value">{{ formatDate(order.createdAt) }}</span>
          </div>
          <div class="info-row">
            <span class="label">奖励金额：</span>
            <span class="value reward">¥{{ order.amount.toFixed(2) }}</span>
          </div>
        </div>

        <!-- 表单数据 -->
        <div v-if="order.formData && Object.keys(order.formData).length > 0" class="form-data">
          <div class="form-title">报名信息</div>
          <div class="form-content">
            <div v-for="(value, key) in order.formData" :key="key" class="form-item">
              <span class="form-label">{{ key }}：</span>
              <span class="form-value">{{ value }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="searched && orders.length === 0" class="empty">
      <div class="empty-icon">📝</div>
      <p>未找到报名记录</p>
      <p class="empty-tip">请确认手机号是否正确</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const phone = ref('')
const orders = ref([])
const campaigns = ref({})
const loading = ref(false)
const searched = ref(false)

// 验证手机号
const isValidPhone = computed(() => {
  return /^1[3-9]\d{9}$/.test(phone.value)
})

// 返回
const goBack = () => {
  router.back()
}

// 查询订单
const searchOrders = async () => {
  if (!isValidPhone.value) {
    alert('请输入正确的手机号')
    return
  }

  loading.value = true
  searched.value = true

  try {
    // 保存手机号到本地存储
    localStorage.setItem('dmh_my_phone', phone.value)
    
    // 先加载活动列表
    await loadCampaigns()

    // 查询订单
    const response = await fetch(`/api/v1/orders?pageSize=100`)
    if (response.ok) {
      const data = await response.json()
      // 过滤出当前手机号的订单
      orders.value = (data.orders || []).filter(order => order.phone === phone.value)
    } else {
      alert('查询失败，请稍后重试')
    }
  } catch (error) {
    console.error('查询订单失败', error)
    alert('查询失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 加载活动列表
const loadCampaigns = async () => {
  try {
    const response = await fetch('/api/v1/campaigns')
    if (response.ok) {
      const data = await response.json()
      const map = {}
      ;(data.campaigns || []).forEach(c => {
        map[c.id] = c
      })
      campaigns.value = map
    }
  } catch (error) {
    console.error('加载活动列表失败', error)
  }
}

// 获取活动名称
const getCampaignName = (campaignId) => {
  return campaigns.value[campaignId]?.name || '未知活动'
}

// 格式化日期
const formatDate = (time) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 状态文本
const statusText = (status) => {
  const map = {
    paid: '已支付',
    pending: '待支付',
    cancelled: '已取消'
  }
  return map[status] || status
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 20px;
}

.header {
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  padding: 16px;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.back-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.title {
  font-size: 18px;
  font-weight: 600;
}

.search-box {
  padding: 16px;
  background-color: #fff;
  display: flex;
  gap: 12px;
}

.phone-input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 16px;
  outline: none;
}

.phone-input:focus {
  border-color: #4f46e5;
}

.search-btn {
  padding: 12px 24px;
  background-color: #4f46e5;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
}

.search-btn:disabled {
  background-color: #d1d5db;
  cursor: not-allowed;
}

.search-btn:active:not(:disabled) {
  background-color: #4338ca;
}

.tips {
  padding: 16px;
  background-color: #fef3c7;
  margin: 0 16px 16px;
  border-radius: 8px;
}

.tips p {
  font-size: 14px;
  color: #92400e;
  margin: 0;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: #999;
  font-size: 14px;
}

.orders-list {
  padding: 0 16px;
}

.order-card {
  background-color: #fff;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.order-id {
  font-size: 12px;
  color: #999;
  font-family: monospace;
}

.order-status {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.order-status.paid {
  background-color: #dcfce7;
  color: #16a34a;
}

.order-status.pending {
  background-color: #fef3c7;
  color: #d97706;
}

.order-status.cancelled {
  background-color: #fee2e2;
  color: #dc2626;
}

.campaign-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.order-info {
  padding: 12px;
  background-color: #f9fafb;
  border-radius: 8px;
  margin-bottom: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  margin-bottom: 8px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-row .label {
  color: #6b7280;
}

.info-row .value {
  color: #111827;
  font-weight: 500;
}

.info-row .value.reward {
  color: #16a34a;
  font-size: 16px;
  font-weight: 600;
}

.form-data {
  border-top: 1px solid #f3f4f6;
  padding-top: 12px;
}

.form-title {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 8px;
  font-weight: 500;
}

.form-content {
  background-color: #f9fafb;
  border-radius: 8px;
  padding: 12px;
}

.form-item {
  font-size: 14px;
  margin-bottom: 6px;
  display: flex;
}

.form-item:last-child {
  margin-bottom: 0;
}

.form-label {
  color: #6b7280;
  min-width: 80px;
}

.form-value {
  color: #111827;
  flex: 1;
  word-break: break-all;
}

.empty {
  text-align: center;
  padding: 80px 20px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty p {
  color: #6b7280;
  font-size: 16px;
  margin-bottom: 8px;
}

.empty-tip {
  color: #9ca3af;
  font-size: 14px;
}
</style>
