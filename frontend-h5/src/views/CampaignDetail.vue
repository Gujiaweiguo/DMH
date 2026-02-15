<template>
  <div class="container">
    <div class="header">
      <button class="back-btn" @click="goBack">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
      </button>
      <h1 class="title">{{ campaign.name }}</h1>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    
    <div v-else-if="campaign" class="content">
      <div class="banner">
        <h1 class="title">{{ campaign.name }}</h1>
        <p class="desc">{{ campaign.description }}</p>
      </div>

      <div class="section">
        <div class="section-title">活动时间</div>
        <div class="section-content">
          {{ formatTime(campaign.startTime) }} ~ {{ formatTime(campaign.endTime) }}
        </div>
      </div>

       <div class="section">
         <div class="section-title">报名奖励</div>
         <div class="section-content">
           <div class="reward">每成功报名奖励 {{ campaign.rewardRule?.toFixed(2) || 0 }} 元</div>
         </div>
       </div>

       <!-- 支付配置 -->
       <div v-if="campaign.paymentConfig" class="section">
         <div class="section-title">支付配置</div>
         <div class="section-content payment-config">
           <div class="payment-item">
             <span class="payment-label">是否需要支付：</span>
             <span :class="['payment-value', campaign.paymentConfig.requirePayment ? 'payment-required' : 'payment-free']">
               {{ campaign.paymentConfig.requirePayment ? '需要支付' : '免费活动' }}
             </span>
           </div>
           <div v-if="campaign.paymentConfig.requirePayment" class="payment-item">
             <span class="payment-label">支付类型：</span>
             <span class="payment-value">
               {{ campaign.paymentConfig.paymentType === 'deposit' ? '订金支付' : '全款支付' }}
             </span>
           </div>
           <div v-if="campaign.paymentConfig.requirePayment" class="payment-item">
             <span class="payment-label">
               {{ campaign.paymentConfig.paymentType === 'deposit' ? '订金金额' : '支付金额' }}：
             </span>
             <span class="payment-value amount">
               {{ campaign.paymentConfig.paymentAmount?.toFixed(2) || 0 }} 元
             </span>
           </div>
         </div>
       </div>

       <!-- 统计数据 -->
       <div class="section">
         <div class="section-title">活动统计</div>
         <div class="section-content stats-grid">
           <div class="stat-item">
             <div class="stat-value">{{ campaign.participantCount || 0 }}</div>
             <div class="stat-label">参与人数</div>
           </div>
           <div class="stat-item">
             <div class="stat-value">{{ campaign.orderCount || 0 }}</div>
             <div class="stat-label">订单数</div>
           </div>
           <div class="stat-item">
             <div class="stat-value">{{ getConversionRate() }}%</div>
             <div class="stat-label">转化率</div>
           </div>
           <div class="stat-item">
             <div class="stat-value">{{ campaign.shareCount || 0 }}</div>
             <div class="stat-label">分享次数</div>
           </div>
         </div>
       </div>

       <div class="footer">
         <button class="btn-primary" @click="goForm">立即报名</button>
         <button class="btn-secondary" @click="share">分享给好友</button>
         <button class="btn-poster" @click="generatePoster">生成海报</button>
         <button class="btn-verify" @click="goToVerify">
           订单核销
         </button>
         <button class="btn-payment" @click="goToPayment">
           支付二维码
         </button>
       </div>
    </div>

    <div v-else class="error">活动不存在或已下线</div>
    
    <van-dialog v-model:show="showPosterDialog" title="活动海报" class="poster-dialog">
      <PosterComponent 
        :data="posterData" 
        :preview="true"
        @close="closePosterDialog"
      />
      <template #footer>
        <van-button @click="downloadPoster" type="primary" :loading="posterLoading">
          下载海报
        </van-button>
        <van-button @click="closePosterDialog">
          关闭
        </van-button>
      </template>
    </van-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast, showConfirmDialog, showLoadingToast, closeToast } from 'vant'
import PosterComponent from '@/components/designer/PosterComponent.vue'
import api from '@/services/api'
import {
  formatTime,
  formatReward,
  formatPaymentAmount,
  getPaymentTypeText,
  getPaymentLabel,
  isPaymentRequired,
  calculateConversionRate,
  isBrandAdmin as checkIsBrandAdmin,
  buildSourceData,
  saveSourceToStorage,
  buildPosterDownloadName,
  buildVerifyRoute,
  buildPaymentRoute,
  buildFormRoute
} from './campaignDetail.logic.js'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const campaign = ref(null)
const showPosterDialog = ref(false)
const posterData = ref(null)
const posterLoading = ref(false)

const isBrandAdmin = computed(() => checkIsBrandAdmin())

const source = ref(buildSourceData(route.query))

const saveSource = () => {
  saveSourceToStorage(source.value)
}

const fetchCampaign = async () => {
  const campaignId = route.params.id
  if (!campaignId) {
    loading.value = false
    return
  }

  try {
    const response = await fetch(`/api/v1/h5/campaigns/${campaignId}`)
    if (response.ok) {
      campaign.value = await response.json()
    }
  } catch (error) {
    console.error('加载活动失败', error)
  } finally {
    loading.value = false
  }
}

