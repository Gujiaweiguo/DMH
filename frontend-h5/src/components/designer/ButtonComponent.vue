<template>
  <div class="button-component">
    <van-button
      block
      :color="data.color || '#667eea'"
      :style="{ borderRadius: `${data.radius || 8}px` }"
      size="large"
    >
      {{ data.text || '立即报名' }}
    </van-button>
  </div>
</template>

<script setup>
defineProps({
  data: { type: Object, required: true },
  preview: { type: Boolean, default: false }
})
</script>

<style scoped>
.button-component {
  padding: 16px 0;
}
</style>
