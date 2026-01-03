<template>
  <div class="title-component">
    <h2 class="title" :style="{ color: data.color }">{{ data.title || '活动标题' }}</h2>
    <p v-if="data.subtitle" class="subtitle">{{ data.subtitle }}</p>
  </div>
</template>

<script setup>
defineProps({
  data: { type: Object, required: true },
  preview: { type: Boolean, default: false }
})
</script>

<style scoped>
.title-component {
  text-align: center;
  padding: 16px 0;
}

.title {
  font-size: 24px;
  font-weight: bold;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 14px;
  color: #666;
  margin: 0;
}
</style>