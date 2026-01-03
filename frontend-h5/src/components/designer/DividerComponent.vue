<template>
  <div class="divider-component">
    <van-divider
      :style="{
        borderColor: data.color || '#e5e5e5',
        borderWidth: `${data.height || 1}px`,
        borderStyle: data.style || 'solid'
      }"
    />
  </div>
</template>

<script setup>
defineProps({
  data: { type: Object, required: true },
  preview: { type: Boolean, default: false }
})
</script>

<style scoped>
.divider-component {
  margin: 8px 0;
}
</style>