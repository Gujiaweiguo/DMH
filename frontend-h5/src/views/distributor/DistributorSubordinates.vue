<template>
  <div class="distributor-subordinates">
    <van-nav-bar title="我的下级" left-arrow @click-left="$router.back()" />

    <div class="subordinates-content">
      <!-- 统计信息 -->
      <div class="stats-header">
        <div class="stat-item">
          <div class="stat-value">{{ subordinates.length }}</div>
          <div class="stat-label">下级人数</div>
        </div>
      </div>

      <!-- 下级列表 -->
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-list
          v-model:loading="loading"
          :finished="finished"
          finished-text="没有更多了"
          @load="loadSubordinates"
        >
          <van-cell
            v-for="sub in subordinates"
            :key="sub.id"
            class="subordinate-item"
            is-link
          >
            <template #title>
              <div class="subordinate-name">{{ sub.username || '用户' }}{{ sub.id }}</div>
            </template>
            <template #label>
              <div class="subordinate-info">
                <p>{{ sub.level }}级分销商 · 加入时间: {{ sub.createdAt }}</p>
              </div>
            </template>
            <template #value>
              <div class="subordinate-stats">
                <p>订单: {{ sub.totalOrders }}</p>
                <p>收益: ¥{{ sub.totalEarnings.toFixed(2) }}</p>
              </div>
            </template>
          </van-cell>
        </van-list>
      </van-pull-refresh>

      <!-- 空状态 -->
      <van-empty v-if="!loading && subordinates.length === 0" description="暂无下级分销商">
        <template #image>
          <van-icon name="friends-o" size="80" color="#ddd" />
        </template>
        <template #description>
          <p>还没有下级分销商</p>
          <p class="tip">分享推广链接，邀请好友成为分销商</p>
        </template>
      </van-empty>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Toast } from 'vant'
import axios from '@/utils/axios'

export default {
  name: 'DistributorSubordinates',
  setup() {
    const route = useRoute()
    const brandId = ref(parseInt(route.query.brandId) || 0)
    const subordinates = ref([])
    const loading = ref(false)
    const finished = ref(false)
    const refreshing = ref(false)
    const page = ref(1)
    const pageSize = 20

    // 加载下级列表
    const loadSubordinates = async () => {
      if (refreshing.value) {
        subordinates.value = []
        page.value = 1
        finished.value = false
      }

      loading.value = true
      try {
		const data = await axios.get(`/distributor/subordinates/${brandId.value}`, {
          params: {
            page: page.value,
            pageSize: pageSize
          }
        })

        if (data.code === 200) {
          const newSubordinates = data.data.subordinates || []
          if (refreshing.value) {
            subordinates.value = newSubordinates
          } else {
            subordinates.value.push(...newSubordinates)
          }

          if (newSubordinates.length < pageSize) {
            finished.value = true
          } else {
            page.value++
          }
        }
      } catch (error) {
        Toast('加载失败')
      } finally {
        loading.value = false
        refreshing.value = false
      }
    }

    // 下拉刷新
    const onRefresh = () => {
      refreshing.value = true
      loadSubordinates()
    }

    onMounted(() => {
      loadSubordinates()
    })

    return {
      subordinates,
      loading,
      finished,
      refreshing,
      loadSubordinates,
      onRefresh
    }
  }
}
</script>

<style scoped>
.distributor-subordinates {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.subordinates-content {
  padding-top: 0;
}

.stats-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 24px;
  text-align: center;
}

.stat-item {
  display: inline-block;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

.subordinate-item {
  margin-bottom: 8px;
}

.subordinate-name {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.subordinate-info p {
  margin: 4px 0;
  font-size: 12px;
  color: #999;
}

.subordinate-stats {
  text-align: right;
}

.subordinate-stats p {
  margin: 2px 0;
  font-size: 12px;
  color: #666;
}

.tip {
  font-size: 14px;
  color: #999;
  margin-top: 8px;
}
</style>
