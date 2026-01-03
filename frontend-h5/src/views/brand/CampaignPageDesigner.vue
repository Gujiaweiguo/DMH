<template>
  <div class="page-designer">
    <!-- 顶部导航栏 -->
    <van-nav-bar
      title="页面设计"
      left-arrow
      @click-left="goBack"
    >
      <template #right>
        <van-button type="primary" size="small" @click="previewPage">
          预览
        </van-button>
        <van-button type="success" size="small" @click="savePage" style="margin-left: 8px;">
          保存活动
        </van-button>
      </template>
    </van-nav-bar>

    <!-- 设计画布 -->
    <div class="design-canvas">
      <!-- 空状态 -->
      <van-empty
        v-if="components.length === 0"
        description="暂无组件，点击下方按钮添加"
      />

      <!-- 组件列表 -->
      <div v-else class="component-list">
        <div
          v-for="(component, index) in components"
          :key="component.id"
          class="canvas-component"
          @click="editComponent(component)"
        >
          <!-- 组件内容 -->
          <component
            :is="getComponentType(component.type)"
            :data="component.data"
            :style="component.style"
            :preview="false"
          />

          <!-- 操作按钮 -->
          <div class="component-actions">
            <van-button
              icon="arrow-up"
              size="mini"
              @click.stop="moveComponent(index, 'up')"
              :disabled="index === 0"
            />
            <van-button
              icon="arrow-down"
              size="mini"
              @click.stop="moveComponent(index, 'down')"
              :disabled="index === components.length - 1"
            />
            <van-button
              icon="edit"
              size="mini"
              type="primary"
              @click.stop="editComponent(component)"
            />
            <van-button
              icon="delete-o"
              size="mini"
              type="danger"
              @click.stop="removeComponent(component.id)"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 底部添加组件按钮 -->
    <div class="bottom-actions">
      <van-button
        type="primary"
        block
        icon="plus"
        @click="showComponentLibrary = true"
      >
        添加组件
      </van-button>
    </div>

    <!-- 组件库弹出层 -->
    <van-action-sheet
      v-model:show="showComponentLibrary"
      title="选择组件"
      :actions="componentLibrary"
      @select="onSelectComponent"
      cancel-text="取消"
    />

    <!-- 组件编辑弹出层 -->
    <van-popup
      v-model:show="showEditor"
      position="bottom"
      :style="{ height: '70%' }"
      round
    >
      <div class="component-editor">
        <van-nav-bar
          :title="`编辑${getComponentName(editingComponent?.type)}`"
          left-text="取消"
          right-text="保存"
          @click-left="closeEditor"
          @click-right="saveComponentData"
        />

        <div class="editor-content">
          <component
            :is="getEditorType(editingComponent?.type)"
            v-if="editingComponent"
            v-model:data="editingComponent.data"
          />
        </div>
      </div>
    </van-popup>

    <!-- 预览弹出层 -->
    <van-popup
      v-model:show="showPreview"
      position="right"
      :style="{ width: '100%', height: '100%' }"
    >
      <div class="preview-container">
        <van-nav-bar
          title="预览"
          left-arrow
          @click-left="showPreview = false"
        />

        <div class="preview-content">
          <div
            v-for="component in components"
            :key="component.id"
            class="preview-component"
          >
            <component
              :is="getComponentType(component.type)"
              :data="component.data"
              :style="component.style"
              :preview="true"
            />
          </div>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Toast, Dialog } from 'vant'
import { campaignApi } from '../../services/brandApi.js'

// 导入组件（稍后创建）
import PosterComponent from '../../components/designer/PosterComponent.vue'
import TitleComponent from '../../components/designer/TitleComponent.vue'
import TimeComponent from '../../components/designer/TimeComponent.vue'
import LocationComponent from '../../components/designer/LocationComponent.vue'
import HighlightComponent from '../../components/designer/HighlightComponent.vue'
import DetailComponent from '../../components/designer/DetailComponent.vue'
import DividerComponent from '../../components/designer/DividerComponent.vue'
import ButtonComponent from '../../components/designer/ButtonComponent.vue'

