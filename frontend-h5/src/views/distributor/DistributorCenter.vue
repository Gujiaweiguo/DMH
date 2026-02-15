<template>
  <div class="distributor-center">
    <div class="header">
      <div class="header-top">
        <h1>分销中心</h1>
        <button class="logout-btn" type="button" @click="handleLogout">退出登录</button>
      </div>
      <div class="brand-selector" v-if="hasDistributor && brands.length > 1">
        <select v-model="selectedBrandId" @change="switchBrand">
          <option v-for="brand in brands" :key="brand.brandId" :value="brand.brandId">
            {{ brand.brandName }}
          </option>
        </select>
      </div>
    </div>

    <!-- 未成为分销商 -->
    <div class="not-distributor" v-if="!hasDistributor">
      <div class="empty-state">
        <van-icon name="gift-o" size="80" color="#999" />
        <h3>您还不是分销商</h3>
        <p>成为分销商，分享推广获得奖励</p>
        <van-button type="primary" size="large" @click="goToApply">
          立即申请
        </van-button>
      </div>

      <!-- 申请状态 -->
      <div class="application-status" v-if="applicationStatus">
        <van-cell-group>
          <van-cell center>
            <template #title>
              <span class="status-label">申请状态</span>
            </template>
            <template #value>
              <van-tag :type="getStatusType(applicationStatus.status)">
                {{ getStatusText(applicationStatus.status) }}
              </van-tag>
            </template>
          </van-cell>
          <van-cell title="申请品牌" :value="applicationStatus.brandName" />
          <van-cell title="申请时间" :value="applicationStatus.createdAt" />
          <van-cell title="申请理由" :value="applicationStatus.reason" />
          <van-cell
            v-if="applicationStatus.reviewNotes"
            title="审核备注"
            :value="applicationStatus.reviewNotes"
          />
        </van-cell-group>

        <div class="action-buttons" v-if="applicationStatus.status === 'rejected'">
          <van-button type="primary" block @click="goToApply">
            重新申请
          </van-button>
        </div>
      </div>
    </div>

    <!-- 已成为分销商 -->
    <div class="distributor-content" v-else>
       <!-- 统计卡片 -->
       <div class="stats-cards">
         <div class="stat-card">
           <div class="stat-value">¥{{ statistics.balance.toFixed(2) }}</div>
           <div class="stat-label">可提现</div>
         </div>
         <div class="stat-card">
           <div class="stat-value">¥{{ statistics.totalEarnings.toFixed(2) }}</div>
           <div class="stat-label">累计收益</div>
         </div>
         <div class="stat-card">
           <div class="stat-value">{{ statistics.totalOrders }}</div>
           <div class="stat-label">订单数</div>
         </div>
         <div class="stat-card">
           <div class="stat-value">{{ statistics.subordinatesCount }}</div>
           <div class="stat-label">下级数</div>
         </div>
       </div>

      <!-- 功能菜单 -->
      <van-cell-group class="menu-group" inset>
        <van-cell
          title="推广工具"
          icon="link-o"
          is-link
          @click="goToPromotion"
        />
        <van-cell
          title="奖励明细"
          icon="gold-coin-o"
          is-link
          @click="goToRewards"
        />
        <van-cell
          title="下级列表"
          icon="friends-o"
          is-link
          @click="goToSubordinates"
        />
        <van-cell
          title="生成海报"
          icon="photo-o"
          is-link
          @click="goToPosterGenerator"
        />
        <van-cell
          title="订单核销"
          icon="checked"
          is-link
          @click="goToOrderVerify"
        />
         <van-cell
           title="推广数据"
           icon="bar-chart-o"
           is-link
           @click="goToStatistics"
         />
         <van-cell
           title="提现管理"
           icon="balance-o"
           is-link
           @click="goToWithdrawals"
         />
       </van-cell-group>

      <!-- 收益趋势 -->
      <div class="earnings-trend" v-if="statistics.monthEarnings > 0">
        <div class="trend-header">
          <span class="trend-title">本月收益</span>
          <span class="trend-value">¥{{ statistics.monthEarnings.toFixed(2) }}</span>
        </div>
        <div class="trend-info">
          今日收益: ¥{{ statistics.todayEarnings.toFixed(2) }}
        </div>
      </div>
    </div>

    <!-- 底部导航 -->
    <van-tabbar v-model="activeTab" v-if="hasDistributor">
      <van-tabbar-item icon="home-o" @click="$router.push('/')">首页</van-tabbar-item>
      <van-tabbar-item icon="gift-o" @click="refresh">分销</van-tabbar-item>
      <van-tabbar-item icon="orders-o" @click="$router.push('/orders')">订单</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { authApi } from "@/services/brandApi.js";
import axios from "@/utils/axios";
import {
  getStatusType,
  getStatusText,
  getDefaultStatistics,
  clearDistributorAuth
} from "./distributorCenter.logic.js";

