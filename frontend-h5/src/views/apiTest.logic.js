/**
 * ApiTest business logic
 */

export const DEFAULT_H5_URL = 'http://localhost:3100'
export const DEFAULT_API_URL = '/api/v1'

export const STATUS_CHECKING = '检测中...'
export const STATUS_OK = '✅ 连接正常'
export const STATUS_FAIL = '❌ 连接失败'
export const STATUS_RUNNING = '✅ 运行正常'

export const LOGIN_TEST_NOT_TESTED = '未测试'
export const LOGIN_TEST_TESTING = '测试中...'
export const LOGIN_TEST_SUCCESS = '✅ 登录成功'
export const LOGIN_TEST_FAIL = '❌ 登录失败'

export const getConnectionStatusText = (ok, status) => {
  if (ok) return STATUS_OK
  if (status) return `❌ 连接异常 (${status})`
  return STATUS_FAIL
}

export const getLoginTestResultText = (result) => {
  if (!result) return ''
  return `用户: ${result.username}, 角色: ${result.roles?.join(', ')}`
}

export const buildLoginTestPayload = () => ({
  username: 'test',
  password: 'test'
})

export const isConnectionOk = (response) => {
  return response.ok || response.status === 401
}

export const getDefaultState = () => ({
  h5Status: STATUS_CHECKING,
  h5Url: DEFAULT_H5_URL,
  apiStatus: STATUS_CHECKING,
  apiUrl: DEFAULT_API_URL,
  loginTestStatus: LOGIN_TEST_NOT_TESTED,
  loginTestResult: '',
  errorMessage: ''
})

export const LOGIN_ROUTE = '/brand/login'
