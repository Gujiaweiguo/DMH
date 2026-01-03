<template>
  <van-form>
    <van-field
      v-model="localData.content"
      label="活动详情"
      type="textarea"
      placeholder="请输入活动详情内容"
      rows="10"
      autosize
      @blur="updateData"
    />
    <van-cell title="提示">
      <template #label>
        <div style="color: #999; font-size: 12px; margin-top: 8px;">
          支持换行和基本文本格式
        </div>
      </template>
    </van-cell>
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