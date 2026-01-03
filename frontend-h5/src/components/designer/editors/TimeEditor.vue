<template>
  <van-form>
    <van-field
      v-model="localData.startTime"
      label="开始时间"
      placeholder="请输入开始时间"
      type="datetime-local"
      @blur="updateData"
    />
    <van-field
      v-model="localData.endTime"
      label="结束时间"
      placeholder="请输入结束时间"
      type="datetime-local"
      @blur="updateData"
    />
    <van-field
      v-model="localData.format"
      label="显示格式"
      placeholder="如: YYYY-MM-DD HH:mm"
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
