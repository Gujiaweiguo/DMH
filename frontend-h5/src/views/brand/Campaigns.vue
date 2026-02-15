<template>
  <div class="brand-campaigns">
    <!-- 顶部导航栏 -->
    <van-nav-bar
      title="活动管理"
      left-arrow
      @click-left="goBack"
    >
      <template #right>
        <van-button
          type="primary"
          size="small"
          @click="createCampaign"
          icon="plus"
        >
          创建
        </van-button>
      </template>
    </van-nav-bar>

    <!-- 筛选标签 -->
    <van-tabs v-model:active="currentStatus" @change="onStatusChange" sticky>
      <van-tab
        v-for="status in statusTabs"
        :key="status.value"
        :title="status.label"
        :name="status.value"
      />
    </van-tabs>

    <!-- 活动列表 -->
    <div class="campaigns-content">
      <!-- 加载状态 -->
      <van-loading v-if="loading" type="spinner" color="#1989fa" vertical>
        加载中...
      </van-loading>

      <!-- 空状态 -->
      <van-empty
        v-else-if="filteredCampaigns.length === 0"
        image="search"
        :description="currentStatus === 'all' ? '暂无活动' : `暂无${getCurrentStatusText()}活动`"
      >
        <van-button
          v-if="currentStatus === 'all'"
          type="primary"
          @click="createCampaign"
          size="small"
        >
          创建第一个活动
        </van-button>
      </van-empty>

      <!-- 活动卡片列表 -->
      <div v-else class="campaign-list">
        <van-cell-group
          v-for="campaign in filteredCampaigns"
          :key="campaign.id"
          inset
          class="campaign-card"
        >
          <van-cell
            :title="campaign.name"
            :label="campaign.description"
            is-link
            @click="viewCampaign(campaign.id)"
          >
            <template #right-icon>
              <van-tag
                :type="getStatusTagType(campaign.status)"
                size="medium"
              >
                {{ getStatusText(campaign.status) }}
              </van-tag>
            </template>
          </van-cell>
          
           <van-cell title="奖励金额" :value="`¥${campaign.rewardRule}`" />
          <van-cell title="参与人数" :value="campaign.participantCount || 0" />
          
          <!-- 支付配置 -->
          <van-cell v-if="campaign.paymentConfig" title="支付方式">
            <template #value>
              <van-tag
                :type="campaign.paymentConfig.requirePayment ? 'primary' : 'success'"
                size="medium"
              >
                {{ campaign.paymentConfig.requirePayment ? '需支付' : '免费' }}
              </van-tag>
            </template>
          </van-cell>
          
          <van-cell
            v-if="campaign.paymentConfig && campaign.paymentConfig.requirePayment"
            title="支付类型"
            :value="campaign.paymentConfig.paymentType === 'deposit' ? '订金' : '全款'"
          />
          
          <van-cell
            v-if="campaign.paymentConfig && campaign.paymentConfig.requirePayment"
            title="支付金额"
            :value="`¥${campaign.paymentConfig.paymentAmount?.toFixed(2) || '0.00'}`"
          />
          
          <van-cell
            title="活动时间"
            :value="`${formatDate(campaign.startTime)} - ${formatDate(campaign.endTime)}`"
          />
          
          <van-cell>
            <template #title>
              <div class="action-buttons">
                <van-button
                  :type="campaign.status === 'active' ? 'warning' : 'success'"
                  size="small"
                  @click="toggleCampaignStatus(campaign)"
                >
                  {{ campaign.status === 'active' ? '暂停' : '启用' }}
                </van-button>
                
                <van-button
                  type="primary"
                  size="small"
                  @click="editCampaign(campaign.id)"
                >
                  编辑
                </van-button>
                
                <van-button
                  type="default"
                  size="small"
                  @click="viewAnalytics(campaign.id)"
                >
                  数据
                </van-button>
              </div>
            </template>
          </van-cell>
        </van-cell-group>
      </div>
    </div>

    <!-- 底部导航 -->
    <van-tabbar v-model="activeTab" @change="onTabChange">
      <van-tabbar-item icon="wap-home-o" to="/brand/dashboard">
        工作台
      </van-tabbar-item>
      <van-tabbar-item icon="fire-o" to="/brand/campaigns">
        活动
      </van-tabbar-item>
      <van-tabbar-item icon="orders-o" to="/brand/orders">
        订单
      </van-tabbar-item>
      <van-tabbar-item icon="user-o" to="/brand/members">
        会员
      </van-tabbar-item>
      <van-tabbar-item icon="friends-o" to="/brand/promoters">
        推广员
      </van-tabbar-item>
      <van-tabbar-item icon="setting-o" to="/brand/settings">
        设置
      </van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Toast } from 'vant'
import { campaignApi } from '../../services/brandApi.js'
import {
  filterCampaignsByStatus,
  findCampaignStatusLabel,
  formatCampaignDate,
  getCampaignStatusActionText,
  getCampaignStatusTagType,
  getCampaignStatusText,
  getFallbackCampaigns,
  getNextCampaignStatus,
} from './campaigns.logic.js'

const router = useRouter()

