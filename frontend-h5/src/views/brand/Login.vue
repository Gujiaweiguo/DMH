<template>
  <div class="brand-login">
    <div class="login-container">
      <!-- Logo和标题 -->
      <div class="login-header">
        <div class="logo">🏢</div>
        <h1 class="title">品牌管理登录</h1>
        <p class="subtitle">DMH数字营销中台</p>
      </div>

      <!-- 登录表单 -->
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label class="form-label">用户名</label>
          <input
            v-model="form.username"
            type="text"
            class="form-input"
            placeholder="请输入用户名"
            required
          >
        </div>

        <div class="form-group">
          <label class="form-label">密码</label>
          <input
            v-model="form.password"
            type="password"
            class="form-input"
            placeholder="请输入密码"
            required
          >
        </div>

        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <button
          type="submit"
          class="login-btn"
          :disabled="loading"
        >
          {{ loginButtonText }}
        </button>
      </form>

      <!-- 测试账号提示 -->
      <div class="test-account">
        <p class="test-title">⚠️ 请使用以下测试账号</p>
        <div class="account-info">
          <p class="test-info"><strong>用户名: brand_manager</strong></p>
          <p class="test-info"><strong>密码: 123456</strong></p>
        </div>
        <button 
          type="button" 
          class="quick-fill-btn"
          @click="quickFill"
        >
          🚀 一键填充测试账号
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../../services/brandApi.js'
import {
  getDefaultForm,
  quickFillTestAccount,
  validateLoginForm,
  hasBrandAccess,
  getFirstBrandId,
  saveLoginInfo,
  getLoginButtonText,
  buildLoginError
} from './login.logic.js'

const router = useRouter()

const form = reactive(getDefaultForm())

const loading = ref(false)
const errorMessage = ref('')

const quickFill = () => {
  quickFillTestAccount(form)
}

const loginButtonText = computed(() => getLoginButtonText(loading.value))

const handleLogin = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    console.log('开始登录...', { username: form.username })
    
    // 调用登录API
    const data = await authApi.login(form.username, form.password)
    
    console.log('登录响应:', data)
    
    // 检查响应数据
    if (!data) {
      throw new Error('登录响应为空')
    }
    
    // 检查用户是否有品牌访问权限
    if (!hasBrandAccess(data)) {
      throw new Error('未绑定品牌，请联系管理员为该账号分配品牌权限')
    }

    // 保存当前品牌ID（默认取第一个）
    const firstBrandId = getFirstBrandId(data.brandIds)
    if (!firstBrandId) {
      throw new Error('未绑定品牌，请联系管理员为该账号分配品牌权限')
    }

    // 保存登录信息
    saveLoginInfo(data, firstBrandId)

    console.log('登录成功，跳转到工作台')
    
    // 跳转到品牌工作台
    router.push('/brand/dashboard')
  } catch (error) {
    console.error('登录失败:', error)
    errorMessage.value = buildLoginError(error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.brand-login {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  background: white;
  border-radius: 20px;
  padding: 40px 30px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  font-size: 48px;
  margin-bottom: 16px;
}

.title {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin: 0 0 8px 0;
}

.subtitle {
  color: #666;
  font-size: 14px;
  margin: 0;
}

.login-form {
  margin-bottom: 30px;
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

.form-input {
  width: 100%;
  padding: 16px;
  border: 2px solid #e1e5e9;
  border-radius: 12px;
  font-size: 16px;
  transition: border-color 0.3s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
  margin-bottom: 20px;
  text-align: center;
}

.login-btn {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.3s;
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.test-account {
  background: #fff3cd;
  border: 2px solid #ffc107;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
}

.test-title {
  font-size: 14px;
  font-weight: 600;
  color: #856404;
  margin: 0 0 12px 0;
}

.account-info {
  background: #fff;
  padding: 12px;
  border-radius: 8px;
  margin: 8px 0;
  border: 1px solid #ffc107;
}

.test-info {
  font-size: 14px;
  color: #333;
  margin: 4px 0;
}

.quick-fill-btn {
  width: 100%;
  padding: 12px 16px;
  background: #ffc107;
  color: #000;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  margin-top: 8px;
  transition: background 0.3s;
}

.quick-fill-btn:hover {
  background: #e0a800;
}
</style>
