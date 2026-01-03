<template>
  <div class="highlight-component">
    <div class="highlight-title">活动亮点</div>
    <div class="highlight-list">
      <div
        v-for="(item, index) in data.items"
        :key="index"
        class="highlight-item"
      >
        <van-icon :name="item.icon || 'star'" size="18" color="#667eea" />
        <span>{{ item.text }}</span>
      </div>
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
.highlight-component {
  padding: 16px;
  background: #f7f8fa;
  border-radius: 8px;
}

.highlight-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 12px;
  color: #333;
}

.highlight-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.highlight-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}
</style>