const campaigns = ref([])
const loading = ref(false)
const currentStatus = ref('all')
const activeTab = ref(1) // 活动页面对应的tab索引

const statusTabs = [
  { value: 'all', label: '全部' },
  { value: 'active', label: '进行中' },
  { value: 'paused', label: '已暂停' },
  { value: 'ended', label: '已结束' }
]

const filteredCampaigns = computed(() => {
  return filterCampaignsByStatus(campaigns.value, currentStatus.value)
})

const getCurrentStatusText = () => {
  return findCampaignStatusLabel(statusTabs, currentStatus.value)
}

const getStatusText = (status) => {
  return getCampaignStatusText(status)
}

const getStatusTagType = (status) => {
  return getCampaignStatusTagType(status)
}

const formatDate = (dateString) => {
  return formatCampaignDate(dateString)
}

const loadCampaigns = async () => {
  loading.value = true
  
  try {
    console.log('开始加载活动列表...')
    const data = await campaignApi.getCampaigns({ page: 1, pageSize: 50 })
    console.log('活动API响应:', data)
    
    campaigns.value = data.campaigns || []
    console.log('活动数据加载完成:', campaigns.value.length, '个活动')
    
  } catch (error) {
    console.error('加载活动失败:', error)
    Toast.fail('加载活动失败，显示示例数据')
  } finally {
    // 无论API成功还是失败，如果没有数据就显示示例数据
    if (campaigns.value.length === 0) {
      console.log('没有活动数据，使用示例数据')
      campaigns.value = getFallbackCampaigns()
      console.log('示例数据已加载:', campaigns.value.length, '个活动')
    }
    
    loading.value = false
    console.log('活动加载完成，loading状态已关闭')
  }
}

const onStatusChange = (name) => {
  currentStatus.value = name
}

const onTabChange = (index) => {
  // 底部导航切换处理
  activeTab.value = index
}

const goBack = () => {
  router.push('/brand/dashboard')
}

const createCampaign = () => {
  router.push('/brand/campaigns/create')
}

const viewCampaign = (id) => {
  router.push(`/campaign/${id}`)
}

const editCampaign = (id) => {
  router.push(`/brand/campaigns/edit/${id}`)
}

const toggleCampaignStatus = async (campaign) => {
  const newStatus = getNextCampaignStatus(campaign.status)
  const actionText = getCampaignStatusActionText(newStatus)
  
  try {
    // 注意：后端暂未实现更新活动API，这里只是模拟状态切换
    // await campaignApi.updateCampaign(campaign.id, { status: newStatus })
    
    // 临时方案：直接更新本地状态
    campaign.status = newStatus
    Toast.success(`${actionText}成功`)
    console.log(`活动 ${campaign.name} 状态已${actionText}`)
  } catch (error) {
    console.error('更新活动状态失败:', error)
    Toast.fail(`${actionText}失败`)
  }
}

const viewAnalytics = (id) => {
  router.push(`/brand/analytics?campaign=${id}`)
}

onMounted(() => {
  loadCampaigns()
})
</script>

<style scoped>
.brand-campaigns {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 50px; /* 为底部导航留出空间 */
}

.campaigns-content {
  padding: 16px;
  min-height: calc(100vh - 200px);
}

.campaign-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.campaign-card {
  margin-bottom: 12px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: flex-start;
  flex-wrap: wrap;
}

.action-buttons .van-button {
  flex: none;
}

/* 自定义Vant组件样式 */
:deep(.van-nav-bar) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

:deep(.van-nav-bar .van-nav-bar__title) {
  color: white;
}

:deep(.van-nav-bar .van-icon) {
  color: white;
}

:deep(.van-tabs__wrap) {
  background: white;
}

:deep(.van-tab) {
  font-weight: 500;
}

:deep(.van-tab--active) {
  color: #667eea;
}

:deep(.van-tabs__line) {
  background: #667eea;
}

:deep(.van-cell-group--inset) {
  margin: 0 0 12px 0;
  border-radius: 12px;
  overflow: hidden;
}

:deep(.van-cell) {
  padding: 12px 16px;
}

:deep(.van-cell__title) {
  font-weight: 500;
  color: #323233;
}

:deep(.van-cell__label) {
  color: #646566;
  margin-top: 4px;
  line-height: 1.4;
}

:deep(.van-cell__value) {
  color: #323233;
  font-weight: 500;
}

:deep(.van-tag) {
  border-radius: 12px;
  font-size: 12px;
  padding: 2px 8px;
}

:deep(.van-empty) {
  padding: 60px 20px;
}

:deep(.van-loading) {
  padding: 60px 20px;
}

:deep(.van-tabbar) {
  background: white;
  border-top: 1px solid #ebedf0;
}

:deep(.van-tabbar-item) {
  color: #646566;
}

:deep(.van-tabbar-item--active) {
  color: #667eea;
}

/* 响应式设计 */
@media (max-width: 375px) {
  .campaigns-content {
    padding: 12px;
  }
  
  .action-buttons {
    gap: 6px;
  }
  
  .action-buttons .van-button {
    font-size: 12px;
    padding: 4px 8px;
  }
}
</style>
