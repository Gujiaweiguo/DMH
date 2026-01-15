<template>
  <div class="brand-dashboard">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-content">
        <div class="brand-info">
          <img :src="brandInfo.logo" alt="品牌logo" class="brand-logo">
          <div>
            <h1 class="brand-name">{{ brandInfo.name }}</h1>
            <p class="welcome-text">品牌工作台</p>
          </div>
        </div>
        <div class="nav-actions">
          <button @click="logout" class="logout-btn">退出</button>
        </div>
      </div>
    </div>

    <!-- 数据概览 -->
    <div class="stats-section">
      <h2 class="section-title">今日数据</h2>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-number">{{ todayStats.orders }}</div>
          <div class="stat-label">新增订单</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">¥{{ todayStats.rewards }}</div>
          <div class="stat-label">奖励发放</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ todayStats.promoters }}</div>
          <div class="stat-label">新增推广员</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ todayStats.campaigns }}</div>
          <div class="stat-label">活跃活动</div>
        </div>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="quick-actions">
      <h2 class="section-title">快捷操作</h2>
      <div class="action-grid">
        <router-link to="/brand/campaigns/create" class="action-card">
          <div class="action-icon">📝</div>
          <div class="action-text">创建活动</div>
        </router-link>
        <router-link to="/brand/materials" class="action-card">
          <div class="action-icon">🎨</div>
          <div class="action-text">素材库</div>
        </router-link>
        <router-link to="/brand/orders" class="action-card">
          <div class="action-icon">📊</div>
          <div class="action-text">订单管理</div>
        </router-link>
        <router-link to="/brand/analytics" class="action-card">
          <div class="action-icon">📈</div>
          <div class="action-text">数据分析</div>
        </router-link>
      </div>
    </div>

    <!-- 活动状态 -->
    <div class="campaigns-section">
      <div class="section-header">
        <h2 class="section-title">我的活动</h2>
        <router-link to="/brand/campaigns" class="view-all">查看全部</router-link>
      </div>
      <div class="campaign-list">
        <div v-for="campaign in recentCampaigns" :key="campaign.id" class="campaign-card">
          <div class="campaign-info">
            <h3 class="campaign-name">{{ campaign.name }}</h3>
            <p class="campaign-desc">{{ campaign.description }}</p>
            <div class="campaign-stats">
              <span class="stat">{{ campaign.orders }}人参与</span>
              <span class="stat">¥{{ campaign.rewards }}奖励</span>
            </div>
          </div>
          <div class="campaign-status">
            <span :class="['status-badge', campaign.status]">
              {{ getStatusText(campaign.status) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部导航 -->
    <div class="bottom-nav">
      <router-link to="/brand/dashboard" class="nav-item active">
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
      <router-link to="/brand/members" class="nav-item">
        <div class="nav-icon">👤</div>
        <div class="nav-text">会员</div>
      </router-link>
      <router-link to="/brand/promoters" class="nav-item">
        <div class="nav-icon">👥</div>
        <div class="nav-text">推广员</div>
      </router-link>
      <router-link to="/brand/settings" class="nav-item">
        <div class="nav-icon">⚙️</div>
        <div class="nav-text">设置</div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const brandInfo = ref({
  name: '示例品牌',
  logo: 'https://api.dicebear.com/7.x/initials/svg?seed=Brand'
})

const todayStats = ref({
  orders: 0,
  rewards: 0,
  promoters: 0,
  campaigns: 0
})

const recentCampaigns = ref([])

const getStatusText = (status) => {
  const statusMap = {
    active: '进行中',
    paused: '已暂停',
    ended: '已结束'
  }
  return statusMap[status] || status
}

const logout = () => {
  localStorage.removeItem('dmh_token')
  localStorage.removeItem('dmh_user_role')
  router.push('/brand/login')
}

const loadDashboardData = async () => {
  try {
    // TODO: 调用API获取仪表板数据
    todayStats.value = {
      orders: 23,
      rewards: 1580,
      promoters: 8,
      campaigns: 3
    }
    
    recentCampaigns.value = [
      {
        id: 1,
        name: '春节特惠活动',
        description: '新春佳节，推荐好友享双重奖励',
        status: 'active',
        orders: 156,
        rewards: 3120
      },
      {
        id: 2,
        name: '会员招募计划',
        description: '招募品牌会员，享受专属优惠',
        status: 'active',
        orders: 89,
        rewards: 1780
      }
    ]
  } catch (error) {
    console.error('加载仪表板数据失败:', error)
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.brand-dashboard {
  min-height: 100vh;
  background: #f5f7fa;
  padding-bottom: 80px;
}

.top-nav {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px 16px;
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo {
  width: 40px;
  height: 40px;
  border-radius: 8px;
}

.brand-name {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
}

.welcome-text {
  font-size: 14px;
  opacity: 0.8;
  margin: 0;
}

.logout-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0 0 16px 0;
  color: #333;
}

.stats-section {
  padding: 20px 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #667eea;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.quick-actions {
  padding: 0 16px 20px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.action-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
  text-decoration: none;
  color: #333;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.action-card:hover {
  transform: translateY(-2px);
}

.action-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.action-text {
  font-size: 14px;
  font-weight: 500;
}

.campaigns-section {
  padding: 0 16px 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.view-all {
  color: #667eea;
  text-decoration: none;
  font-size: 14px;
}

.campaign-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.campaign-card {
  background: white;
  padding: 16px;
  border-radius: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.campaign-name {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 4px 0;
}

.campaign-desc {
  font-size: 14px;
  color: #666;
  margin: 0 0 8px 0;
}

.campaign-stats {
  display: flex;
  gap: 12px;
}

.stat {
  font-size: 12px;
  color: #999;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.active {
  background: #e8f5e8;
  color: #4caf50;
}

.status-badge.paused {
  background: #fff3e0;
  color: #ff9800;
}

.status-badge.ended {
  background: #fce4ec;
  color: #e91e63;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  display: flex;
  border-top: 1px solid #eee;
  padding: 8px 0;
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: #999;
  padding: 8px;
}

.nav-item.active {
  color: #667eea;
}

.nav-icon {
  font-size: 20px;
  margin-bottom: 4px;
}

.nav-text {
  font-size: 12px;
}
</style>
