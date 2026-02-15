<template>
  <div class="members-page">
    <van-nav-bar title="会员管理" left-arrow @click-left="$router.back()" />
    
    <!-- 搜索和筛选 -->
    <div class="search-section">
      <van-search
        v-model="searchKeyword"
        placeholder="搜索昵称/手机号"
        @search="handleSearch"
      />
      
      <van-dropdown-menu>
        <van-dropdown-item v-model="filters.source" :options="sourceOptions" @change="handleSearch" />
        <van-dropdown-item v-model="filters.status" :options="statusOptions" @change="handleSearch" />
      </van-dropdown-menu>
    </div>

    <!-- 会员列表 -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="loadMembers"
      >
        <div v-for="member in members" :key="member.id" class="member-card" @click="viewDetail(member.id)">
          <div class="member-header">
            <van-image
              round
              width="50"
              height="50"
              :src="member.avatar || '/default-avatar.png'"
            />
            <div class="member-info">
              <div class="member-name">{{ member.nickname || '未设置昵称' }}</div>
              <div class="member-phone">{{ member.phone || '未绑定手机' }}</div>
            </div>
            <van-tag :type="member.status === 'active' ? 'success' : 'danger'">
              {{ member.status === 'active' ? '正常' : '禁用' }}
            </van-tag>
          </div>
          
          <div class="member-stats">
            <div class="stat-item">
              <div class="stat-value">{{ member.totalOrders }}</div>
              <div class="stat-label">订单数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">¥{{ member.totalPayment.toFixed(2) }}</div>
              <div class="stat-label">消费金额</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ member.participatedCampaigns }}</div>
              <div class="stat-label">参与活动</div>
            </div>
          </div>
          
          <div class="member-footer">
            <span class="member-source">来源: {{ member.source || '未知' }}</span>
            <span class="member-time">注册: {{ formatDate(member.createdAt) }}</span>
          </div>
        </div>
      </van-list>
    </van-pull-refresh>

    <!-- 底部操作栏 -->
    <div class="bottom-actions">
      <van-button type="primary" block @click="showExportDialog">申请导出会员数据</van-button>
    </div>

    <!-- 导出申请弹窗 -->
    <van-dialog
      v-model:show="exportDialogVisible"
      title="申请导出会员数据"
      show-cancel-button
      @confirm="submitExportRequest"
    >
      <van-field
        v-model="exportReason"
        rows="3"
        autosize
        type="textarea"
        placeholder="请说明导出原因（必填）"
      />
    </van-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { showToast, showSuccessToast, showFailToast } from 'vant';
import {
  buildExportRequestBody,
  buildMembersQueryParams,
  formatMemberDate,
  getBrandIdFromStorage,
  mergeMembersPage,
} from './members.logic.js';

const router = useRouter();

const members = ref([]);
const loading = ref(false);
const finished = ref(false);
const refreshing = ref(false);
const currentPage = ref(1);
const pageSize = 20;

const searchKeyword = ref('');
const filters = reactive({
  source: '',
  status: '',
});

const sourceOptions = [
  { text: '全部来源', value: '' },
  { text: '微信', value: 'wechat' },
  { text: '支付宝', value: 'alipay' },
  { text: '其他', value: 'other' },
];

const statusOptions = [
  { text: '全部状态', value: '' },
  { text: '正常', value: 'active' },
  { text: '禁用', value: 'disabled' },
];

const exportDialogVisible = ref(false);
const exportReason = ref('');

const getBrandId = () => {
  return getBrandIdFromStorage(localStorage.getItem('brandId'), localStorage.getItem('dmh_user_info'));
};

// 加载会员列表
const loadMembers = async () => {
  if (loading.value || finished.value) return;
  
  loading.value = true;
  try {
    const token = localStorage.getItem('dmh_token');
    const brandId = getBrandId();
    
    const params = buildMembersQueryParams({
      page: currentPage.value,
      pageSize,
      brandId,
      keyword: searchKeyword.value,
      source: filters.source,
      status: filters.status,
    });

    const response = await fetch(`/api/v1/members?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (!response.ok) throw new Error('加载失败');
    
    const data = await response.json();
    
    const merged = mergeMembersPage({
      currentMembers: members.value,
      incomingMembers: data.members,
      page: currentPage.value,
      pageSize,
    });
    members.value = merged.members;
    finished.value = merged.finished;
    currentPage.value = merged.nextPage;
  } catch (error) {
    showFailToast(error.message || '加载失败');
  } finally {
    loading.value = false;
    refreshing.value = false;
  }
};

// 刷新
const onRefresh = () => {
  currentPage.value = 1;
  finished.value = false;
  members.value = [];
  loadMembers();
};

// 搜索
const handleSearch = () => {
  onRefresh();
};

// 查看详情
const viewDetail = (memberId) => {
  router.push(`/brand/members/${memberId}`);
};

// 显示导出对话框
const showExportDialog = () => {
  exportReason.value = '';
  exportDialogVisible.value = true;
};

// 提交导出申请
const submitExportRequest = async () => {
  if (!exportReason.value.trim()) {
    showToast('请填写导出原因');
    return;
  }

  try {
    const token = localStorage.getItem('dmh_token');
    const brandId = getBrandId();
    if (!brandId) {
      showFailToast('未找到品牌信息，无法提交导出申请');
      return;
    }
    
    const response = await fetch('/api/v1/members/export-requests', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(buildExportRequestBody({
        brandId,
        reason: exportReason.value,
        keyword: searchKeyword.value,
        source: filters.source,
        status: filters.status,
      })),
    });

    if (!response.ok) throw new Error('申请失败');
    
    showSuccessToast('导出申请已提交，请等待审批');
    exportDialogVisible.value = false;
  } catch (error) {
    showFailToast(error.message || '申请失败');
  }
};

// 格式化日期
const formatDate = (dateStr) => {
  return formatMemberDate(dateStr);
};

onMounted(() => {
  loadMembers();
});
</script>

<style scoped>
.members-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 70px;
}

.search-section {
  background: white;
  margin-bottom: 10px;
}

.member-card {
  background: white;
  margin: 10px;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.member-header {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.member-info {
  flex: 1;
  margin-left: 12px;
}

.member-name {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 4px;
}

.member-phone {
  font-size: 14px;
  color: #666;
}

.member-stats {
  display: flex;
  justify-content: space-around;
  padding: 15px 0;
  border-top: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.member-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  font-size: 12px;
  color: #999;
}

.bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 10px;
  background: white;
  box-shadow: 0 -2px 4px rgba(0,0,0,0.1);
}
</style>
