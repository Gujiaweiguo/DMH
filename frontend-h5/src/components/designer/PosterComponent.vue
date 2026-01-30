<template>
  <div class="poster-component">
    <van-image
      v-if="data.imageUrl"
      :src="data.imageUrl"
      fit="cover"
      class="poster-image"
    />
    <div v-else class="poster-placeholder">
      <van-icon name="photo-o" size="48" />
      <p>暂无海报图片</p>
    </div>
    
    <div v-if="data.qrcodeUrl" class="qrcode-section">
      <div class="qrcode-label">支付二维码</div>
      <div class="qrcode-content">
        <div class="qrcode-box">
          <div v-if="data.qrcodeBase64" class="qrcode-image">
            <img :src="data.qrcodeBase64" alt="支付二维码" />
          </div>
          <div v-else class="qrcode-loading">
            <van-loading size="24" />
            <span>加载中...</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  data: {
    type: Object,
    required: true
  },
  preview: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['download', 'close'])

const downloadPoster = () => {
  emit('download')
}
</script>

<style scoped>
.poster-component {
  width: 100%;
  min-height: 200px;
}

.poster-image {
  width: 100%;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.poster-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  background: #f5f5f5;
  border-radius: 8px;
  color: #999;
}

.poster-placeholder p {
  margin-top: 12px;
  font-size: 14px;
}

.qrcode-section {
  margin-top: 16px;
  padding: 16px;
  background: #f9f9f9;
  border-radius: 8px;
}

.qrcode-label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 12px;
  text-align: center;
}

.qrcode-content {
  display: flex;
  justify-content: center;
}

.qrcode-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px;
  background: white;
  border-radius: 8px;
}

.qrcode-image img {
  width: 200px;
  height: 200px;
  display: block;
}

.qrcode-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #999;
}

.qrcode-loading span {
  font-size: 14px;
}
</style>
