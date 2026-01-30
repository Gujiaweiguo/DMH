<template>
  <div class="verification-records">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <h1 class="nav-title">核销记录</h1>
    </div>

    <!-- 筛选器 -->
    <div class="filters">
      <div class="filter-tabs">
        <button
          v-for="status in statusTabs"
          :key="status.value"
          @click="currentStatus = status.value"
          :class="['filter-tab', { active: currentStatus === status.value }]"
        >
          {{ status.label }}
        </button>
      </div>

      <div class="date-filter">
        <input
          v-model="dateRange.start"
          type="date"
          class="date-input"
          placeholder="开始日期"
        >
        <span class="date-separator">至</span>
        <input
          v-model="dateRange.end"
          type="date"
          class="date-input"
          placeholder="结束日期"
        >
      </div>

      <div class="search-filter">
        <input
          v-model="searchKeyword"
          type="text"
          class="search-input"
          placeholder="搜索订单号或手机号"
        >
      </div>
    </div>

    <!-- 核销记录列表 -->
    <div class="records-list">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="filteredRecords.length === 0" class="empty-state">
        <div class="empty-icon">📋</div>
        <p class="empty-text">暂无核销记录</p>
      </div>

      <div v-else class="record-cards">
        <div
          v-for="record in filteredRecords"
          :key="record.orderId"
          class="record-card"
        >
          <div class="card-header">
            <div class="record-info">
              <h3 class="order-id">订单 #{{ record.orderId }}</h3>
              <span class="record-time">{{ formatDateTime(record.verifiedAt) }}</span>
            </div>
            <span :class="['status-badge', record.verificationStatus]">
              {{ getVerificationStatusText(record.verificationStatus) }}
            </span>
          </div>

          <div class="card-content">
            <div class="verification-details">
              <div class="detail-row">
                <span class="detail-label">订单状态:</span>
                <span class="detail-value">{{ record.orderStatus }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">订单金额:</span>
                <span class="detail-value amount">¥{{ record.orderAmount?.toFixed(2) || '0.00' }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">用户手机:</span>
                <span class="detail-value">{{ record.userPhone }}</span>
              </div>
              <div class="detail-row" v-if="record.verifiedBy">
                <span class="detail-label">核销人:</span>
                <span class="detail-value">{{ record.verifiedByName }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">核销时间:</span>
                <span class="detail-value">{{ formatDateTime(record.verifiedAt) }}</span>
              </div>
            </div>

            <div class="verification-code" v-if="record.verificationCode">
              <span class="code-label">核销码:</span>
              <span class="code-value">{{ record.verificationCode }}</span>
            </div>
          </div>

          <div class="card-actions">
            <button
              v-if="record.verificationStatus === 'verified'"
              @click="unverifyOrder(record.orderId)"
              class="action-btn unverify"
            >
              取消核销
            </button>
            <button
              @click="viewOrderDetail(record.orderId)"
              class="action-btn detail"
            >
              查看订单
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="stats-section">
      <h2 class="stats-title">核销统计</h2>
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-number">{{ recordStats.total }}</div>
          <div class="stat-label">总核销数</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ recordStats.today }}</div>
          <div class="stat-label">今日核销</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">¥{{ recordStats.totalAmount?.toFixed(2) || '0.00' }}</div>
          <div class="stat-label">核销总金额</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Toast } from 'vant'
import { orderApi } from '../../services/brandApi.js'

const router = useRouter()

const records = ref([])
const loading = ref(false)
const currentStatus = ref('all')
const searchKeyword = ref('')
const dateRange = ref({
  start: '',
  end: ''
})

const statusTabs = [
  { value: 'all', label: '全部' },
  { value: 'verified', label: '已核销' },
  { value: 'unverified', label: '未核销' },
  { value: 'cancelled', label: '已取消' }
]

const recordStats = ref({
  total: 0,
  today: 0,
  totalAmount: 0
})

const filteredRecords = computed(() => {
  let result = records.value

  if (currentStatus.value !== 'all') {
    result = result.filter(r => r.verificationStatus === currentStatus.value)
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(r =>
      r.orderId.toString().includes(keyword) ||
      r.userPhone?.includes(keyword)
    )
  }

  return result
})

const getVerificationStatusText = (status) => {
  const statusMap = {
    verified: '已核销',
    unverified: '未核销',
    cancelled: '已取消'
  }
  return statusMap[status] || status
}

const formatDateTime = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadRecords = async () => {
  loading.value = true
  try {
    // 调用真实API获取核销记录
    const resp = await orderApi.getVerificationRecords()
    if (resp && resp.records) {
      // 合并订单信息到核销记录
      records.value = await enrichRecordsWithOrderInfo(resp.records)
      calculateStats()
    }
  } catch (error) {
    console.error('加载核销记录失败:', error)
    Toast.fail('加载失败，请重试')
  } finally {
    loading.value = false
  }
}