// 导入编辑器（稍后创建）
import PosterEditor from '../../components/designer/editors/PosterEditor.vue'
import TitleEditor from '../../components/designer/editors/TitleEditor.vue'
import TimeEditor from '../../components/designer/editors/TimeEditor.vue'
import LocationEditor from '../../components/designer/editors/LocationEditor.vue'
import HighlightEditor from '../../components/designer/editors/HighlightEditor.vue'
import DetailEditor from '../../components/designer/editors/DetailEditor.vue'
import DividerEditor from '../../components/designer/editors/DividerEditor.vue'
import ButtonEditor from '../../components/designer/editors/ButtonEditor.vue'

const router = useRouter()
const route = useRoute()

const components = ref([])
const showComponentLibrary = ref(false)
const showEditor = ref(false)
const showPreview = ref(false)
const editingComponent = ref(null)
const loading = ref(false)

// 组件库配置
const componentLibrary = [
  { name: '活动海报', value: 'poster', icon: 'photo-o' },
  { name: '活动标题', value: 'title', icon: 'font-o' },
  { name: '活动时间', value: 'time', icon: 'clock-o' },
  { name: '活动地点', value: 'location', icon: 'location-o' },
  { name: '活动亮点', value: 'highlight', icon: 'star-o' },
  { name: '活动详情', value: 'detail', icon: 'description' },
  { name: '分割线', value: 'divider', icon: 'minus' },
  { name: '报名按钮', value: 'button', icon: 'plus' }
]

// 获取组件类型
const getComponentType = (type) => {
  const typeMap = {
    poster: PosterComponent,
    title: TitleComponent,
    time: TimeComponent,
    location: LocationComponent,
    highlight: HighlightComponent,
    detail: DetailComponent,
    divider: DividerComponent,
    button: ButtonComponent
  }
  return typeMap[type] || 'div'
}

// 获取编辑器类型
const getEditorType = (type) => {
  const typeMap = {
    poster: PosterEditor,
    title: TitleEditor,
    time: TimeEditor,
    location: LocationEditor,
    highlight: HighlightEditor,
    detail: DetailEditor,
    divider: DividerEditor,
    button: ButtonEditor
  }
  return typeMap[type]
}

// 获取组件名称
const getComponentName = (type) => {
  const component = componentLibrary.find(c => c.value === type)
  return component ? component.name : '组件'
}

// 获取默认数据
const getDefaultData = (type) => {
  const defaults = {
    poster: { imageUrl: '', alt: '活动海报' },
    title: { title: '活动标题', subtitle: '', color: '#333333' },
    time: { startTime: '', endTime: '', format: 'YYYY-MM-DD HH:mm' },
    location: { address: '', detail: '', mapUrl: '' },
    highlight: { items: [{ icon: 'star', text: '亮点1' }] },
    detail: { content: '活动详情内容' },
    divider: { style: 'solid', color: '#e5e5e5', height: 1 },
    button: { text: '立即报名', color: '#667eea', action: 'register' }
  }
  return defaults[type] || {}
}

// 选择组件
const onSelectComponent = (action) => {
  const component = {
    id: `${action.value}_${Date.now()}`,
    type: action.value,
    order: components.value.length,
    data: getDefaultData(action.value),
    style: {}
  }
  components.value.push(component)
  Toast.success(`已添加${action.name}`)
  showComponentLibrary.value = false
}

// 编辑组件
const editComponent = (component) => {
  editingComponent.value = { ...component }
  showEditor.value = true
}

// 保存组件数据
const saveComponentData = () => {
  const index = components.value.findIndex(c => c.id === editingComponent.value.id)
  if (index > -1) {
    components.value[index] = { ...editingComponent.value }
    Toast.success('保存成功')
  }
  closeEditor()
}

