<template>
  <div class="payment-qrcode-container">
    <div class="header">
      <button class="back-btn" @click="goBack">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
      </button>
      <h1 class="title">支付二维码</h1>
    </div>

    <div v-if="loading" class="loading-container">
      <van-loading size="40" />
      <p>加载中...</p>
    </div>

    <div v-else-if="error" class="error-container">
      <van-icon name="warning-o" size="48" color="#ee0a0f" />
      <p class="error-message">{{ error }}</p>
      <button class="btn-primary" @click="loadQrcode">重新加载</button>
      <button class="btn-secondary" @click="goBack">返回</button>
    </div>

    <div v-else class="content">
      <div class="qrcode-section">
        <div class="qrcode-box">
          <div v-if="qrcodeData.qrcodeBase64" class="qrcode-image">
            <img :src="qrcodeData.qrcodeBase64" alt="支付二维码" />
          </div>
          <div v-else class="qrcode-loading">
            <van-loading size="32" />
            <span>加载中...</span>
          </div>
        </div>
      </div>

      <div class="info-section">
        <div class="info-item">
          <div class="info-label">活动名称</div>
          <div class="info-value">{{ qrcodeData.campaignName || '-' }}</div>
        </div>
        <div class="info-item">
          <div class="info-label">金额</div>
          <div class="info-value amount">{{ displayAmount(qrcodeData.amount) }}</div>
        </div>
        <div class="info-item">
          <div class="info-label">过期时间</div>
          <div class="info-value">{{ formatExpireTime(qrcodeData.expireAt) }}</div>
        </div>
      </div>

      <div class="tips-section">
        <van-notice-bar color="#1989fa" left-icon="info-o">
          请使用微信扫描二维码完成支付
        </van-notice-bar>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import api from '@/services/api'
import {
  formatExpireTime,
  formatAmount,
  getDefaultQrcodeData,
  buildQrcodeData,
  validateCampaignId
} from './paymentQrcode.logic.js'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const error = ref(null)
const qrcodeData = ref(getDefaultQrcodeData())

const loadQrcode = async () => {
  const campaignId = route.params.id
  const validation = validateCampaignId(campaignId)
  if (!validation.valid) {
    error.value = validation.error
    return
  }

  loading.value = true
  error.value = null
  showLoadingToast({
    message: '加载中...',
    forbidClick: true,
    duration: 0
  })

  try {
    const response = await api.getPaymentQrcode(campaignId)
    const data = buildQrcodeData(response)
    
    if (data) {
      qrcodeData.value = data
    } else {
      error.value = response.msg || '加载失败'
    }
  } catch (err) {
    console.error('加载支付二维码失败', err)
    error.value = err.message || '加载失败，请重试'
  } finally {
    loading.value = false
    closeToast()
  }
}

const goBack = () => {
  router.go(-1)
}

const displayAmount = (amount) => {
  return amount ? `¥${formatAmount(amount)}` : '-'
}

onMounted(() => {
  loadQrcode()
})
</script>

<style scoped>
.payment-qrcode-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  background: white;
  padding: 12px 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.back-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: #333;
}

.loading-container, .error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  min-height: 400px;
}

.loading-container p, .error-message {
  margin-top: 16px;
  font-size: 14px;
  color: #666;
}

.error-icon {
  margin-bottom: 16px;
}

.btn-primary, .btn-secondary {
  margin-top: 24px;
  margin-right: 12px;
  width: 140px;
}

.content {
  padding: 24px;
}

.qrcode-section {
  background: white;
  border-radius: 12px;
  padding: 32px;
  margin-bottom: 24px;
  text-align: center;
}

.qrcode-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px;
}

.qrcode-image img {
  width: 280px;
  height: 280px;
  display: block;
}

.qrcode-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.qrcode-loading span {
  font-size: 14px;
  color: #999;
}

.info-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 14px;
  color: #999;
}

.info-value {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.info-value.amount {
  color: #ee0a0f;
  font-size: 24px;
  font-weight: 600;
}

.tips-section {
  margin-top: 16px;
}
</style>
