<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RefreshCw, Image, Download, Share, Clock, BarChart3 } from 'lucide-vue-next';
import { posterApi, type PosterRecord } from '../services/posterApi';

const records = ref<PosterRecord[]>([]);
const total = ref(0);
const loading = ref(false);

const fetchRecords = async () => {
  loading.value = true;
  try {
    const response = await posterApi.getPosterRecords();
    records.value = response.records;
    total.value = response.total;
  } catch (error) {
    console.error('Failed to fetch poster records:', error);
    alert('加载海报记录失败');
  } finally {
    loading.value = false;
  }
};

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN');
};

const getRecordTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    campaign: '活动海报',
    distributor: '分销员海报',
  };
  return typeMap[type] || type;
};

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    success: '成功',
    failed: '失败',
    pending: '生成中',
  };
  return statusMap[status] || status;
};

const getStatusClass = (status: string) => {
  const classMap: Record<string, string> = {
    success: 'bg-green-100 text-green-800',
    failed: 'bg-red-100 text-red-800',
    pending: 'bg-yellow-100 text-yellow-800',
  };
  return classMap[status] || '';
};

const formatFileSize = (size: string) => {
  const num = parseFloat(size);
  if (num < 1024) return `${num} B`;
  if (num < 1024 * 1024) return `${(num / 1024).toFixed(2)} KB`;
  return `${(num / (1024 * 1024)).toFixed(2)} MB`;
};

const handlePreview = (url: string) => {
  window.open(url, '_blank');
};

const handleDownload = (url: string, templateName: string) => {
  const link = document.createElement('a');
  link.href = url;
  link.download = `${templateName}_poster.jpg`;
  link.click();
};

onMounted(() => {
  fetchRecords();
});
</script>

<template>
  <div class="poster-records-view">
    <div class="header">
      <h1 class="title">海报生成记录</h1>
      <button class="btn-refresh" @click="fetchRecords">
        <RefreshCw :size="20" :class="{ spinning: loading }" />
        <span>刷新</span>
      </button>
    </div>

    <div class="stats">
      <div class="stat-card">
        <div class="stat-icon" style="background: #3b82f6;">
          <Image :size="24" />
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ total }}</div>
          <div class="stat-label">海报生成总数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #10b981;">
          <Download :size="24" />
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ records.reduce((sum, r) => sum + r.downloadCount, 0) }}</div>
          <div class="stat-label">总下载次数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #8b5cf6;">
          <Share :size="24" />
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ records.reduce((sum, r) => sum + r.shareCount, 0) }}</div>
          <div class="stat-label">总分享次数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #f59e0b;">
          <Clock :size="24" />
        </div>
        <div class="stat-content">
          <div class="stat-value">
            {{ Math.round(records.reduce((sum, r) => sum + r.generationTime, 0) / (records.length || 1)) }}ms
          </div>
          <div class="stat-label">平均生成时间</div>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="table-container">
      <table class="records-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>类型</th>
            <th>活动/分销员ID</th>
            <th>模板名称</th>
            <th>海报预览</th>
            <th>文件大小</th>
            <th>生成时间</th>
            <th>下载/分享</th>
            <th>状态</th>
            <th>操作</th>
            <th>创建时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in records" :key="record.id">
            <td>{{ record.id }}</td>
            <td>{{ getRecordTypeText(record.recordType) }}</td>
            <td class="id-cell">
              {{ record.recordType === 'campaign' ? record.campaignId : record.distributorId }}
            </td>
            <td>{{ record.templateName }}</td>
            <td class="preview-cell">
              <img
                :src="record.thumbnailUrl || record.posterUrl"
                :alt="record.templateName"
                class="thumbnail"
                @click="handlePreview(record.posterUrl)"
              />
            </td>
            <td>{{ formatFileSize(record.fileSize) }}</td>
            <td>{{ record.generationTime }}ms</td>
            <td class="counts">
              <span class="count-item">
                <Download :size="14" />
                {{ record.downloadCount }}
              </span>
              <span class="count-item">
                <Share :size="14" />
                {{ record.shareCount }}
              </span>
            </td>
            <td>
              <span class="status-badge" :class="getStatusClass(record.status)">
                {{ getStatusText(record.status) }}
              </span>
            </td>
            <td class="actions">
              <button
                class="btn-action"
                @click="handlePreview(record.posterUrl)"
                title="预览"
              >
                <Image :size="16" />
              </button>
              <button
                class="btn-action"
                @click="handleDownload(record.posterUrl, record.templateName)"
                title="下载"
              >
                <Download :size="16" />
              </button>
            </td>
            <td>{{ formatDate(record.createdAt) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.poster-records-view {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
}

.btn-refresh {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background 0.2s;
}

.btn-refresh:hover {
  background: #2563eb;
}

.btn-refresh .spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: white;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a1a;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  margin-top: 4px;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #6b7280;
}

.table-container {
  background: white;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  overflow-x: auto;
}

.records-table {
  width: 100%;
  border-collapse: collapse;
}

.records-table th {
  background: #f9fafb;
  padding: 12px;
  text-align: left;
  font-size: 13px;
  font-weight: 600;
  color: #6b7280;
  border-bottom: 1px solid #e5e7eb;
  white-space: nowrap;
}

.records-table td {
  padding: 12px;
  border-bottom: 1px solid #f3f4f6;
  font-size: 14px;
  color: #374151;
}

.id-cell {
  font-family: monospace;
  color: #3b82f6;
  font-weight: 500;
}

.preview-cell {
  padding: 8px;
}

.thumbnail {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 6px;
  cursor: pointer;
  transition: transform 0.2s;
}

.thumbnail:hover {
  transform: scale(1.05);
}

.counts {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.count-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #6b7280;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn-action {
  padding: 6px;
  border: none;
  background: #f3f4f6;
  cursor: pointer;
  color: #6b7280;
  border-radius: 4px;
  transition: all 0.2s;
}

.btn-action:hover {
  background: #e5e7eb;
  color: #3b82f6;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}
</style>
