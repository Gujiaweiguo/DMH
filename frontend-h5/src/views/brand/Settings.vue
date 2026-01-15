<template>
  <div class="brand-settings">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <h1 class="nav-title">品牌设置</h1>
    </div>

    <!-- 品牌信息 -->
    <div class="settings-section">
      <h2 class="section-title">品牌信息</h2>
      <div class="settings-card">
        <div class="brand-logo-section">
          <div class="logo-preview">
            <img :src="brandInfo.logo" alt="品牌logo" class="logo-img">
          </div>
          <div class="logo-actions">
            <button @click="uploadLogo" class="upload-btn">更换Logo</button>
            <p class="upload-hint">建议尺寸: 200x200px</p>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">品牌名称</label>
          <input
            v-model="brandInfo.name"
            type="text"
            class="form-input"
            placeholder="请输入品牌名称"
          >
        </div>

        <div class="form-group">
          <label class="form-label">品牌描述</label>
          <textarea
            v-model="brandInfo.description"
            class="form-textarea"
            placeholder="请输入品牌描述"
            rows="3"
          ></textarea>
        </div>

        <div class="form-group">
          <label class="form-label">联系电话</label>
          <input
            v-model="brandInfo.phone"
            type="tel"
            class="form-input"
            placeholder="请输入联系电话"
          >
        </div>

        <div class="form-group">
          <label class="form-label">联系邮箱</label>
          <input
            v-model="brandInfo.email"
            type="email"
            class="form-input"
            placeholder="请输入联系邮箱"
          >
        </div>

        <button @click="saveBrandInfo" :disabled="saving" class="save-btn">
          {{ saving ? '保存中...' : '保存品牌信息' }}
        </button>
      </div>
    </div>

    <!-- 奖励设置 -->
    <div class="settings-section">
      <h2 class="section-title">奖励设置</h2>
      <div class="settings-card">
        <div class="form-group">
          <label class="form-label">默认奖励比例</label>
          <div class="input-with-unit">
            <input
              v-model.number="rewardSettings.defaultRate"
              type="number"
              class="form-input"
              placeholder="20"
              min="0"
              max="100"
              step="0.1"
            >
            <span class="input-unit">%</span>
          </div>
          <p class="form-hint">推广员获得的奖励占订单金额的比例</p>
        </div>

        <div class="form-group">
          <label class="form-label">最低提现金额</label>
          <div class="input-with-unit">
            <input
              v-model.number="rewardSettings.minWithdraw"
              type="number"
              class="form-input"
              placeholder="100"
              min="1"
              step="1"
            >
            <span class="input-unit">元</span>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">奖励结算方式</label>
          <div class="radio-group">
            <label class="radio-option">
              <input
                v-model="rewardSettings.settlementType"
                type="radio"
                value="instant"
                class="radio-input"
              >
              <span class="radio-label">即时结算</span>
              <span class="radio-desc">支付成功后立即发放奖励</span>
            </label>
            <label class="radio-option">
              <input
                v-model="rewardSettings.settlementType"
                type="radio"
                value="daily"
                class="radio-input"
              >
              <span class="radio-label">每日结算</span>
              <span class="radio-desc">每天固定时间统一结算</span>
            </label>
          </div>
        </div>

        <button @click="saveRewardSettings" :disabled="saving" class="save-btn">
          {{ saving ? '保存中...' : '保存奖励设置' }}
        </button>
      </div>
    </div>

    <!-- 通知设置 -->
    <div class="settings-section">
      <h2 class="section-title">通知设置</h2>
      <div class="settings-card">
        <div class="notification-item">
          <div class="notification-info">
            <h4 class="notification-title">新订单通知</h4>
            <p class="notification-desc">有新订单时发送通知</p>
          </div>
          <label class="switch">
            <input
              v-model="notificationSettings.newOrder"
              type="checkbox"
              class="switch-input"
            >
            <span class="switch-slider"></span>
          </label>
        </div>

        <div class="notification-item">
          <div class="notification-info">
            <h4 class="notification-title">推广员注册通知</h4>
            <p class="notification-desc">有新推广员注册时发送通知</p>
          </div>
          <label class="switch">
            <input
              v-model="notificationSettings.newPromoter"
              type="checkbox"
              class="switch-input"
            >
            <span class="switch-slider"></span>
          </label>
        </div>

        <div class="notification-item">
          <div class="notification-info">
            <h4 class="notification-title">每日数据报告</h4>
            <p class="notification-desc">每日发送数据统计报告</p>
          </div>
          <label class="switch">
            <input
              v-model="notificationSettings.dailyReport"
              type="checkbox"
              class="switch-input"
            >
            <span class="switch-slider"></span>
          </label>
        </div>

        <div class="form-group">
          <label class="form-label">通知邮箱</label>
          <input
            v-model="notificationSettings.email"
            type="email"
            class="form-input"
            placeholder="请输入接收通知的邮箱"
          >
        </div>

        <button @click="saveNotificationSettings" :disabled="saving" class="save-btn">
          {{ saving ? '保存中...' : '保存通知设置' }}
        </button>
      </div>
    </div>

    <!-- 数据同步设置 -->
    <div class="settings-section">
      <h2 class="section-title">数据同步</h2>
      <div class="settings-card">
        <div class="sync-status">
          <div class="status-info">
            <h4 class="status-title">同步状态</h4>
            <p class="status-desc">与外部系统的数据同步状态</p>
          </div>
          <span :class="['status-badge', syncSettings.status]">
            {{ getSyncStatusText(syncSettings.status) }}
          </span>
        </div>

        <div class="form-group">
          <label class="form-label">同步频率</label>
          <select v-model="syncSettings.frequency" class="form-select">
            <option value="realtime">实时同步</option>
            <option value="hourly">每小时</option>
            <option value="daily">每日</option>
            <option value="manual">手动同步</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">同步数据类型</label>
          <div class="checkbox-group">
            <label class="checkbox-option">
              <input
                v-model="syncSettings.dataTypes"
                type="checkbox"
                value="orders"
                class="checkbox-input"
              >
              <span class="checkbox-label">订单数据</span>
            </label>
            <label class="checkbox-option">
              <input
                v-model="syncSettings.dataTypes"
                type="checkbox"
                value="users"
                class="checkbox-input"
              >
              <span class="checkbox-label">用户数据</span>
            </label>
            <label class="checkbox-option">
              <input
                v-model="syncSettings.dataTypes"
                type="checkbox"
                value="rewards"
                class="checkbox-input"
              >
              <span class="checkbox-label">奖励数据</span>
            </label>
          </div>
        </div>

        <div class="sync-actions">
          <button @click="testSync" class="test-btn">测试连接</button>
          <button @click="manualSync" class="sync-btn">立即同步</button>
        </div>

        <button @click="saveSyncSettings" :disabled="saving" class="save-btn">
          {{ saving ? '保存中...' : '保存同步设置' }}
        </button>
      </div>
    </div>

    <!-- 账户管理 -->
    <div class="settings-section">
      <h2 class="section-title">账户管理</h2>
      <div class="settings-card">
        <div class="account-item">
          <div class="account-info">
            <h4 class="account-title">修改密码</h4>
            <p class="account-desc">定期修改密码以保证账户安全</p>
          </div>
          <button @click="showChangePassword = true" class="action-btn">
            修改密码
          </button>
        </div>

        <div class="account-item">
          <div class="account-info">
            <h4 class="account-title">数据导出</h4>
            <p class="account-desc">导出品牌相关的所有数据</p>
          </div>
          <button @click="exportAllData" class="action-btn">
            导出数据
          </button>
        </div>

        <div class="account-item danger">
          <div class="account-info">
            <h4 class="account-title">注销账户</h4>
            <p class="account-desc">永久删除账户和所有相关数据</p>
          </div>
          <button @click="showDeleteAccount = true" class="action-btn danger">
            注销账户
          </button>
        </div>
      </div>
    </div>

    <!-- 修改密码模态框 -->
    <div v-if="showChangePassword" class="modal-overlay" @click="showChangePassword = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>修改密码</h3>
          <button @click="showChangePassword = false" class="close-btn">✕</button>
        </div>
        
        <div class="password-form">
          <div class="form-group">
            <label>当前密码</label>
            <input
              v-model="passwordForm.oldPassword"
              type="password"
              class="form-input"
              placeholder="请输入当前密码"
            >
          </div>
          <div class="form-group">
            <label>新密码</label>
            <input
              v-model="passwordForm.newPassword"
              type="password"
              class="form-input"
              placeholder="请输入新密码"
            >
          </div>
          <div class="form-group">
            <label>确认新密码</label>
            <input
              v-model="passwordForm.confirmPassword"
              type="password"
              class="form-input"
              placeholder="请再次输入新密码"
            >
          </div>
        </div>

        <div class="modal-actions">
          <button @click="showChangePassword = false" class="cancel-btn">取消</button>
          <button @click="changePassword" class="confirm-btn">确认修改</button>
        </div>
      </div>
    </div>

    <!-- 底部导航 -->
    <div class="bottom-nav">
      <router-link to="/brand/dashboard" class="nav-item">
        <div class="nav-icon">🏠</div>
        <div class="nav-text">工作台</div>
      </router-link>
      <router-link to="/brand/campaigns" class="nav-item">
        <div class="nav-icon">🎯</div>
        <div class="nav-text">活动</div>
      </router-link>
      <router-link to="/brand/orders" class="nav-item">
        <div class="nav-icon">📋</div>
        <div class="nav-text">订单</div>
      </router-link>
      <router-link to="/brand/members" class="nav-item">
        <div class="nav-icon">👤</div>
        <div class="nav-text">会员</div>
      </router-link>
      <router-link to="/brand/promoters" class="nav-item">
        <div class="nav-icon">👥</div>
        <div class="nav-text">推广员</div>
      </router-link>
      <router-link to="/brand/settings" class="nav-item active">
        <div class="nav-icon">⚙️</div>
        <div class="nav-text">设置</div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const saving = ref(false)
