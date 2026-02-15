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
import {
  DEFAULT_H5_URL,
  DEFAULT_API_URL,
  STATUS_CHECKING,
  STATUS_OK,
  STATUS_FAIL,
  STATUS_RUNNING,
  LOGIN_TEST_NOT_TESTED,
  LOGIN_TEST_TESTING,
  LOGIN_TEST_SUCCESS,
  LOGIN_TEST_FAIL,
  getConnectionStatusText,
  getLoginTestResultText,
  buildLoginTestPayload,
  isConnectionOk,
  getDefaultState,
  LOGIN_ROUTE
} from './apiTest.logic.js'

const router = useRouter()

const h5Status = ref(STATUS_CHECKING)
const h5Url = ref(DEFAULT_H5_URL)
const apiStatus = ref(STATUS_CHECKING)
const apiUrl = ref(DEFAULT_API_URL)
const loginTestStatus = ref(LOGIN_TEST_NOT_TESTED)
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
      body: JSON.stringify(buildLoginTestPayload())
    })
    
    if (isConnectionOk(response)) {
      apiStatus.value = STATUS_OK
      Toast.success('API连接正常')
    } else {
      apiStatus.value = getConnectionStatusText(false, response.status)
      Toast.fail(`连接异常: ${response.status}`)
    }
  } catch (error) {
    apiStatus.value = STATUS_FAIL
    errorMessage.value = error.message
    Toast.fail('连接失败')
  }
}

const testLogin = async () => {
  Toast.loading('测试登录...')
  loginTestStatus.value = LOGIN_TEST_TESTING
  errorMessage.value = ''
  loginTestResult.value = ''
  
  try {
    const result = await authApi.login('brand_admin', 'password')
    loginTestStatus.value = LOGIN_TEST_SUCCESS
    loginTestResult.value = getLoginTestResultText(result)
    Toast.success('登录测试成功')
  } catch (error) {
    loginTestStatus.value = LOGIN_TEST_FAIL
    errorMessage.value = error.message
    Toast.fail('登录测试失败')
  }
}

const goToLogin = () => {
  router.push(LOGIN_ROUTE)
}

onMounted(() => {
  h5Status.value = STATUS_RUNNING
  
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