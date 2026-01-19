<template>
  <div class="brand-distributors">
    <div class="top-nav">
      <div class="nav-header">
        <h1 class="nav-title">分销管理</h1>
        <div class="nav-actions">
          <button class="action-btn secondary" @click="goToLevelRewards">奖励设置</button>
          <button class="action-btn primary" @click="goToApproval">
            去审批
            <span v-if="pendingCount > 0" class="count-badge">{{ pendingCount }}</span>
          </button>
        </div>
      </div>
      <div class="nav-subtitle">
        <span>总计 {{ total }} 人</span>
        <span class="dot">·</span>
        <span>启用 {{ activeCount }} 人</span>
        <span class="dot">·</span>
        <span>停用 {{ suspendedCount }} 人</span>
      </div>
    </div>

    <div class="filters">
      <div class="filter-row">
        <select v-model="filters.status" class="filter-select" @change="refresh">
          <option value="">全部状态</option>
          <option value="active">启用</option>
          <option value="suspended">停用</option>
        </select>
        <select v-model.number="filters.level" class="filter-select" @change="refresh">
          <option :value="0">全部级别</option>
          <option :value="1">一级</option>
          <option :value="2">二级</option>
          <option :value="3">三级</option>
        </select>
      </div>

      <div class="filter-row">
        <input
          v-model="filters.keyword"
          class="search-input"
          placeholder="搜索用户名/手机号"
          @keyup.enter="refresh"
        >
        <button class="search-btn" @click="refresh">查询</button>
      </div>
    </div>

    <div class="list">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="distributors.length === 0" class="empty-state">
        <div class="empty-icon">🧾</div>
        <p class="empty-text">暂无分销数据</p>
      </div>

      <div v-else class="cards">
        <div v-for="d in distributors" :key="d.id" class="card">
          <div class="card-header">
            <div class="user">
              <div class="avatar">{{ (d.username || '用').charAt(0) }}</div>
              <div class="info">
                <div class="name-line">
                  <span class="name">{{ d.username || '-' }}</span>
                  <span :class="['status-badge', d.status]">{{ getStatusText(d.status) }}</span>
                </div>
                <div class="meta">
                  <span class="pill">级别 {{ d.level || 1 }}</span>
                  <span class="pill">下级 {{ d.subordinatesCount || 0 }}</span>
                </div>
              </div>
            </div>
            <div class="money">
              <div class="money-value">¥{{ (d.totalEarnings || 0).toFixed(2) }}</div>
              <div class="money-label">累计收益</div>
            </div>
          </div>

          <div class="card-body">
            <div class="row">
              <span class="label">所属品牌</span>
              <span class="value">{{ d.brandName || '-' }}</span>
            </div>
            <div class="row">
              <span class="label">上级</span>
              <span class="value">{{ d.parentName || '-' }}</span>
            </div>
            <div class="row">
              <span class="label">加入时间</span>
              <span class="value">{{ formatTime(d.createdAt) }}</span>
            </div>
          </div>

          <div class="card-actions">
            <button class="mini-btn" @click="openLevelModal(d)">调级</button>
            <button
              class="mini-btn danger"
              v-if="d.status === 'active'"
              @click="openStatusModal(d, 'suspended')"
            >
              停用
            </button>
            <button
              class="mini-btn"
              v-else
              @click="openStatusModal(d, 'active')"
            >
              启用
            </button>
          </div>
        </div>
      </div>

      <div v-if="hasMore && !loading" class="load-more">
        <button class="load-more-btn" @click="loadMore">加载更多</button>
      </div>
    </div>

    <div v-if="showLevelModal" class="modal-overlay" @click="closeModals">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>调整分销级别</h3>
          <button class="close-btn" @click="closeModals">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>分销商</label>
            <div class="form-static">{{ currentDistributor?.username || '-' }}</div>
          </div>
          <div class="form-group">
            <label>级别</label>
            <select v-model.number="levelForm.level" class="form-select">
              <option :value="1">一级</option>
              <option :value="2">二级</option>
              <option :value="3">三级</option>
            </select>
          </div>
        </div>
        <div class="modal-actions">
          <button class="action-btn secondary" @click="closeModals">取消</button>
          <button class="action-btn primary" @click="submitLevel" :disabled="submitting">保存</button>
        </div>
      </div>
    </div>

    <div v-if="showStatusModal" class="modal-overlay" @click="closeModals">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ statusForm.status === 'active' ? '启用分销商' : '停用分销商' }}</h3>
          <button class="close-btn" @click="closeModals">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>分销商</label>
            <div class="form-static">{{ currentDistributor?.username || '-' }}</div>
          </div>
          <div class="form-group">
            <label>原因（可选）</label>
            <textarea
              v-model="statusForm.reason"
              class="form-textarea"
              rows="3"
              placeholder="例如：违规、冻结、恢复等"
            ></textarea>
          </div>
        </div>
        <div class="modal-actions">
          <button class="action-btn secondary" @click="closeModals">取消</button>
          <button class="action-btn primary" @click="submitStatus" :disabled="submitting">确认</button>
        </div>
      </div>
    </div>

    <div class="bottom-nav">
      <router-link to="/brand/dashboard" class="nav-item">
        <div class="nav-icon">🏠</div>
        <div class="nav-text">工作台</div>
      </router-link>
      <router-link to="/brand/campaigns" class="nav-item">
        <div class="nav-icon">🎯</div>
        <div class="nav-text">活动</div>
      </router-link>
      <router-link to="/brand/orders" class="nav-item">
        <div class="nav-icon">📋</div>
        <div class="nav-text">订单</div>
      </router-link>
      <router-link to="/brand/distributors" class="nav-item active">
        <div class="nav-icon">🧭</div>
        <div class="nav-text">分销</div>
      </router-link>
      <router-link to="/brand/promoters" class="nav-item">
        <div class="nav-icon">👥</div>
        <div class="nav-text">推广员</div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const distributors = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const hasMore = ref(false)
