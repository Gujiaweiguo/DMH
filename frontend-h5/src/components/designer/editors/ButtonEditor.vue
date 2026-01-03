<template>
  <van-form>
    <van-field
      v-model="localData.text"
      label="按钮文字"
      placeholder="请输入按钮文字"
      @blur="updateData"
    />
    
    <van-field label="按钮颜色">
      <template #input>
        <input
          type="color"
          v-model="localData.color"
          @change="updateData"
          style="width: 100%; height: 40px; border: none;"
        />
      </template>
    </van-field>
    
    <van-field
      v-model.number="localData.radius"
      label="圆角大小"
      type="number"
      placeholder="0-50"
      @blur="updateData"
    >
      <template #right-icon>
        <span>px</span>
      </template>
    </van-field>
    
    <van-field label="按钮动作">
      <template #input>
        <van-radio-group v-model="localData.action" @change="updateData">
          <van-radio name="register">报名</van-radio>
          <van-radio name="share">分享</van-radio>
          <van-radio name="custom">自定义</van-radio>
        </van-radio-group>
      </template>
    </van-field>
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