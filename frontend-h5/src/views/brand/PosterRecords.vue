<template>
  <div class="poster-records">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <h1 class="nav-title">海报生成记录</h1>
    </div>

    <!-- 筛选器 -->
    <div class="filters">
      <div class="filter-tabs">
        <button
          v-for="type in TYPE_TABS"
          :key="type.value"
          @click="currentType = type.value"
          :class="['filter-tab', { active: currentType === type.value }]"
        >
          {{ type.label }}
        </button>
      </div>

      <div class="date-filter">
        <input
          v-model="dateRange.start"
          type="date"
          class="date-input"
          placeholder="开始日期"
        >
        <span class="date-separator">至</span>
        <input
          v-model="dateRange.end"
          type="date"
          class="date-input"
          placeholder="结束日期"
        >
      </div>

      <div class="search-filter">
        <input
          v-model="searchKeyword"
          type="text"
          class="search-input"
          placeholder="搜索活动名称或分销员"
        >
      </div>
    </div>

    <!-- 海报记录列表 -->
    <div class="records-list">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="filteredRecords.length === 0" class="empty-state">
        <div class="empty-icon">🖼️</div>
        <p class="empty-text">暂无海报生成记录</p>
      </div>

      <div v-else class="record-cards">
        <div
          v-for="record in filteredRecords"
          :key="record.id"
          class="record-card"
        >
          <div class="card-header">
            <div class="record-info">
              <h3 class="poster-title">{{ record.campaignName || '活动海报' }}</h3>
              <span class="record-time">{{ formatDateTime(record.createdAt) }}</span>
            </div>
            <span class="type-badge">{{ getRecordTypeLabel(record.recordType) }}</span>
          </div>

          <div class="card-content">
            <div class="poster-preview">
              <img
                v-if="record.posterUrl"
                :src="record.posterUrl"
                :alt="record.campaignName"
                class="poster-image"
                @click="previewPoster(record.posterUrl)"
              >
              <div v-else class="no-preview">
                <span class="no-preview-icon">📷</span>
              </div>
            </div>

            <div class="record-details">
              <div class="detail-row">
                <span class="detail-label">模板:</span>
                <span class="detail-value">{{ record.templateName || '-' }}</span>
              </div>
              <div class="detail-row" v-if="record.recordType === 'distributor'">
                <span class="detail-label">分销员:</span>
                <span class="detail-value">{{ record.distributorName || '-' }}</span>
              </div>
              <div class="detail-row" v-if="record.recordType === 'campaign'">
                <span class="detail-label">活动ID:</span>
                <span class="detail-value">{{ record.campaignId }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">文件大小:</span>
                <span class="detail-value">{{ record.fileSize || '-' }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">生成耗时:</span>
                <span class="detail-value">{{ record.generationTime || '-' }}</span>
              </div>
              <div class="detail-row" v-if="record.downloadCount">
                <span class="detail-label">下载次数:</span>
                <span class="detail-value">{{ record.downloadCount }}</span>
              </div>
            </div>

            <div class="card-actions">
              <button
                v-if="record.posterUrl"
                @click="downloadPoster(record.posterUrl, record.campaignName)"
                class="action-btn download"
              >
                下载海报
              </button>
              <button
                @click="regeneratePoster(record.id)"
                :disabled="record.regenerating"
                class="action-btn regenerate"
              >
                {{ getRegenerateButtonText(record.regenerating) }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="stats-section">
      <h2 class="stats-title">生成统计</h2>
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-number">{{ recordStats.total }}</div>
          <div class="stat-label">总生成数</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ recordStats.campaign }}</div>
          <div class="stat-label">活动海报</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ recordStats.distributor }}</div>
          <div class="stat-label">分销商海报</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ recordStats.today }}</div>
          <div class="stat-label">今日生成</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ recordStats.totalDownloads }}</div>
          <div class="stat-label">总下载次数</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Toast } from 'vant'
import { posterApi } from '../../services/brandApi.js'
import {
  TYPE_TABS,
  formatDateTime,
  getRecordTypeLabel,
  filterRecords,
  calculateStats,
  buildPosterFileName,
  getDefaultDateRange,
  getDefaultStats,
  getRegenerateButtonText
} from './posterRecords.logic.js'

const router = useRouter()

const records = ref([])
const loading = ref(false)
const currentType = ref('all')
const searchKeyword = ref('')
const dateRange = ref(getDefaultDateRange())

