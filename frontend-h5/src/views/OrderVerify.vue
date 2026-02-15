<template>
  <div class="container">
    <div class="header">
      <button class="back-btn" @click="goBack">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
      </button>
      <h1 class="title">订单核销</h1>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="!orderInfo" class="error">
      <div class="error-icon">❌</div>
      <p class="error-msg">{{ errorMessage }}</p>
      <button class="btn-primary" @click="goBack">返回</button>
    </div>

    <div v-else class="content">
      <!-- 扫码输入 -->
      <div class="scan-section">
        <div class="scan-title">请扫描订单核销码</div>
        <div class="scan-input-wrapper">
          <input 
            v-model="verificationCode" 
            type="text" 
            placeholder="请输入核销码" 
            class="scan-input"
            @keyup.enter="handleScan"
            autofocus
          />
          <button 
            class="scan-btn" 
            @click="handleScan" 
            :disabled="!verificationCode || scanning"
          >
            {{ scanning ? '扫描中...' : '扫码' }}
          </button>
        </div>
      </div>

      <!-- 订单信息 -->
      <div v-if="orderInfo" class="order-info">
        <div class="info-card">
          <div class="info-header">
            <span class="info-title">订单信息</span>
            <span class="info-status" :class="orderInfo.verificationStatus">
              {{ verificationStatusText }}
            </span>
          </div>
          <div class="info-body">
            <div class="info-row">
              <span class="label">订单号：</span>
              <span class="value">{{ orderInfo.id }}</span>
            </div>
            <div class="info-row">
              <span class="label">手机号：</span>
              <span class="value">{{ orderInfo.phone }}</span>
            </div>
            <div class="info-row">
              <span class="label">报名时间：</span>
              <span class="value">{{ formatDate(orderInfo.createdAt) }}</span>
            </div>
            <div class="info-row" v-if="orderInfo.formData">
              <span class="label">报名信息：</span>
              <div class="form-data-value">{{ formatFormDataDisplay(orderInfo.formData) }}</div>
            </div>
          </div>
        </div>

        <!-- 核销操作按钮 -->
        <div class="actions">
          <button 
            v-if="orderInfo.verificationStatus === 'unverified'" 
            class="btn-verify" 
            @click="verifyOrder" 
            :disabled="processing"
          >
            核销订单
          </button>
          <button 
            v-if="orderInfo.verificationStatus === 'verified'" 
            class="btn-unverify" 
            @click="unverifyOrder" 
            :disabled="processing"
          >
            取消核销
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast, showConfirmDialog, showLoadingToast, closeToast } from 'vant'
import { api, orderApi } from '@/services/api'
import {
  getVerificationStatusText,
  formatDate as formatDateUtil,
  formatFormData,
  canVerify as canVerifyOrder,
  canUnverify as canUnverifyOrder,
  isValidVerificationCode,
  buildScanButtonText
} from './orderVerify.logic.js'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const scanning = ref(false)
const verificationCode = ref('')
const orderInfo = ref(null)
const processing = ref(false)
const errorMessage = ref('')

const verificationStatusText = computed(() => {
  return getVerificationStatusText(orderInfo.value?.verificationStatus)
})

const scanButtonText = computed(() => buildScanButtonText(scanning.value))

const handleScan = async () => {
  if (!isValidVerificationCode(verificationCode.value)) {
    showToast('请输入核销码')
    return
  }

  scanning.value = true
  processing.value = true

  try {
    const response = await orderApi.scanOrderCode(verificationCode.value)
    orderInfo.value = response
    showToast('订单信息加载成功')
  } catch (error) {
    console.error('扫码失败', error)
    errorMessage.value = error.message || '订单不存在或核销码无效'
    showToast(error.message || '订单不存在')
  } finally {
    scanning.value = false
    processing.value = false
  }
}

const verifyOrder = () => {
  if (!orderInfo.value) {
    return
  }

  showConfirmDialog({
    title: '确认核销',
    message: `确认核销订单 ${orderInfo.value.id}?`,
    confirmButtonText: '确认',
    cancelButtonText: '取消'
  }).then(async () => {
    processing.value = true

    try {
      await orderApi.verifyOrder(verificationCode.value)
      orderInfo.value.verificationStatus = 'verified'
      showToast('订单核销成功')
    } catch (error) {
      console.error('核销失败', error)
      showToast(error.message || '核销失败')
    } finally {
      processing.value = false
    }
  })
}

const unverifyOrder = () => {
  if (!orderInfo.value) {
    return
  }

  showConfirmDialog({
    title: '取消核销',
    message: `确认取消订单 ${orderInfo.value.id} 的核销状态?`,
    confirmButtonText: '确认',
    cancelButtonText: '取消'
  }).then(async () => {
    processing.value = true

    try {
      await orderApi.unverifyOrder(verificationCode.value)
      orderInfo.value.verificationStatus = 'unverified'
      showToast('已取消核销')
    } catch (error) {
      console.error('取消核销失败', error)
      showToast(error.message || '取消核销失败')
    } finally {
      processing.value = false
    }
  })
}

const formatDate = (time) => formatDateUtil(time)

const formatFormDataDisplay = (formDataStr) => formatFormData(formDataStr)

const goBack = () => {
  router.back()
}

onMounted(() => {
  const code = route.query.code
  if (code) {
    verificationCode.value = code
    handleScan()
  }
})
</script>

<style scoped>
.container {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 16px;
}

.header {
  display: flex;
  align-items: center;
  padding: 16px 0;
  margin-bottom: 20px;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.back-btn {
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  padding: 8px;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #999;
  font-size: 16px;
}

.error {
  text-align: center;
  padding: 40px 20px;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-msg {
  color: #f56c6c;
  font-size: 16px;
  margin-bottom: 24px;
}

.content {
  max-width: 600px;
  margin: 0 auto;
}

.scan-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.scan-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  text-align: center;
}

.scan-input-wrapper {
  display: flex;
  gap: 12px;
}

.scan-input {
  flex: 1;
  padding: 14px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 16px;
  outline: none;
}

.scan-input:focus {
  border-color: #4f46e5;
}

.scan-btn {
  padding: 14px 24px;
  background-color: #4f46e5;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
}

.scan-btn:disabled {
  background-color: #d1d5db;
  cursor: not-allowed;
}

.order-info {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.info-card {
  padding: 20px;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.info-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.info-status {
  padding: 6px 16px;
  border-radius: 20px;
  font-weight: 600;
  font-size: 14px;
}

.info-status.unverified {
  background-color: #dcfce7;
  color: #16a34a;
}

.info-status.verified {
  background-color: #52c41a;
  color: #fff;
}

.info-status.cancelled {
  background-color: #fee2e2;
  color: #dc2626;
}

.info-body {
  margin-bottom: 20px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 14px;
}

.label {
  color: #666;
  min-width: 100px;
}

.value {
  color: #333;
  text-align: right;
}

.form-data-value {
  background-color: #f9fafb;
  padding: 8px 12px;
  border-radius: 6px;
  color: #111827;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}

.actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  padding-top: 24px;
}

.btn-verify, .btn-unverify {
  flex: 1;
  padding: 14px 32px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 16px;
  cursor: pointer;
  border: none;
}

.btn-verify {
  background-color: #52c41a;
  color: #fff;
}

.btn-unverify {
  background-color: #fff;
  border: 1px solid #dc2626;
  color: #dc2626;
}

.btn-primary {
  padding: 12px 24px;
  background: #666;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}
</style>
