<template>
  <div class="order-verification">
    <!-- 顶部导航栏 -->
    <van-nav-bar
      title="订单核销"
      left-text="返回"
      left-arrow
      @click-left="goBack"
    />

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 扫码区域 -->
      <div class="scan-section">
        <div class="section-header">
          <h2>1. 扫描核销码</h2>
          <p class="section-desc">请扫描订单核销码或手动输入</p>
        </div>

        <div class="scan-container">
          <van-loading v-if="state.loading.scan" size="32px">
            启动摄像头...
          </van-loading>

          <div v-else-if="!state.scannedOrder" class="camera-placeholder">
            <div class="camera-icon">
              <van-icon name="scan" size="48px" />
            </div>
            <div class="camera-text">点击下方按钮打开摄像头</div>
            <div class="camera-actions">
              <van-button
                type="primary"
                size="small"
                :disabled="state.isScanning"
                @click="startScan"
              >
                开始扫码
              </van-button>
              <van-button
                type="default"
                size="small"
                :disabled="!state.isScanning"
                @click="stopScan"
              >
                停止扫码
              </van-button>
            </div>
            <div v-show="state.isScanning" id="qr-reader" class="qr-reader"></div>
          </div>

          <div v-else class="scanned-result">
            <van-cell-group>
              <van-cell title="订单号" :value="state.scannedOrder.orderCode || '-'" />
              <van-cell title="订单状态" :value="getOrderStatusText(state.scannedOrder.verifyStatus)" />
            </van-cell-group>

            <van-button type="primary" size="large" block @click="scanAgain">
              重新扫描
            </van-button>
          </div>
        </div>

        <!-- 手动输入区域 -->
        <van-cell-group>
          <van-field
            v-model="state.orderCode"
            label="核销码"
            placeholder="请输入订单核销码"
            clearable
          />
          <van-button
            type="primary"
            size="large"
            block
            :loading="state.loading.verify"
            @click="verifyOrderCode"
          >
            验证订单
          </van-button>
        </van-cell-group>
      </div>

      <!-- 订单详情区域 -->
      <div v-if="state.scannedOrder" class="detail-section">
        <div class="section-header">
          <h2>2. 订单详情</h2>
        </div>

        <van-cell-group>
          <van-cell title="手机号码" :value="state.scannedOrder.userPhone || '-'" />
          <van-cell title="支付状态" :value="getPaymentStatusText(state.scannedOrder.paymentStatus)" />
          <van-cell title="订单状态" :value="state.scannedOrder.orderStatus || '-'" />
          <van-cell title="核销状态" :value="getVerifyStatusText(state.scannedOrder.verifyStatus)" />
          <van-cell title="核销时间" :value="state.scannedOrder.verifiedAt || '未核销'" />
        </van-cell-group>
      </div>

      <!-- 操作区域 -->
      <div v-if="state.scannedOrder && state.scannedOrder.verifyStatus === 'unverified'" class="action-section">
        <div class="section-header">
          <h2>3. 核销订单</h2>
        </div>

        <van-cell-group>
          <van-field
            v-model="state.verifyNotes"
            type="textarea"
            label="核销备注"
            placeholder="请输入核销备注（可选）"
            rows="3"
          />
        </van-cell-group>

        <div class="action-buttons">
          <van-button
            type="success"
            size="large"
            block
            :loading="state.loading.verify"
            @click="confirmVerify"
          >
            确认核销
          </van-button>

          <van-button
            type="danger"
            size="large"
            block
            :loading="state.loading.cancel"
            @click="showCancelDialog"
          >
            取消核销
          </van-button>
        </div>
      </div>

      <!-- 已核销区域 -->
      <div v-if="state.scannedOrder && state.scannedOrder.verifyStatus === 'verified'" class="verified-section">
        <van-notice-bar type="success" left-icon="success">
          订单已核销
        </van-notice-bar>

        <van-cell-group>
          <van-cell title="核销人" :value="state.scannedOrder.verifiedBy || '-'" />
          <van-cell title="核销备注" :value="state.scannedOrder.notes || '-'" />
          <van-cell title="核销时间" :value="state.scannedOrder.verifiedAt || '-'" />
        </van-cell-group>

        <div v-if="state.scannedOrder.cancelledAt" class="cancelled-info">
          <van-notice-bar type="danger" left-icon="warning">
            已取消核销
          </van-notice-bar>

          <van-cell-group>
            <van-cell title="取消人" :value="state.scannedOrder.cancelledBy || '-'" />
            <van-cell title="取消原因" :value="state.scannedOrder.cancelReason || '-'" />
            <van-cell title="取消时间" :value="state.scannedOrder.cancelledAt || '-'" />
          </van-cell-group>
        </div>
      </div>
    </div>

    <!-- 取消核销对话框 -->
    <van-dialog
      v-model:show="state.showCancelDialog"
      title="取消核销"
      show-cancel-button
      @confirm="cancelVerify"
    >
      <van-field
        v-model="state.cancelReason"
        type="textarea"
        label="取消原因"
        placeholder="请输入取消原因（必填）"
        rows="3"
      />
    </van-dialog>
  </div>
