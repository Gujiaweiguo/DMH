<template>
  <div class="member-detail-page">
    <van-nav-bar title="会员详情" left-arrow @click-left="$router.back()" />
    
    <van-loading v-if="loading" class="loading" />
    
    <div v-else-if="member" class="content">
      <!-- 基本信息 -->
      <div class="section">
        <div class="member-header">
          <van-image
            round
            width="80"
            height="80"
            :src="member.avatar || '/default-avatar.png'"
          />
          <div class="member-basic">
            <div class="member-name">{{ member.nickname || '未设置昵称' }}</div>
            <div class="member-id">ID: {{ member.id }}</div>
            <van-tag :type="member.status === 'active' ? 'success' : 'danger'">
              {{ member.status === 'active' ? '正常' : '禁用' }}
            </van-tag>
          </div>
        </div>
        
        <van-cell-group>
          <van-cell title="手机号" :value="member.phone || '未绑定'" />
          <van-cell title="性别" :value="getGenderText(member.gender)" />
          <van-cell title="来源" :value="member.source || '未知'" />
          <van-cell title="注册时间" :value="formatDateTime(member.createdAt)" />
        </van-cell-group>
      </div>

      <!-- 消费统计 -->
      <div class="section">
        <div class="section-title">消费统计</div>
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-value">{{ member.totalOrders }}</div>
            <div class="stat-label">订单数</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">¥{{ member.totalPayment.toFixed(2) }}</div>
            <div class="stat-label">消费金额</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">¥{{ member.totalReward.toFixed(2) }}</div>
            <div class="stat-label">奖励金额</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ member.participatedCampaigns }}</div>
            <div class="stat-label">参与活动</div>
          </div>
        </div>
        
        <van-cell-group>
          <van-cell title="首次下单" :value="formatDateTime(member.firstOrderAt)" />
          <van-cell title="最后下单" :value="formatDateTime(member.lastOrderAt)" />
        </van-cell-group>
      </div>

      <!-- 标签 -->
      <div class="section" v-if="member.tags && member.tags.length > 0">
        <div class="section-title">会员标签</div>
        <div class="tags-container">
          <van-tag
            v-for="tag in member.tags"
            :key="tag.id"
            :color="tag.color"
            size="medium"
            style="margin: 5px;"
          >
            {{ tag.name }}
          </van-tag>
        </div>
      </div>

      <!-- 关联品牌 -->
      <div class="section" v-if="member.brands && member.brands.length > 0">
        <div class="section-title">关联品牌</div>
        <van-cell-group>
          <van-cell
            v-for="brand in member.brands"
            :key="brand.brandId"
            :title="brand.brandName"
            :label="`首次参与: ${formatDateTime(brand.createdAt)}`"
          />
        </van-cell-group>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { showFailToast } from 'vant';

const route = useRoute();
const member = ref(null);
const loading = ref(true);

// 加载会员详情
const loadMemberDetail = async () => {
  loading.value = true;
  try {
    const token = localStorage.getItem('dmh_token');
    const memberId = route.params.id;
    
    const response = await fetch(`/api/v1/members/${memberId}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (!response.ok) throw new Error('加载失败');
    
    member.value = await response.json();
  } catch (error) {
    showFailToast(error.message || '加载失败');
  } finally {
    loading.value = false;
  }
};

// 格式化日期时间
const formatDateTime = (dateStr) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
};

// 获取性别文本
const getGenderText = (gender) => {
  const genderMap = {
    0: '未知',
    1: '男',
    2: '女',
  };
  return genderMap[gender] || '未知';
};

onMounted(() => {
  loadMemberDetail();
});
</script>

<style scoped>
.member-detail-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}

.content {
  padding-bottom: 20px;
}

.section {
  margin-bottom: 10px;
  background: white;
}

.section-title {
  padding: 15px;
  font-size: 16px;
  font-weight: bold;
  border-bottom: 1px solid #f0f0f0;
}

.member-header {
  display: flex;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.member-basic {
  flex: 1;
  margin-left: 15px;
}

.member-name {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 8px;
}

.member-id {
  font-size: 14px;
  color: #999;
  margin-bottom: 8px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.stat-card {
  text-align: center;
  padding: 15px;
  background: #f8f8f8;
  border-radius: 8px;
}

.stat-value {
  font-size: 20px;
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 13px;
  color: #666;
}

.tags-container {
  padding: 15px;
}
</style>
