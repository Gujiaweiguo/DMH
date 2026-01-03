<template>
  <div class="api-test">
    <van-nav-bar title="API连接测试" />
    
    <van-cell-group inset title="连接状态">
      <van-cell title="H5服务" :value="h5Status" :label="h5Url" />
      <van-cell title="API代理" :value="apiStatus" :label="apiUrl" />
    </van-cell-group>
    
    <van-cell-group inset title="测试结果">
      <van-cell title="登录API测试" :value="loginTestStatus" />
      <van-cell v-if="loginTestResult" title="返回数据" :label="loginTestResult" />
      <van-cell v-if="errorMessage" title="错误信息" :label="errorMessage" />
    </van-cell-group>
    
    <div class="test-actions">
      <van-button type="primary" block @click="testApiConnection">
        测试API连接
      </van-button>
      
      <van-button type="success" block @click="testLogin" style="margin-top: 12px;">
        测试登录接口
      </van-button>
      
      <van-button plain block @click="goToLogin" style="margin-top: 12px;">
        返回登录页面
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Toast } from 'vant'
import { authApi } from '../services/brandApi.js'

const router = useRouter()

const h5Status = ref('检测中...')
const h5Url = ref('http://localhost:3100')
const apiStatus = ref('检测中...')
const apiUrl = ref('/api/v1')
const loginTestStatus = ref('未测试')
const loginTestResult = ref('')
const errorMessage = ref('')

const testApiConnection = async () => {
  Toast.loading('测试中...')
  
  try {
    // 测试基本连接
    const response = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: 'test',
        password: 'test'
      })
    })
    
    if (response.ok || response.status === 401) {
      apiStatus.value = '✅ 连接正常'
      Toast.success('API连接正常')
    } else {
      apiStatus.value = `❌ 连接异常 (${response.status})`
      Toast.fail(`连接异常: ${response.status}`)
    }
  } catch (error) {
    apiStatus.value = '❌ 连接失败'
    errorMessage.value = error.message
    Toast.fail('连接失败')
  }
}

const testLogin = async () => {
  Toast.loading('测试登录...')
  loginTestStatus.value = '测试中...'
  errorMessage.value = ''
  loginTestResult.value = ''
  
  try {
    const result = await authApi.login('brand_admin', 'password')
    loginTestStatus.value = '✅ 登录成功'
    loginTestResult.value = `用户: ${result.username}, 角色: ${result.roles?.join(', ')}`
    Toast.success('登录测试成功')
  } catch (error) {
    loginTestStatus.value = '❌ 登录失败'
    errorMessage.value = error.message
    Toast.fail('登录测试失败')
  }
}

const goToLogin = () => {
  router.push('/brand/login')
}

onMounted(() => {
  h5Status.value = '✅ 运行正常'
  
  // 自动测试API连接
  setTimeout(() => {
    testApiConnection()
  }, 1000)
})
</script>

<style scoped>
.api-test {
  min-height: 100vh;
  background: #f7f8fa;
}

.test-actions {
  padding: 16px;
  margin-top: 20px;
}
</style>