<template>
  <van-form>
    <van-field
      v-model="localData.address"
      label="活动地址"
      placeholder="请输入活动地址"
      @blur="updateData"
    />
    <van-field
      v-model="localData.detail"
      label="详细说明"
      type="textarea"
      placeholder="详细地址说明（可选）"
      rows="3"
      @blur="updateData"
    />
    <van-field
      v-model="localData.mapUrl"
      label="地图链接"
      placeholder="地图链接（可选）"
      @blur="updateData"
    />
  </van-form>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  data: { type: Object, required: true }
})

const emit = defineEmits(['update:data'])

const localData = ref({ ...props.data })

const updateData = () => {
  emit('update:data', localData.value)
}

watch(() => props.data, (newData) => {
  localData.value = { ...newData }
}, { deep: true })
</script>