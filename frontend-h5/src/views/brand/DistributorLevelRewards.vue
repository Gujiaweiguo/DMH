<template>
  <div class="level-rewards">
    <div class="top-nav">
      <div class="nav-left">
        <button class="back-btn" @click="back">←</button>
        <h1 class="nav-title">分销奖励设置</h1>
      </div>
      <button class="save-btn" @click="save" :disabled="saving">保存</button>
    </div>

    <div class="content">
      <div class="hint">
        这里配置不同分销级别的奖励比例（%），用于奖励计算。
      </div>

      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else class="cards">
        <div class="card" v-for="row in form.rewards" :key="row.level">
          <div class="card-title">级别 {{ row.level }}</div>
          <div class="field">
            <label>奖励比例（%）</label>
            <input
              v-model.number="row.rewardPercentage"
              type="number"
              min="0"
              max="100"
              step="0.1"
              class="input"
              placeholder="例如 5.5"
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  getDefaultRewards,
  buildSavePayload,
  mergeRewardsFromServer
} from './distributorLevelRewards.logic.js'

const router = useRouter()

const loading = ref(false)
const saving = ref(false)

const form = reactive({
  rewards: getDefaultRewards()
})

const getCurrentBrandId = () => {
  const fromStorage = Number(localStorage.getItem('dmh_current_brand_id'))
  if (Number.isFinite(fromStorage) && fromStorage > 0) return fromStorage
  try {
    const info = JSON.parse(localStorage.getItem('dmh_user_info') || '{}')
    const firstBrandId = Array.isArray(info.brandIds) && info.brandIds.length > 0 ? Number(info.brandIds[0]) : 0
    if (Number.isFinite(firstBrandId) && firstBrandId > 0) {
      localStorage.setItem('dmh_current_brand_id', String(firstBrandId))
      return firstBrandId
    }
  } catch {
    // ignore
  }
  return 0
}

const back = () => router.back()

const load = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('dmh_token')
    const brandId = getCurrentBrandId()
    if (!brandId) {
      alert('未选择品牌，请重新登录')
      router.push('/brand/login')
      return
    }

    const response = await fetch(`/api/v1/brands/${brandId}/distributor/level-rewards`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (!response.ok) return
    const data = await response.json()
    const merged = mergeRewardsFromServer(form.rewards, data.rewards)
    form.rewards.splice(0, form.rewards.length, ...merged)
  } catch (error) {
    console.error('加载奖励配置失败:', error)
  } finally {
    loading.value = false
  }
}

const save = async () => {
  saving.value = true
  try {
    const token = localStorage.getItem('dmh_token')
    const brandId = getCurrentBrandId()
    if (!brandId) {
      alert('未选择品牌，请重新登录')
      router.push('/brand/login')
      return
    }

    const payload = buildSavePayload(form.rewards)

    const response = await fetch(`/api/v1/brands/${brandId}/distributor/level-rewards`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    })

    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      alert(`保存失败：${data.message || '未知错误'}`)
      return
    }
    alert('保存成功')
  } catch (error) {
    console.error('保存奖励配置失败:', error)
    alert('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.level-rewards {
  min-height: 100vh;
  background: #f5f7fa;
}

.top-nav {
  position: sticky;
  top: 0;
  z-index: 10;
  background: white;
  border-bottom: 1px solid #e6ebf2;
  padding: 12px 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.back-btn {
  border: none;
  background: #f1f3f5;
  padding: 8px 10px;
  border-radius: 10px;
  font-weight: 800;
  cursor: pointer;
}

.nav-title {
  margin: 0;
  font-size: 18px;
  font-weight: 800;
  color: #333;
}

.save-btn {
  border: none;
  background: #667eea;
  color: white;
  padding: 10px 14px;
  border-radius: 12px;
  font-weight: 800;
  cursor: pointer;
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.content {
  padding: 14px;
}

.hint {
  background: #fff;
  border: 1px solid #edf1f7;
  border-radius: 14px;
  padding: 12px;
  color: #666;
  font-size: 13px;
  margin-bottom: 12px;
}

.loading {
  text-align: center;
  padding: 30px 0;
  color: #666;
}

.loading-spinner {
  width: 28px;
  height: 28px;
  border: 3px solid #e6ebf2;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 10px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card {
  background: white;
  border: 1px solid #edf1f7;
  border-radius: 16px;
  padding: 14px;
}

.card-title {
  font-weight: 900;
  color: #333;
  margin-bottom: 12px;
}

.field label {
  display: block;
  font-size: 13px;
  color: #666;
  margin-bottom: 8px;
}

.input {
  width: 100%;
  padding: 12px;
  border-radius: 12px;
  border: 1px solid #e6ebf2;
  font-size: 14px;
  box-sizing: border-box;
}
</style>

