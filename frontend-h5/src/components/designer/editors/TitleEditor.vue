<template>
  <van-form>
    <van-field
      v-model="localData.title"
      label="活动标题"
      placeholder="请输入活动标题"
      @blur="updateData"
    />
    <van-field
      v-model="localData.subtitle"
      label="副标题"
      placeholder="副标题（可选）"
      @blur="updateData"
    />
    <van-field label="文字颜色">
      <template #input>
        <input
          type="color"
          v-model="localData.color"
          @change="updateData"
          style="width: 100%; height: 40px; border: none;"
        />
      </template>
    </van-field>
    <van-cell title="AI调色">
      <template #right-icon>
        <van-button type="primary" size="small" @click="aiColor">
          AI调色
        </van-button>
      </template>
    </van-cell>
  </van-form>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Toast } from 'vant'

const props = defineProps({
  data: { type: Object, required: true }
})

const emit = defineEmits(['update:data'])

const localData = ref({ ...props.data })

const updateData = () => {
  emit('update:data', localData.value)
}

const aiColor = () => {
  Toast('AI调色功能开发中')
}

watch(() => props.data, (newData) => {
  localData.value = { ...newData }
}, { deep: true })
</script>