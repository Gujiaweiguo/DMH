<template>
  <van-form>
    <div v-for="(item, index) in localData.items" :key="index" class="highlight-item-editor">
      <van-field
        v-model="item.text"
        :label="`亮点${index + 1}`"
        placeholder="请输入亮点内容"
        @blur="updateData"
      >
        <template #button>
          <van-button
            size="small"
            type="danger"
            @click="removeItem(index)"
          >
            删除
          </van-button>
        </template>
      </van-field>
      <van-field
        v-model="item.icon"
        label="图标"
        placeholder="图标名称（如：star）"
        @blur="updateData"
      />
    </div>
    
    <van-button
      block
      type="primary"
      @click="addItem"
      style="margin-top: 16px;"
    >
      添加亮点
    </van-button>
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

const addItem = () => {
  if (!localData.value.items) {
    localData.value.items = []
  }
  localData.value.items.push({ icon: 'star', text: '' })
  updateData()
}

const removeItem = (index) => {
  localData.value.items.splice(index, 1)
  updateData()
}

watch(() => props.data, (newData) => {
  localData.value = { ...newData }
}, { deep: true })
</script>

<style scoped>
.highlight-item-editor {
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebedf0;
}
</style>