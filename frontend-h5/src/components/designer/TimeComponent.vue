<template>
  <div class="time-component">
    <van-icon name="clock-o" size="20" />
    <div class="time-content">
      <div class="time-label">活动时间</div>
      <div class="time-value">
        {{ formatTime(data.startTime) }} - {{ formatTime(data.endTime) }}
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  data: { type: Object, required: true },
  preview: { type: Boolean, default: false }
})

const formatTime = (time) => {
  if (!time) return '待设置'
  return new Date(time).toLocaleString('zh-CN', {
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.time-component {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.time-content {
  flex: 1;
}

.time-label {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.time-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}
</style>