const showChangePassword = ref(false)
const showDeleteAccount = ref(false)

const brandInfo = reactive({
  name: '示例品牌',
  description: '专业的数字营销服务品牌',
  logo: 'https://api.dicebear.com/7.x/initials/svg?seed=Brand',
  phone: '400-123-4567',
  email: 'contact@brand.com'
})

const rewardSettings = reactive({
  defaultRate: 20,
  minWithdraw: 100,
  settlementType: 'instant'
})

const notificationSettings = reactive({
  newOrder: true,
  newPromoter: true,
  dailyReport: false,
  email: 'admin@brand.com'
})

const syncSettings = reactive({
  status: 'connected',
  frequency: 'realtime',
  dataTypes: ['orders', 'users']
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const getSyncStatusText = (status) => {
  const statusMap = {
    connected: '已连接',
    disconnected: '未连接',
    error: '连接错误'
  }
  return statusMap[status] || status
}

const uploadLogo = () => {
  // TODO: 实现Logo上传功能
  alert('Logo上传功能开发中...')
}

const saveBrandInfo = async () => {
  saving.value = true
  try {
    // TODO: 调用API保存品牌信息
    await new Promise(resolve => setTimeout(resolve, 1000))
    alert('品牌信息保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

const saveRewardSettings = async () => {
  saving.value = true
  try {
    // TODO: 调用API保存奖励设置
    await new Promise(resolve => setTimeout(resolve, 1000))
    alert('奖励设置保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

const saveNotificationSettings = async () => {
  saving.value = true
  try {
    // TODO: 调用API保存通知设置
    await new Promise(resolve => setTimeout(resolve, 1000))
    alert('通知设置保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

const saveSyncSettings = async () => {
  saving.value = true
  try {
    // TODO: 调用API保存同步设置
    await new Promise(resolve => setTimeout(resolve, 1000))
    alert('同步设置保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

const testSync = async () => {
  try {
    // TODO: 测试同步连接
    await new Promise(resolve => setTimeout(resolve, 2000))
    alert('连接测试成功')
  } catch (error) {
    console.error('测试失败:', error)
    alert('连接测试失败')
  }
}

const manualSync = async () => {
  try {
    // TODO: 手动同步数据
    await new Promise(resolve => setTimeout(resolve, 3000))
    alert('数据同步完成')
  } catch (error) {
    console.error('同步失败:', error)
    alert('数据同步失败')
  }
}

const changePassword = async () => {
  if (!passwordForm.oldPassword || !passwordForm.newPassword) {
    alert('请填写完整的密码信息')
    return
  }
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    alert('两次输入的新密码不一致')
    return
  }
  
  try {
    // TODO: 调用API修改密码
    await new Promise(resolve => setTimeout(resolve, 1000))
    alert('密码修改成功')
    showChangePassword.value = false
    
    // 重置表单
    Object.assign(passwordForm, {
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
  } catch (error) {
    console.error('修改密码失败:', error)
    alert('修改密码失败')
  }
}

const exportAllData = () => {
  // TODO: 实现数据导出功能
  alert('数据导出功能开发中...')
}

const loadSettings = async () => {
  try {
    // TODO: 从API加载设置数据
    console.log('加载设置数据')
  } catch (error) {
    console.error('加载设置失败:', error)
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.brand-settings {
  min-height: 100vh;
  background: #f5f7fa;
  padding-bottom: 80px;
}

.top-nav {
  background: white;
  padding: 16px;
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  z-index: 10;
}

.nav-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
  color: #333;
}

.settings-section {
  padding: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 16px 0;
  color: #333;
}

.settings-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.brand-logo-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.logo-preview {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  overflow: hidden;
  border: 2px solid #e1e5e9;
}

.logo-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.logo-actions {
  flex: 1;
}

.upload-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  margin-bottom: 8px;
}

.upload-hint {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-input,
.form-textarea,
.form-select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e1e5e9;
  border-radius: 12px;
  font-size: 16px;
  transition: border-color 0.3s;
  box-sizing: border-box;
}

.form-input:focus,
.form-textarea:focus,
.form-select:focus {
  outline: none;
  border-color: #667eea;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
}

.form-hint {
  font-size: 12px;
  color: #666;
  margin-top: 4px;
}

.input-with-unit {
  position: relative;
}

.input-unit {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #666;
  font-size: 16px;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.radio-option {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border: 2px solid #e1e5e9;
  border-radius: 12px;
  cursor: pointer;
  transition: border-color 0.3s;
}

.radio-option:has(.radio-input:checked) {
  border-color: #667eea;
  background: #f8f9ff;
}

.radio-input {
  margin-top: 2px;
}

.radio-label {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.radio-desc {
  font-size: 14px;
  color: #666;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.checkbox-option {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}

.checkbox-input {
  width: 18px;
  height: 18px;
}

.checkbox-label {
  font-size: 14px;
  color: #333;
}

.notification-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.notification-item:last-child {
  border-bottom: none;
}

.notification-info {
  flex: 1;
}

.notification-title {
  font-size: 14px;
  font-weight: 500;
  margin: 0 0 4px 0;
  color: #333;
}

.notification-desc {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.switch {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 28px;
}

.switch-input {
  opacity: 0;
  width: 0;
  height: 0;
}

.switch-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: 0.4s;
  border-radius: 28px;
}

.switch-slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  transition: 0.4s;
  border-radius: 50%;
}

.switch-input:checked + .switch-slider {
  background-color: #667eea;
}

.switch-input:checked + .switch-slider:before {
  transform: translateX(20px);
}

.sync-status {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 20px;
}

.status-info {
  flex: 1;
}

.status-title {
  font-size: 14px;
  font-weight: 500;
  margin: 0 0 4px 0;
  color: #333;
}

.status-desc {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.connected {
  background: #e8f5e8;
  color: #4caf50;
}

.status-badge.disconnected {
  background: #fff3e0;
  color: #ff9800;
}

.status-badge.error {
  background: #fce4ec;
  color: #e91e63;
}

.sync-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.test-btn,
.sync-btn {
  flex: 1;
  padding: 12px;
  border: 2px solid #e1e5e9;
  background: white;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.test-btn:hover,
.sync-btn:hover {
  border-color: #667eea;
  background: #f8f9ff;
}

.account-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.account-item:last-child {
  border-bottom: none;
}

.account-item.danger {
  border-color: #fce4ec;
}

.account-info {
  flex: 1;
}

.account-title {
  font-size: 14px;
  font-weight: 500;
  margin: 0 0 4px 0;
  color: #333;
}

.account-desc {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.action-btn {
  background: white;
  border: 2px solid #e1e5e9;
  color: #333;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn:hover {
  border-color: #667eea;
  background: #f8f9ff;
}

.action-btn.danger {
  border-color: #f44336;
  color: #f44336;
}

.action-btn.danger:hover {
  background: #fce4ec;
}

.save-btn {
  width: 100%;
  background: #667eea;
  color: white;
  border: none;
  padding: 16px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.3s;
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 20px;
  width: 100%;
  max-width: 400px;
}

.modal-header {
  padding: 20px 20px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #999;
}

.password-form {
  padding: 20px;
}

.modal-actions {
  padding: 20px;
  display: flex;
  gap: 12px;
}

.cancel-btn,
.confirm-btn {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.confirm-btn {
  background: #667eea;
  color: white;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  display: flex;
  border-top: 1px solid #eee;
  padding: 8px 0;
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: #999;
  padding: 8px;
}

.nav-item.active {
  color: #667eea;
}

.nav-icon {
  font-size: 20px;
  margin-bottom: 4px;
}

.nav-text {
  font-size: 12px;
}
</style>
