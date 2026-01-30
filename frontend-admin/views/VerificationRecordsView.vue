<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RefreshCw, Clock, User, CheckCircle, XCircle } from 'lucide-vue-next';
import { orderApi, type VerificationRecord } from '../services/orderApi';

const records = ref<VerificationRecord[]>([]);
const total = ref(0);
const loading = ref(false);

const fetchRecords = async () => {
  loading.value = true;
  try {
    const response = await orderApi.getVerificationRecords();
    records.value = response.records;
    total.value = response.total;
  } catch (error) {
    console.error('Failed to fetch verification records:', error);
    alert('加载核销记录失败');
  } finally {
    loading.value = false;
  }
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleString('zh-CN');
};

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    verified: '已核销',
    unverified: '未核销',
    cancelled: '已取消',
  };
  return statusMap[status] || status;
};

const getStatusClass = (status: string) => {
  const classMap: Record<string, string> = {
    verified: 'bg-green-100 text-green-800',
    unverified: 'bg-yellow-100 text-yellow-800',
    cancelled: 'bg-red-100 text-red-800',
  };
  return classMap[status] || '';
};

const getMethodText = (method: string) => {
  const methodMap: Record<string, string> = {
    scan: '扫码核销',
    manual: '手动核销',
  };
  return methodMap[method] || method;
};

onMounted(() => {
  fetchRecords();
});
</script>

<template>
  <div class="verification-records-view">
    <div class="header">
      <h1 class="title">核销记录查询</h1>
      <button class="btn-refresh" @click="fetchRecords">
        <RefreshCw :size="20" :class="{ spinning: loading }" />
        <span>刷新</span>
      </button>
    </div>

    <div class="stats">
      <div class="stat-card">
        <div class="stat-icon" style="background: #3b82f6;">
          <Clock :size="24" />
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ total }}</div>
          <div class="stat-label">核销记录总数</div>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else class="table-container">
      <table class="records-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>订单ID</th>
            <th>核销状态</th>
            <th>核销时间</th>
            <th>核销人ID</th>
            <th>核销码</th>
            <th>核销方式</th>
            <th>备注</th>
            <th>创建时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in records" :key="record.id">
            <td>{{ record.id }}</td>
            <td class="order-id">{{ record.orderId }}</td>
            <td>
              <span class="status-badge" :class="getStatusClass(record.verificationStatus)">
                {{ getStatusText(record.verificationStatus) }}
              </span>
            </td>
            <td>{{ formatDate(record.verifiedAt) }}</td>
            <td>{{ record.verifiedBy || '-' }}</td>
            <td class="code">{{ record.verificationCode }}</td>
            <td>{{ getMethodText(record.verificationMethod) }}</td>
            <td class="remark">{{ record.remark || '-' }}</td>
            <td>{{ formatDate(record.createdAt) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.verification-records-view {
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
}

.records-table td {
  padding: 12px;
  border-bottom: 1px solid #f3f4f6;
  font-size: 14px;
  color: #374151;
}

.order-id {
  font-family: monospace;
  color: #3b82f6;
  font-weight: 500;
}

.code {
  font-family: monospace;
  color: #6b7280;
  font-size: 13px;
}

.remark {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}
</style>