const enrichRecordsWithOrderInfo = async (verificationRecords) => {
  // 获取所有订单ID
  const orderIds = verificationRecords.map(r => r.orderId)
  if (orderIds.length === 0) return verificationRecords

  try {
    // 批量获取订单信息
    const ordersResp = await orderApi.getOrders()
    const ordersMap = {}
    if (ordersResp && ordersResp.orders) {
      ordersResp.orders.forEach(order => {
        ordersMap[order.id] = order
      })
    }

    // 合并信息
    return verificationRecords.map(record => {
      const order = ordersMap[record.orderId]
      return {
        ...record,
        orderStatus: order?.status || '',
        orderAmount: order?.amount || 0,
        userPhone: order?.phone || '',
        verifiedByName: `用户${record.verifiedBy || '-'}`
      }
    })
  } catch (error) {
    console.error('获取订单信息失败:', error)
    // 即使失败，也返回基本的核销记录
    return verificationRecords
  }
}

const calculateStats = () => {
  const verifiedRecords = records.value.filter(r => r.verificationStatus === 'verified')
  const today = new Date().toLocaleDateString('zh-CN')
  
  recordStats.value = {
    total: verifiedRecords.length,
    today: verifiedRecords.filter(r => {
      const recordDate = new Date(r.verifiedAt).toLocaleDateString('zh-CN')
      return recordDate === today
    }).length,
    totalAmount: verifiedRecords.reduce((sum, r) => sum + (r.orderAmount || 0), 0)
  }
}

const viewOrderDetail = (orderId) => {
  router.push(`/brand/order-detail/${orderId}`)
}

const unverifyOrder = async (orderId) => {
  try {
    await orderApi.unverifyOrder('', { orderId })
    Toast.success('取消核销成功')
    await loadRecords()
  } catch (error) {
    console.error('取消核销失败:', error)
    Toast.fail('取消核销失败，请重试')
  }
}

onMounted(() => {
  loadRecords()
})
</script>

<style scoped>
.verification-records {
  min-height: 100vh;
  background: #f5f7fa;
}

.top-nav {
  background: white;
  padding: 16px;
  border-bottom: 1px solid #eee;
}

.nav-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.filters {
  background: white;
  padding: 16px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  border-bottom: 1px solid #eee;
}

.filter-tabs {
  display: flex;
  gap: 8px;
}

.filter-tab {
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background: white;
  color: #6b7280;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-tab.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.filter-tab:hover:not(.active) {
  background: #f3f4f6;
}

.date-filter {
  display: flex;
  align-items: center;
  gap: 8px;
}

.date-input {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  width: 150px;
}

.date-separator {
  color: #9ca3af;
}

.search-filter {
  flex: 1;
}

.search-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
}

.records-list {
  padding: 16px;
  max-width: 800px;
  margin: 0 auto;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #6b7280;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #9ca3af;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 14px;
}

.record-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.record-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f3f4f6;
}

.record-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.order-id {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.record-time {
  font-size: 12px;
  color: #9ca3af;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.verified {
  background: #dcfce7;
  color: #166534;
}

.status-badge.unverified {
  background: #fef3c7;
  color: #dc2626;
}

.status-badge.cancelled {
  background: #fee2e2;
  color: #991b1b;
}

.card-content {
  padding: 16px;
}

.verification-details {
  display: grid;
  gap: 12px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f9fafb;
}

.detail-label {
  font-size: 14px;
  color: #6b7280;
}

.detail-value {
  font-size: 14px;
  color: #1f2937;
  font-weight: 500;
}

.detail-value.amount {
  color: #07c160;
  font-weight: 600;
}

.verification-code {
  margin-top: 16px;
  padding: 12px;
  background: #f8fafc;
  border-radius: 6px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.code-label {
  font-size: 14px;
  color: #6b7280;
}

.code-value {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #1f2937;
}

.card-actions {
  display: flex;
  gap: 8px;
  padding: 16px;
  border-top: 1px solid #f3f4f6;
}

.action-btn {
  flex: 1;
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn.unverify {
  background: white;
  color: #dc2626;
  border: 1px solid #dc2626;
}

.action-btn.unverify:hover {
  background: #fef2f2;
}

.action-btn.detail {
  background: #667eea;
  color: white;
}

.action-btn.detail:hover {
  background: #5a67d8;
}

.stats-section {
  max-width: 800px;
  margin: 20px auto;
  padding: 20px;
  background: white;
  border-radius: 12px;
}

.stats-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 16px 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-item {
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  text-align: center;
}

.stat-number {
  font-size: 24px;
  font-weight: 600;
  color: #667eea;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
}
</style>
