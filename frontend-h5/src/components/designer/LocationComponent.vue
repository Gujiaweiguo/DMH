<template>
  <div class="location-component">
    <van-icon name="location-o" size="20" />
    <div class="location-content">
      <div class="location-label">活动地点</div>
      <div class="location-value">{{ data.address || '待设置' }}</div>
      <div v-if="data.detail" class="location-detail">{{ data.detail }}</div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  data: { type: Object, required: true },
  preview: { type: Boolean, default: false }
})
</script>

<style scoped>
.location-component {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.location-content {
  flex: 1;
}

.location-label {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.location-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
  margin-bottom: 4px;
}

.location-detail {
  font-size: 12px;
  color: #666;
}
</style>