// 关闭编辑器
const closeEditor = () => {
  showEditor.value = false
  editingComponent.value = null
}

// 移动组件
const moveComponent = (index, direction) => {
  if (direction === 'up' && index > 0) {
    [components.value[index], components.value[index - 1]] = 
      [components.value[index - 1], components.value[index]]
  } else if (direction === 'down' && index < components.value.length - 1) {
    [components.value[index], components.value[index + 1]] = 
      [components.value[index + 1], components.value[index]]
  }
  // 更新order
  components.value.forEach((c, i) => c.order = i)
}

// 删除组件
const removeComponent = (id) => {
  Dialog.confirm({
    title: '确认删除',
    message: '确定要删除这个组件吗？'
  }).then(() => {
    const index = components.value.findIndex(c => c.id === id)
    if (index > -1) {
      components.value.splice(index, 1)
      Toast.success('删除成功')
    }
  }).catch(() => {
    // 取消删除
  })
}

// 预览页面
const previewPage = () => {
  if (components.value.length === 0) {
    Toast('请先添加组件')
    return
  }
  showPreview.value = true
}

// 保存页面
const savePage = async () => {
  if (components.value.length === 0) {
    Toast('请先添加组件')
    return
  }

  Toast.loading({ message: '保存中...', forbidClick: true, duration: 0 })
  
  try {
    const pageConfig = {
      components: components.value,
      theme: {
        primaryColor: '#667eea',
        backgroundColor: '#ffffff',
        textColor: '#333333'
      }
    }
    
    await campaignApi.savePageConfig(route.params.id, pageConfig)
    
    Toast.success('保存成功')
    setTimeout(() => {
      router.push('/brand/campaigns')
    }, 1000)
  } catch (error) {
    console.error('保存失败:', error)
    Toast.fail(error.message || '保存失败')
  }
}

// 加载页面配置
const loadPageConfig = async () => {
  loading.value = true
  Toast.loading({ message: '加载中...', forbidClick: true, duration: 0 })
  
  try {
    const config = await campaignApi.getPageConfig(route.params.id)
    
    if (config && config.components) {
      components.value = config.components
      Toast.success('加载成功')
    } else {
      Toast.clear()
    }
  } catch (error) {
    console.error('加载页面配置失败:', error)
    // 如果是404错误，说明还没有配置，不显示错误
    if (error.status !== 404) {
      Toast.fail('加载失败')
    } else {
      Toast.clear()
    }
  } finally {
    loading.value = false
  }
}

// 返回
const goBack = () => {
  if (components.value.length > 0) {
    Dialog.confirm({
      title: '确认返回',
      message: '有未保存的更改，确定要返回吗？'
    }).then(() => {
      router.back()
    }).catch(() => {
      // 取消返回
    })
  } else {
    router.back()
  }
}

// 组件挂载时加载配置
onMounted(() => {
  loadPageConfig()
})
</script>

<style scoped>
.page-designer {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 70px;
}

.design-canvas {
  min-height: calc(100vh - 120px);
  padding: 16px;
}

.component-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.canvas-component {
  background: white;
  border-radius: 8px;
  padding: 16px;
  position: relative;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.component-actions {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  gap: 4px;
  background: rgba(255, 255, 255, 0.9);
  padding: 4px;
  border-radius: 4px;
}

.bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px 16px;
  background: white;
  border-top: 1px solid #ebedf0;
  z-index: 100;
}

.component-editor {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.editor-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.preview-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: white;
}

.preview-content {
  flex: 1;
  overflow-y: auto;
}

.preview-component {
  padding: 16px;
  border-bottom: 1px solid #ebedf0;
}

/* 自定义Vant样式 */
:deep(.van-nav-bar) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

:deep(.van-nav-bar .van-nav-bar__title),
:deep(.van-nav-bar .van-nav-bar__text),
:deep(.van-nav-bar .van-icon) {
  color: white;
}
</style>