<template>
  <van-form>
    <van-field
      v-model="localData.imageUrl"
      label="图片URL"
      placeholder="请输入图片URL"
      @blur="updateData"
    />
    <van-cell title="或选择素材">
      <template #right-icon>
        <van-button type="primary" size="small" @click="selectMaterial">
          选择素材
        </van-button>
      </template>
    </van-cell>
    <van-field
      v-model="localData.alt"
      label="图片描述"
      placeholder="图片描述（可选）"
      @blur="updateData"
    />
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

const selectMaterial = () => {
  Toast('素材管理功能开发中')
}

watch(() => props.data, (newData) => {
  localData.value = { ...newData }
}, { deep: true })
</script>