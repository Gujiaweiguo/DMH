<template>
  <div class="detail-component">
    <div class="detail-title">活动详情</div>
    <div class="detail-content" v-html="data.content || '暂无详情'"></div>
  </div>
</template>

<script setup>
defineProps({
  data: { type: Object, required: true },
  preview: { type: Boolean, default: false }
})
</script>

<style scoped>
.detail-component {
  padding: 16px;
}

.detail-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 12px;
  color: #333;
}

.detail-content {
  font-size: 14px;
  line-height: 1.6;
  color: #666;
  white-space: pre-wrap;
}
</style>