const goToVerify = () => {
  router.push({
    path: '/verify',
    query: { campaignId: campaign.value?.id }
  })
}

const goToPayment = () => {
  if (!campaign.value || !campaign.value.id) {
    showToast('活动信息加载中')
    return
  }
  
  router.push({
    path: `/campaign/${campaign.value.id}/payment`,
  })
}

const goForm = () => {
  if (isBrandAdmin.value) {
    router.push({
      path: '/verify',
      query: { campaignId: campaign.value?.id }
    })
  } else {
    router.push(`/campaign/${route.params.id}/form`)
  }
}

 const share = () => {
  showConfirmDialog({
    title: '分享提示',
    message: '请使用浏览器的分享功能',
    showCancelButton: false
  })
}

const generatePoster = async (forceRefresh = false) => {
  if (!campaign.value || !campaign.value.id) {
    showToast('活动信息加载中')
    return
  }
  
  posterLoading.value = true
  showLoadingToast({
    message: '生成中...',
    forbidClick: true,
    duration: 0
  })
  
  try {
    const response = await api.generateCampaignPoster(campaign.value.id)
    if (response && response.posterUrl) {
      posterData.value = response
      showPosterDialog.value = true
      showToast('海报生成成功')
    } else {
      showToast('海报生成失败')
    }
  } catch (error) {
    console.error('生成海报失败', error)
    showToast(error.message || '海报生成失败')
  } finally {
    posterLoading.value = false
    closeToast()
  }
}

const closePosterDialog = () => {
  showPosterDialog.value = false
}

const downloadPoster = () => {
  if (!posterData.value || !posterData.value.posterUrl) {
    showToast('海报尚未生成')
    return
  }

  const link = document.createElement('a')
  link.href = posterData.value.posterUrl
  link.download = `poster_${campaign.value.id}.png`
  link.target = '_blank'
  link.click()
  showToast('开始下载')
}

// 计算转化率
const getConversionRate = () => {
  if (!campaign.value) return 0
  return calculateConversionRate(campaign.value.participantCount, campaign.value.orderCount)
}

onMounted(() => {
  saveSource()
  fetchCampaign()
})
</script>

<style scoped>
.container {
  padding: 16px;
  padding-bottom: 80px;
}

.loading, .error {
  text-align: center;
  padding: 60px 20px;
  color: #999;
  font-size: 14px;
}

.banner {
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  border-radius: 16px;
  padding: 24px;
  color: #fff;
  margin-bottom: 24px;
}

.title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 8px;
}

.desc {
  font-size: 14px;
  opacity: 0.9;
  line-height: 1.6;
}

.section {
  background-color: #fff;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
}

.section-title {
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 8px;
  color: #666;
}

.section-content {
  font-size: 14px;
  color: #333;
}

.reward {
  font-size: 20px;
  color: #16a34a;
  font-weight: 600;
}

.footer {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 12px 16px;
  display: flex;
  gap: 12px;
  background-color: #fff;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.1);
}

.btn-primary, .btn-secondary, .btn-poster, .btn-verify, .btn-payment {
  flex: 1;
  border-radius: 999px;
  padding: 14px 0;
  text-align: center;
  font-size: 16px;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background-color: #4f46e5;
  color: #fff;
}

.btn-primary:active {
  background-color: #4338ca;
}

.btn-secondary {
  background-color: #fff;
  color: #666;
  border: 1px solid #e5e7eb;
}

.btn-secondary:active {
  background-color: #f9fafb;
}

.btn-poster {
  background: linear-gradient(135deg, #9333ea, #7c3aed);
  color: #fff;
}

.btn-poster:active {
  background: linear-gradient(135deg, #7c3aed, #5a67d8);
}

.btn-verify {
  background: linear-gradient(135deg, #f093fb, #f5576c);
  color: #fff;
}

.btn-verify:active {
  background: linear-gradient(135deg, #f5576c, #dc2626);
}

.btn-payment {
  background: linear-gradient(135deg, #07c160, #f59e0b);
  color: #fff;
}

.btn-payment:active {
  background: linear-gradient(135deg, #f59e0b, #b91a1a);
}

.poster-dialog {
  .van-dialog__content {
    max-width: 400px;
  }

  .poster-preview {
    padding: 16px;
    text-align: center;
  }
}

/* 支付配置样式 */
.payment-config {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.payment-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  background: #f8fafc;
  border-radius: 6px;
}

.payment-label {
  font-size: 14px;
  color: #6b7280;
  flex-shrink: 0;
}

.payment-value {
  font-size: 14px;
  color: #1f2937;
  font-weight: 500;
  flex-shrink: 0;
}

.payment-value.amount {
  font-size: 18px;
  color: #07c160;
  font-weight: 600;
}

.payment-value.payment-required {
  color: #f59e0b;
}

.payment-value.payment-free {
  color: #10b981;
}

/* 统计数据样式 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 10px;
  background: #f8fafc;
  border-radius: 6px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
}

</style>
