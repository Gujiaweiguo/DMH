<template>
  <div class="poster-generator">
    <!-- 顶部导航栏 -->
    <van-nav-bar
      title="海报生成器"
      left-text="返回"
      left-arrow
      @click-left="goBack"
    />

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 模板选择区域 -->
      <div class="template-section">
        <div class="section-header">
          <h2>1. 选择海报模板</h2>
        </div>
        
        <van-loading v-if="state.loading.templates" size="24px">
          加载模板中...
        </van-loading>

        <div v-else class="template-grid">
          <div 
            v-for="template in state.templates" 
            :key="template.id"
            class="template-item"
            :class="{ 'selected': state.selectedTemplateId === template.id }"
            @click="selectTemplate(template)"
          >
            <div class="template-thumbnail">
              <van-image :src="template.thumbnailUrl" fit="cover" />
            </div>
            <div class="template-info">
              <div class="template-name">{{ template.name }}</div>
              <div class="template-desc">{{ template.desc }}</div>
            </div>
          </div>
        </div>

        <van-empty v-if="state.templates.length === 0" description="暂无模板">
          <van-button size="small" type="primary" @click="loadTemplates">
            加载模板
          </van-button>
        </van-empty>
      </div>

      <!-- 海报预览区域 -->
      <div class="preview-section">
        <div class="section-header">
          <h2>2. 预览海报效果</h2>
        </div>

        <div class="preview-container">
          <van-loading v-if="state.loading.preview" size="32px">
            生成预览中...
          </van-loading>

          <div v-else class="poster-preview" ref="previewContainer">
            <!-- 海报图片 -->
            <van-image 
              v-if="state.previewUrl" 
              :src="state.previewUrl" 
              fit="contain"
              class="poster-image"
              @click="showImageDetail"
            />
            
            <!-- 空状态 -->
            <van-empty v-if="!state.previewUrl" description="请选择模板后预览" />
            
            <!-- 缩放和旋转控制 -->
            <div v-if="state.previewUrl" class="preview-controls">
              <van-button-group>
                <van-button size="small" @click="zoomIn" icon="plus" />
                <van-button size="small" @click="zoomOut" icon="minus" />
                <van-button size="small" @click="resetZoom" icon="replay" />
              </van-button-group>
              
              <van-button-group>
                <van-button size="small" @click="rotateLeft" icon="arrow-left" />
                <van-button size="small" @click="rotateRight" icon="arrow" />
                <van-button size="small" @click="resetRotation" icon="replay" />
              </van-button-group>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-section">
        <div class="section-header">
          <h2>3. 生成海报</h2>
        </div>

        <div class="action-buttons">
          <van-button 
            type="primary" 
            size="large" 
            block
            :loading="state.loading.generate"
            :disabled="!state.selectedTemplateId || !state.previewUrl"
            @click="generatePoster"
          >
            生成海报
          </van-button>

          <van-button 
            type="success" 
            size="large" 
            block
            :disabled="!state.posterUrl"
            @click="downloadPoster"
            icon="down"
          >
            下载海报
          </van-button>

          <van-button 
            type="warning" 
            size="large" 
            block
            :disabled="!state.posterUrl"
            @click="sharePoster"
            icon="share"
          >
            分享海报
          </van-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Toast, Dialog, ImagePreview } from 'vant'
import { posterApi } from '../../services/brandApi'

const router = useRouter()

// 响应式状态
const state = reactive({
  loading: {
    templates: false,
    preview: false,
    generate: false
  },
  templates: [],
  selectedTemplateId: null,
  previewUrl: '',
  posterUrl: '',
  qrcodeUrl: '',
  scale: 1,
  rotation: 0
})

// 获取路由参数
const routeParams = router.currentRoute.value.params
const campaignId = routeParams.id

// 临时使用模拟模板数据（等后端API可用后移除）
const loadTemplates = async () => {
  state.loading.templates = true
  try {
    const response = await posterApi.getPosterTemplates()
    state.templates = response.templates || []
    if (state.templates.length > 0) {
      selectTemplate(state.templates[0])
    }
  } catch (error) {
    console.error('加载模板失败:', error)
    Toast.fail('加载模板失败')
  } finally {
    state.loading.templates = false
  }
}

// 选择模板
const selectTemplate = async (template) => {
  state.selectedTemplateId = template.id
  state.previewUrl = ''
  state.posterUrl = ''
  
  try {
    state.loading.preview = true
    state.previewUrl = template.previewImage
    state.loading.preview = false
  } catch (error) {
    state.loading.preview = false
    console.error('生成预览失败:', error)
    Toast.fail('生成预览失败')
  }
}

