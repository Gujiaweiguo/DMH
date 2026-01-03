<template>
  <van-form>
    <van-field label="分割线颜色">
      <template #input>
        <input
          type="color"
          v-model="localData.color"
          @change="updateData"
          style="width: 100%; height: 40px; border: none;"
        />
      </template>
    </van-field>
    
    <van-field label="线条样式">
      <template #input>
        <van-radio-group v-model="localData.style" @change="updateData">
          <van-radio name="solid">实线</van-radio>
          <van-radio name="dashed">虚线</van-radio>
          <van-radio name="dotted">点线</van-radio>
        </van-radio-group>
      </template>
    </van-field>
    
    <van-field
      v-model.number="localData.height"
      label="线条粗细"
      type="number"
      placeholder="1-10"
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