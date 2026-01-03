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
</script>

<style scoped>
.poster-component {
  width: 100%;
  min-height: 200px;
}

.poster-image {
  width: 100%;
  border-radius: 8px;
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
</style>