export default {
	name: "DistributorCenter",
	setup() {
		const router = useRouter();
		const hasDistributor = ref(false);
		const brands = ref([]);
		const selectedBrandId = ref(null);
		const applicationStatus = ref(null);
		const activeTab = ref(1);
		const statistics = ref(getDefaultStatistics());

		// 加载分销商状态
		const loadDistributorStatus = async () => {
			try {
			const data = await axios.get("/distributor/dashboard");
				if (data.code === 200) {
					hasDistributor.value = data.data.hasDistributor;
					if (data.data.hasDistributor) {
						brands.value = data.data.brands || [];
						if (brands.value.length > 0) {
							selectedBrandId.value = brands.value[0].brandId;
							loadStatistics(selectedBrandId.value);
						}
					}
				}
			} catch (error) {
				console.error("获取分销商状态失败:", error);
			}
		};

		// 加载统计数据
		const loadStatistics = async (brandId) => {
			try {
				const data = await axios.get(
					`/distributor/statistics/${brandId}`,
				);
				if (data.code === 200) {
					statistics.value = data.data || {};
				}
			} catch (error) {
				console.error("获取统计数据失败:", error);
			}
		};

		// 加载申请状态
		const loadApplicationStatus = async () => {
			try {
			const data = await axios.get("/distributor/applications");
				if (data.code === 200 && data.data.applications.length > 0) {
					applicationStatus.value = data.data.applications[0];
				}
			} catch (error) {
				console.error("获取申请状态失败:", error);
			}
		};

		// 切换品牌
		const switchBrand = () => {
			loadStatistics(selectedBrandId.value);
		};

		const handleLogout = async () => {
			try {
				await authApi.logout();
			} catch (error) {
				console.warn("登出失败:", error);
			}

			clearDistributorAuth();
			router.push("/distributor/login");
		};

		// 跳转到申请页面
		const goToApply = () => {
			router.push("/distributor/apply");
		};

		// 跳转到推广工具
		const goToPromotion = () => {
			router.push(`/distributor/promotion?brandId=${selectedBrandId.value}`);
		};

		// 跳转到奖励明细
		const goToRewards = () => {
			router.push(`/distributor/rewards?brandId=${selectedBrandId.value}`);
		};

		// 跳转到下级列表
		const goToSubordinates = () => {
			router.push(`/distributor/subordinates?brandId=${selectedBrandId.value}`);
		};

		// 跳转到统计数据
		const goToStatistics = () => {
			router.push(`/distributor/statistics?brandId=${selectedBrandId.value}`);
		};

		// 跳转到提现管理
		// biome-ignore lint/correctness/noUnusedVariables: used in template
		const goToWithdrawals = () => {
			router.push(`/distributor/withdrawals?brandId=${selectedBrandId.value}`);
		};

		const goToPosterGenerator = () => {
			const campaignId = localStorage.getItem("currentCampaignId") || "1";
			router.push(`/poster-generator/${campaignId}`);
		};

		const goToOrderVerify = () => {
			router.push("/verify");
		};

		// 刷新
		const refresh = () => {
			loadDistributorStatus();
			if (!hasDistributor.value) {
				loadApplicationStatus();
			}
		};

		// 获取状态类型 (re-export from logic)
		const getStatusTypeLocal = getStatusType;

		// 获取状态文本 (re-export from logic)
		const getStatusTextLocal = getStatusText;

		onMounted(() => {
			loadDistributorStatus();
			loadApplicationStatus();
		});

		return {
			hasDistributor,
			brands,
			selectedBrandId,
			statistics,
			applicationStatus,
			activeTab,
			switchBrand,
			goToApply,
			goToPromotion,
			goToRewards,
			goToSubordinates,
			goToPosterGenerator,
			goToOrderVerify,
			goToStatistics,
			handleLogout,
			refresh,
			getStatusType: getStatusTypeLocal,
			getStatusText: getStatusTextLocal,
		};
	},
};
</script>

<style scoped>
.distributor-center {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px;
}

.header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.header h1 {
  margin: 0 0 10px 0;
  font-size: 24px;
}

.logout-btn {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.4);
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 12px;
  cursor: pointer;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.brand-selector select {
  padding: 8px 12px;
  border-radius: 4px;
  border: none;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  font-size: 14px;
}

.not-distributor {
  padding: 20px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 8px;
}

.empty-state h3 {
  margin: 20px 0 10px;
  font-size: 18px;
  color: #333;
}

.empty-state p {
  color: #999;
  margin-bottom: 30px;
}

.application-status {
  margin-top: 20px;
}

.status-label {
  font-weight: 500;
}

.action-buttons {
  margin-top: 20px;
}

.distributor-content {
  padding: 20px 0;
}

.stats-cards {
  display: flex;
  justify-content: space-around;
  background: white;
  padding: 20px 0;
  margin: 0 16px 16px;
  border-radius: 8px;
}

.stat-card {
  text-align: center;
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #667eea;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #999;
}

.menu-group {
  margin-bottom: 16px;
}

.earnings-trend {
  background: white;
  margin: 0 16px;
  padding: 16px;
  border-radius: 8px;
}

.trend-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.trend-title {
  font-size: 14px;
  color: #666;
}

.trend-value {
  font-size: 20px;
  font-weight: bold;
  color: #f56c6c;
}

.trend-info {
  font-size: 12px;
  color: #999;
}
</style>