const recordStats = ref(getDefaultStats())

const filteredRecords = computed(() => {
  return filterRecords(records.value, currentType.value, searchKeyword.value)
})

const loadRecords = async () => {
  loading.value = true
  try {
    const resp = await posterApi.getPosterRecords()
    if (resp && resp.records) {
      records.value = resp.records
      recordStats.value = calculateStats(records.value)
    }
  } catch (error) {
    console.error('加载海报记录失败:', error)
    Toast.fail('加载失败，请重试')
  } finally {
    loading.value = false
  }
}

const downloadPoster = async (url, campaignName) => {
  try {
    const link = document.createElement('a')
    link.href = url
    link.download = buildPosterFileName(campaignName)
    link.target = '_blank'
    link.click()
    Toast.success('开始下载海报')
    
    // 更新下载次数
    const record = records.value.find(r => r.posterUrl === url)
    if (record) {
      record.downloadCount = (record.downloadCount || 0) + 1
    }
  } catch (error) {
    console.error('下载海报失败:', error)
    Toast.fail('下载失败，请重试')
  }
}

const previewPoster = (url) => {
  window.open(url, '_blank')
}

const regeneratePoster = async (recordId) => {
  try {
    const record = records.value.find(r => r.id === recordId)
    if (record) {
      record.regenerating = true
    }
    
    await posterApi.regeneratePoster(recordId)
    Toast.success('海报重新生成成功')
    await loadRecords()
  } catch (error) {
    console.error('重新生成海报失败:', error)
    Toast.fail('重新生成失败，请重试')
    const record = records.value.find(r => r.id === recordId)
    if (record) {
      record.regenerating = false
    }
  }
}

onMounted(() => {
  loadRecords()
})
</script>

<style scoped>
.poster-records {
  min-height: 100vh;
  background: #f5f7fa;
}

.top-nav {
  background: white;
  padding: 16px;
  border-bottom: 1px solid #eee;
}

.nav-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.filters {
  background: white;
  padding: 16px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  border-bottom: 1px solid #eee;
}

.filter-tabs {
  display: flex;
  gap: 8px;
}

.filter-tab {
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background: white;
  color: #6b7280;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-tab.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.filter-tab:hover:not(.active) {
  background: #f3f4f6;
}

.date-filter {
  display: flex;
  align-items: center;
  gap: 8px;
}

.date-input {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  width: 150px;
}

.date-separator {
  color: #9ca3af;
}

.search-filter {
  flex: 1;
}

.search-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
}

.records-list {
  padding: 16px;
  max-width: 800px;
  margin: 0 auto;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #6b7280;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #9ca3af;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 14px;
}

.record-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.record-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f3f4f6;
}

.record-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.poster-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.record-time {
  font-size: 12px;
  color: #9ca3af;
}

.type-badge {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.type-badge.campaign {
  background: #dbeafe;
  color: #1e40af;
}

.type-badge.distributor {
  background: #fce7f3;
  color: #c2410c;
}

.card-content {
  padding: 16px;
}

.poster-preview {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

.poster-image {
  width: 200px;
  height: 200px;
  object-fit: contain;
  border-radius: 8px;
  border: 2px solid #e5e7eb;
  cursor: pointer;
  transition: all 0.2s;
}

.poster-image:hover {
  transform: scale(1.05);
}

.no-preview {
  width: 200px;
  height: 200px;
  background: #f8fafc;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.no-preview-icon {
  font-size: 48px;
  color: #cbd5e1;
}

.record-details {
  display: grid;
  gap: 12px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f9fafb;
}

.detail-label {
  font-size: 14px;
  color: #6b7280;
}

.detail-value {
  font-size: 14px;
  color: #1f2937;
  font-weight: 500;
}

.card-actions {
  display: flex;
  gap: 8px;
  padding: 16px;
  border-top: 1px solid #f3f4f6;
}

.action-btn {
  flex: 1;
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn.download {
  background: #10b981;
  color: white;
}

.action-btn.download:hover {
  background: #059669;
}

.action-btn.regenerate {
  background: #667eea;
  color: white;
}

.action-btn.regenerate:hover {
  background: #5a67d8;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.stats-section {
  max-width: 800px;
  margin: 20px auto;
  padding: 20px;
  background: white;
  border-radius: 12px;
}

.stats-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 16px 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-item {
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  text-align: center;
}

.stat-number {
  font-size: 24px;
  font-weight: 600;
  color: #667eea;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
}
</style>