const pendingCount = ref(0)

const filters = reactive({
  keyword: '',
  status: '',
  level: 0
})

const showLevelModal = ref(false)
const showStatusModal = ref(false)
const currentDistributor = ref(null)

const levelForm = reactive({
  level: 1
})

const statusForm = reactive({
  status: 'active',
  reason: ''
})

const getCurrentBrandId = () => {
  const fromStorage = Number(localStorage.getItem('dmh_current_brand_id'))
  if (Number.isFinite(fromStorage) && fromStorage > 0) return fromStorage
  try {
    const info = JSON.parse(localStorage.getItem('dmh_user_info') || '{}')
    const firstBrandId = Array.isArray(info.brandIds) && info.brandIds.length > 0 ? Number(info.brandIds[0]) : 0
    if (Number.isFinite(firstBrandId) && firstBrandId > 0) {
      localStorage.setItem('dmh_current_brand_id', String(firstBrandId))
      return firstBrandId
    }
  } catch {
    // ignore
  }
  return 0
}

const formatTime = (timeString) => {
  if (!timeString) return '-'
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const getStatusText = (status) => {
  const map = {
    active: '启用',
    suspended: '停用'
  }
  return map[status] || status || '-'
}

const activeCount = computed(() => distributors.value.filter(d => d.status === 'active').length)
const suspendedCount = computed(() => distributors.value.filter(d => d.status === 'suspended').length)

const goToApproval = () => router.push('/brand/distributor-approval')
const goToLevelRewards = () => router.push('/brand/distributor-level-rewards')

const refresh = async () => {
  page.value = 1
  await loadList(true)
}

const loadMore = async () => {
  page.value += 1
  await loadList(false)
}

const loadPendingCount = async () => {
  const token = localStorage.getItem('dmh_token')
  const brandId = getCurrentBrandId()
  if (!brandId) return

  try {
    const response = await fetch(`/api/v1/brands/${brandId}/distributor/applications?page=1&pageSize=1&status=pending`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (!response.ok) return
    const data = await response.json()
    pendingCount.value = Number(data.total) || 0
  } catch {
    // ignore
  }
}

const loadList = async (replace) => {
  loading.value = true
  try {
    const token = localStorage.getItem('dmh_token')
    const brandId = getCurrentBrandId()
    if (!brandId) {
      alert('未选择品牌，请重新登录')
      router.push('/brand/login')
      return
    }

    const query = new URLSearchParams()
    query.set('page', String(page.value))
    query.set('pageSize', String(pageSize))
    query.set('brandId', String(brandId))
    if (filters.keyword.trim()) query.set('keyword', filters.keyword.trim())
    if (filters.status) query.set('status', filters.status)
    if (filters.level > 0) query.set('level', String(filters.level))

    const response = await fetch(`/api/v1/brands/${brandId}/distributors?${query.toString()}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      const data = await response.json()
      const list = Array.isArray(data.distributors) ? data.distributors : []
      total.value = Number(data.total) || 0
      distributors.value = replace ? list : distributors.value.concat(list)
      hasMore.value = distributors.value.length < total.value
      return
    }

    if (replace) {
      distributors.value = []
      total.value = 0
      hasMore.value = false
    }
  } catch (error) {
    console.error('加载分销列表失败:', error)
  } finally {
    loading.value = false
  }
}

const openLevelModal = (d) => {
  currentDistributor.value = d
  levelForm.level = d.level || 1
  showLevelModal.value = true
}

const openStatusModal = (d, nextStatus) => {
  currentDistributor.value = d
  statusForm.status = nextStatus
  statusForm.reason = ''
  showStatusModal.value = true
}

const closeModals = () => {
  showLevelModal.value = false
  showStatusModal.value = false
  currentDistributor.value = null
}

const submitLevel = async () => {
  if (!currentDistributor.value) return
  submitting.value = true
  try {
    const token = localStorage.getItem('dmh_token')
    const response = await fetch(`/api/v1/brands/distributors/${currentDistributor.value.id}/level`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ level: levelForm.level })
    })
    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      alert(`保存失败：${data.message || '未知错误'}`)
      return
    }
    closeModals()
    await refresh()
  } catch (error) {
    console.error('更新级别失败:', error)
    alert('保存失败，请重试')
  } finally {
    submitting.value = false
  }
}

const submitStatus = async () => {
  if (!currentDistributor.value) return
  submitting.value = true
  try {
    const token = localStorage.getItem('dmh_token')
    const response = await fetch(`/api/v1/brands/distributors/${currentDistributor.value.id}/status`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        status: statusForm.status,
        reason: statusForm.reason
      })
    })
    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      alert(`保存失败：${data.message || '未知错误'}`)
      return
    }
    closeModals()
    await refresh()
  } catch (error) {
    console.error('更新状态失败:', error)
    alert('保存失败，请重试')
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  await loadPendingCount()
  await loadList(true)
})
</script>

<style scoped>
.brand-distributors {
  min-height: 100vh;
  background: #f5f7fa;
  padding-bottom: 80px;
}

.top-nav {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 16px;
}

.nav-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.nav-title {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
}

.nav-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.nav-subtitle {
  margin-top: 10px;
  font-size: 12px;
  opacity: 0.95;
}

.dot {
  margin: 0 6px;
  opacity: 0.8;
}

.action-btn {
  border: none;
  border-radius: 10px;
  padding: 10px 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.action-btn.primary {
  background: rgba(255, 255, 255, 0.95);
  color: #5b4bb7;
}

.action-btn.secondary {
  background: rgba(255, 255, 255, 0.18);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.25);
}

.count-badge {
  background: #ff4757;
  color: white;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 999px;
  line-height: 18px;
}

.filters {
  padding: 12px 16px 6px;
}

.filter-row {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.filter-select {
  flex: 1;
  padding: 12px;
  border-radius: 12px;
  border: 1px solid #e6ebf2;
  background: white;
  font-size: 14px;
}

.search-input {
  flex: 1;
  padding: 12px;
  border-radius: 12px;
  border: 1px solid #e6ebf2;
  background: white;
  font-size: 14px;
}

.search-btn {
  padding: 0 14px;
  border-radius: 12px;
  border: none;
  background: #667eea;
  color: white;
  font-weight: 600;
  cursor: pointer;
}

.list {
  padding: 0 16px 16px;
}

.loading {
  text-align: center;
  padding: 30px 0;
  color: #666;
}

.loading-spinner {
  width: 28px;
  height: 28px;
  border: 3px solid #e6ebf2;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 10px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: #666;
}

.empty-icon {
  font-size: 40px;
  margin-bottom: 10px;
}

.cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card {
  background: white;
  border-radius: 16px;
  padding: 14px;
  border: 1px solid #edf1f7;
}

.card-header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.user {
  display: flex;
  gap: 12px;
  align-items: center;
  min-width: 0;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
}

.info {
  min-width: 0;
}

.name-line {
  display: flex;
  gap: 8px;
  align-items: center;
}

.name {
  font-size: 16px;
  font-weight: 700;
  color: #333;
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 999px;
  background: #f1f3f5;
  color: #555;
}

.status-badge.active {
  background: #e8fff1;
  color: #1e9d53;
}

.status-badge.suspended {
  background: #fff1f0;
  color: #d4380d;
}

.meta {
  margin-top: 6px;
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.pill {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 999px;
  background: #f5f7fa;
  color: #555;
}

.money {
  text-align: right;
}

.money-value {
  font-size: 16px;
  font-weight: 800;
  color: #333;
}

.money-label {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.card-body {
  margin-top: 12px;
  border-top: 1px solid #f1f3f5;
  padding-top: 12px;
}

.row {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 13px;
}

.label {
  color: #888;
}

.value {
  color: #333;
  font-weight: 600;
}

.card-actions {
  margin-top: 10px;
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.mini-btn {
  border: 1px solid #e6ebf2;
  background: white;
  padding: 8px 12px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
}

.mini-btn.danger {
  border-color: #ffd7d7;
  color: #d4380d;
}

.load-more {
  text-align: center;
  padding: 14px 0 6px;
}

.load-more-btn {
  border: none;
  background: #f1f3f5;
  color: #333;
  padding: 10px 14px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
}

.modal-content {
  width: 100%;
  max-width: 420px;
  background: white;
  border-radius: 16px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  background: #f8f9fa;
}

.modal-header h3 {
  margin: 0;
  font-size: 16px;
}

.close-btn {
  border: none;
  background: transparent;
  font-size: 18px;
  cursor: pointer;
}

.modal-body {
  padding: 16px;
}

.form-group {
  margin-bottom: 14px;
}

.form-group label {
  display: block;
  font-size: 13px;
  color: #666;
  margin-bottom: 8px;
}

.form-static {
  padding: 12px;
  border-radius: 12px;
  background: #f5f7fa;
  font-weight: 700;
  color: #333;
}

.form-select,
.form-textarea {
  width: 100%;
  padding: 12px;
  border-radius: 12px;
  border: 1px solid #e6ebf2;
  font-size: 14px;
  box-sizing: border-box;
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  padding: 12px 16px 16px;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 70px;
  background: white;
  border-top: 1px solid #e6ebf2;
  display: flex;
  justify-content: space-around;
  align-items: center;
  z-index: 999;
}

.nav-item {
  text-decoration: none;
  color: #666;
  text-align: center;
  font-size: 12px;
}

.nav-item.active {
  color: #667eea;
}

.nav-icon {
  font-size: 20px;
  margin-bottom: 4px;
}
</style>