// 生成海报
const generatePoster = async () => {
  if (!state.selectedTemplateId) {
    Toast('请先选择模板')
    return
  }

  state.loading.generate = true
  try {
    const response = await posterApi.generatePoster(campaignId, 0)
    state.posterUrl = response.posterUrl
    Toast.success('海报生成成功')
  } catch (error) {
    console.error('生成海报失败:', error)
    Toast.fail('海报生成失败')
  } finally {
    state.loading.generate = false
  }
}

// 下载海报
const downloadPoster = () => {
  if (!state.posterUrl) {
    Toast('请先生成海报')
    return
  }

  try {
    // 创建下载链接
    const link = document.createElement('a')
    link.href = state.posterUrl
    link.download = `poster_${Date.now()}.png`
    link.click()
    Toast.success('开始下载')
  } catch (error) {
    console.error('下载失败:', error)
    Toast.fail('下载失败，请长按图片保存')
  }
}

// 分享海报
const sharePoster = () => {
  if (!state.posterUrl) {
    Toast('请先生成海报')
    return
  }

  try {
    // 检查微信环境
    const isWeixin = /micromessenger/i.test(navigator.userAgent)
    
    if (isWeixin) {
      // 微信环境，调用微信分享
      if (window.wx) {
        window.wx.updateAppMessageShareData({
          title: '活动海报',
          desc: '扫码参与活动',
          link: window.location.href,
          imgUrl: state.posterUrl,
          success: () => {
            Toast.success('分享成功')
          },
          fail: () => {
            Toast.fail('分享失败')
          }
        })
      } else {
        Toast.fail('微信分享功能不可用')
      }
    } else {
      // 非微信环境，复制链接
      const shareUrl = window.location.href

      // 尝试使用现代Clipboard API
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(shareUrl)
          .then(() => {
            Toast.success('链接已复制，可分享给好友')
          })
          .catch((err) => {
            console.error('Clipboard API失败:', err)
            // 降级到传统复制方法
            fallbackCopyText(shareUrl)
          })
      } else {
        // 降级到传统复制方法
        fallbackCopyText(shareUrl)
      }
    }
  } catch (error) {
    console.error('分享失败:', error)
    Toast.fail('分享失败')
  }
}

// 降级复制方法（兼容旧浏览器）
const fallbackCopyText = (text) => {
  try {
    const textArea = document.createElement('textarea')
    textArea.value = text
    textArea.style.position = 'fixed'
    textArea.style.left = '-9999px'
    textArea.style.top = '0'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()
    const successful = document.execCommand('copy')
    document.body.removeChild(textArea)
    if (successful) {
      Toast.success('链接已复制，可分享给好友')
    } else {
      Toast.fail('复制失败，请手动复制链接')
    }
  } catch (err) {
    console.error('降级复制失败:', err)
    Toast.fail('复制失败，请手动复制链接')
  }
}

// 显示图片详情
const showImageDetail = () => {
  ImagePreview({
    images: [state.previewUrl],
    startPosition: 0,
    closeable: true
  })
}

// 缩放控制
const zoomIn = () => {
  if (state.scale < 3) {
    state.scale += 0.25
  }
}

const zoomOut = () => {
  if (state.scale > 0.5) {
    state.scale -= 0.25
  }
}

const resetZoom = () => {
  state.scale = 1
}

// 旋转控制
const rotateLeft = () => {
  state.rotation -= 90
}

const rotateRight = () => {
  state.rotation += 90
}

const resetRotation = () => {
  state.rotation = 0
}

// 返回
const goBack = () => {
  router.back()
}

// 生命周期
onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.poster-generator {
  min-height: 100vh;
  background: #f5f5f5;
}

.main-content {
  padding-bottom: 80px;
}

.section-header {
  padding: 20px;
  background: white;
  border-bottom: 1px solid #e5e5e5;
}

.section-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #323233;
}

.template-section,
.preview-section,
.action-section {
  margin-bottom: 20px;
  background: white;
  border-radius: 12px;
  overflow: hidden;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 16px;
}

.template-item {
  border: 2px solid #e5e5e5;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
}

.template-item.selected {
  border-color: #1989fa;
  background: #f0f7ff;
}

.template-item:active {
  transform: scale(0.98);
}

.template-thumbnail {
  width: 100%;
  height: 120px;
  background: #f8f8fa;
  display: flex;
  align-items: center;
  justify-content: center;
}

.template-info {
  padding: 12px;
}

.template-name {
  font-size: 16px;
  font-weight: 600;
  color: #323233;
  margin-bottom: 4px;
}

.template-desc {
  font-size: 13px;
  color: #969796;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-container {
  background: #f8f8fa;
  border-radius: 8px;
  padding: 20px;
  min-height: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.poster-preview {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.poster-image {
  max-width: 100%;
  max-height: 400px;
  transition: transform 0.3s ease;
}

.preview-controls {
  margin-top: 20px;
  display: flex;
  gap: 8px;
  justify-content: center;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
}

.action-buttons .van-button {
  height: 48px;
  font-size: 16px;
  font-weight: 600;
}
</style>