</template>

<script setup>
import { Html5Qrcode } from 'html5-qrcode'
import { Toast } from 'vant'
import { onMounted, onUnmounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { orderApi } from '../../services/brandApi'
import {
  buildUnverifiedOrderState,
  buildVerifiedOrderState,
  extractCodeFromText,
  getOrderStatusText as mapOrderStatusText,
  getPaymentStatusText as mapPaymentStatusText,
  getVerifyStatusText as mapVerifyStatusText,
  isOperationSuccess,
} from './orderVerification.logic'

const router = useRouter()

// 响应式状态
const state = reactive({
  loading: {
    scan: false,
    verify: false,
    cancel: false
  },
  orderCode: '',
  scannedOrder: null,
  verifyNotes: '',
  cancelReason: '',
  showCancelDialog: false,
  isScanning: false
})

let qrScanner = null

// 获取路由参数


// 扫描订单核销码
const scanOrderCode = async () => {
  state.loading.scan = true
  state.isScanning = true
  
  try {
    const response = await orderApi.scanOrderCode(state.orderCode)
    
    if (response.code === 0 && response.data) {
      state.scannedOrder = response.data
      Toast.success('订单扫描成功')
    } else {
      Toast.fail(response.msg || '订单不存在或核销码无效')
      state.scannedOrder = null
    }
  } catch (error) {
    console.error('扫描订单失败:', error)
    Toast.fail('扫描失败，请重试')
    state.scannedOrder = null
  } finally {
    state.loading.scan = false
    state.isScanning = false
  }
}

// biome-ignore lint/correctness/noUnusedVariables: used in template
const startScan = async () => {
  if (state.isScanning) {
    return
  }

  state.isScanning = true
  state.loading.scan = true

  try {
    if (!qrScanner) {
      qrScanner = new Html5Qrcode('qr-reader')
    }

    await qrScanner.start(
      { facingMode: 'environment' },
      { fps: 10, qrbox: { width: 250, height: 250 } },
      async (decodedText) => {
        const code = extractCodeFromText(decodedText)
        if (!code) {
          Toast.fail('二维码内容无效')
          return
        }

        state.orderCode = code
        await stopScan()
        verifyOrderCode()
      },
      () => {}
    )
  } catch (error) {
    console.error('启动扫码失败:', error)
    Toast.fail('无法启动摄像头，请检查权限')
    state.isScanning = false
  } finally {
    state.loading.scan = false
  }
}

const stopScan = async () => {
  if (!qrScanner || !state.isScanning) {
    state.isScanning = false
    return
  }

  try {
    await qrScanner.stop()
    await qrScanner.clear()
  } catch (error) {
    console.error('停止扫码失败:', error)
  } finally {
    state.isScanning = false
  }
}

// 重新扫描
// biome-ignore lint/correctness/noUnusedVariables: used in template
const scanAgain = () => {
  state.scannedOrder = null
  state.orderCode = ''
  state.verifyNotes = ''
}

// 验证订单码
const verifyOrderCode = () => {
  if (!state.orderCode.trim()) {
    Toast('请输入核销码')
    return
  }

  scanOrderCode()
}

// 确认核销
// biome-ignore lint/correctness/noUnusedVariables: used in template
const confirmVerify = async () => {
  if (!state.scannedOrder) {
    Toast('请先扫描订单')
    return
  }

  state.loading.verify = true
  
  try {
    const response = await orderApi.verifyOrder(state.orderCode, {
      notes: state.verifyNotes
    })
    
    if (isOperationSuccess(response)) {
      state.scannedOrder = buildVerifiedOrderState(
        state.scannedOrder,
        state.verifyNotes,
        new Date().toISOString(),
      )
      Toast.success('订单核销成功')
      state.verifyNotes = ''
    } else {
      Toast.fail(response.msg || '核销失败')
    }
  } catch (error) {
    console.error('核销订单失败:', error)
    Toast.fail('核销失败，请重试')
  } finally {
    state.loading.verify = false
  }
}

// 显示取消对话框
// biome-ignore lint/correctness/noUnusedVariables: used in template
const showCancelDialog = () => {
  state.showCancelDialog = true
}

// 取消核销
// biome-ignore lint/correctness/noUnusedVariables: used in template
const cancelVerify = async () => {
  if (!state.cancelReason.trim()) {
    Toast('请输入取消原因')
    return
  }

  state.loading.cancel = true
  
  try {
    const response = await orderApi.unverifyOrder(state.orderCode, {
      reason: state.cancelReason
    })
    
    if (isOperationSuccess(response)) {
      state.scannedOrder = buildUnverifiedOrderState(
        state.scannedOrder,
        state.cancelReason,
        new Date().toISOString(),
      )
      Toast.success('已取消核销')
      state.cancelReason = ''
      state.showCancelDialog = false
    } else {
      Toast.fail(response.msg || '取消失败')
    }
  } catch (error) {
    console.error('取消核销失败:', error)
    Toast.fail('取消失败，请重试')
  } finally {
    state.loading.cancel = false
  }
}

// biome-ignore lint/correctness/noUnusedVariables: used in template
const getOrderStatusText = (status) => mapOrderStatusText(status)

// biome-ignore lint/correctness/noUnusedVariables: used in template
const getPaymentStatusText = (status) => mapPaymentStatusText(status)

// biome-ignore lint/correctness/noUnusedVariables: used in template
const getVerifyStatusText = (status) => mapVerifyStatusText(status)

// 返回
// biome-ignore lint/correctness/noUnusedVariables: used in template
const goBack = () => {
  router.back()
}

// 生命周期
onMounted(() => {
  console.log('OrderVerification页面已加载')
})

onUnmounted(() => {
  stopScan()
})
</script>

<style scoped>
.order-verification {
  min-height: 100vh;
  background: #f5f5f5;
}

.main-content {
  padding-bottom: 80px;
}

.section-header {
  padding: 20px;
  background: white;
  border-bottom: 1px solid #e5e5e5;
}

.section-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #323233;
}

.section-desc {
  margin: 8px 0 0;
  font-size: 14px;
  color: #969796;
}

.scan-section,
.detail-section,
.action-section,
.verified-section {
  margin-bottom: 20px;
  background: white;
  border-radius: 12px;
  overflow: hidden;
}

.scan-container {
  padding: 20px;
  background: #f8f8fa;
  border-radius: 8px;
  text-align: center;
}

.camera-placeholder {
  padding: 40px 20px;
}

.camera-icon {
  margin-bottom: 16px;
  color: #1989fa;
}

.camera-text {
  font-size: 14px;
  color: #969796;
}

.camera-actions {
  margin-top: 16px;
  display: flex;
  gap: 12px;
  justify-content: center;
}

.qr-reader {
  width: 100%;
  max-width: 320px;
  margin: 16px auto 0;
}

.scanned-result {
  margin-top: 20px;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
}

.action-buttons .van-button {
  height: 48px;
  font-size: 16px;
  font-weight: 600;
}

.verified-section {
  background: #e8f5e9;
}

.cancelled-info {
  margin-top: 20px;
  background: #fff3cd;
  padding: 15px;
  border-radius: 8px;
}
</style>
