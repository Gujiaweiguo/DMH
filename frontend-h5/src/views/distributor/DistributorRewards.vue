<template>
  <div class="distributor-rewards">
    <van-nav-bar title="奖励明细" left-arrow @click-left="$router.back()" />

    <div class="rewards-content">
      <!-- 筛选条件 -->
      <van-dropdown-menu>
        <van-dropdown-item v-model="filterLevel" :options="levelOptions" @change="loadRewards" />
        <van-dropdown-item v-model="filterTime" :options="timeOptions" @change="loadRewards" />
      </van-dropdown-menu>

      <!-- 奖励列表 -->
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-list
          v-model:loading="loading"
          :finished="finished"
          finished-text="没有更多了"
          @load="loadRewards"
        >
          <van-cell
            v-for="reward in rewards"
            :key="reward.id"
            class="reward-item"
          >
            <template #title>
              <div class="reward-title">
                <span class="reward-level">{{ reward.level }}级奖励</span>
                <span class="reward-rate">{{ reward.rewardRate }}%</span>
              </div>
            </template>
            <template #value>
              <span class="reward-amount">+¥{{ reward.amount.toFixed(2) }}</span>
            </template>
            <template #label>
              <div class="reward-info">
                <p>订单号: {{ reward.orderId }}</p>
                <p v-if="reward.fromUsername">来自: {{ reward.fromUsername }}</p>
                <p>时间: {{ reward.createdAt }}</p>
              </div>
            </template>
          </van-cell>
        </van-list>
      </van-pull-refresh>

      <!-- 空状态 -->
      <van-empty v-if="!loading && rewards.length === 0" description="暂无奖励记录" />
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Toast } from 'vant'
import axios from '@/utils/axios'
import {
  PAGE_SIZE,
  getLevelOptions,
  getTimeOptions,
  calculateDateRange,
  getDefaultFilters,
  shouldFinishLoading,
  mergeRewardsList
} from './distributorRewards.logic.js'

export default {
  name: 'DistributorRewards',
  setup() {
    const route = useRoute()
    const brandId = ref(parseInt(route.query.brandId) || 0)
    const rewards = ref([])
    const loading = ref(false)
    const finished = ref(false)
    const refreshing = ref(false)
    const page = ref(1)
    const filterLevel = ref(getDefaultFilters().level)
    const filterTime = ref(getDefaultFilters().time)

    const levelOptions = getLevelOptions()
    const timeOptions = getTimeOptions()

    // 加载奖励列表
    const loadRewards = async () => {
      if (refreshing.value) {
        rewards.value = []
        page.value = 1
        finished.value = false
      }

      loading.value = true
      try {
        const params = {
          page: page.value,
          pageSize: PAGE_SIZE
        }

        if (filterLevel.value > 0) {
          params.level = filterLevel.value
        }

        const dateRange = calculateDateRange(filterTime.value)
        if (dateRange) {
          params.startDate = dateRange.startDate
          params.endDate = dateRange.endDate
        }

		const data = await axios.get(`/distributor/rewards/${brandId.value}`, { params })

        if (data.code === 200) {
          const newRewards = data.data.rewards || []
          rewards.value = mergeRewardsList(rewards.value, newRewards, refreshing.value)

          if (shouldFinishLoading(newRewards.length, PAGE_SIZE)) {
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
      loadRewards()
    }

    onMounted(() => {
      loadRewards()
    })

    return {
      rewards,
      loading,
      finished,
      refreshing,
      filterLevel,
      filterTime,
      levelOptions,
      timeOptions,
      loadRewards,
      onRefresh
    }
  }
}
</script>

<style scoped>
.distributor-rewards {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.rewards-content {
  padding-top: 0;
}

.reward-item {
  margin-bottom: 8px;
}

.reward-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.reward-level {
  font-size: 14px;
  color: #666;
}

.reward-rate {
  font-size: 12px;
  color: #999;
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 4px;
}

.reward-amount {
  font-size: 16px;
  font-weight: bold;
  color: #f56c6c;
}

.reward-info p {
  margin: 4px 0;
  font-size: 12px;
  color: #999;
}
</style>
