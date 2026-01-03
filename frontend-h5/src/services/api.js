// API 基础配置
const API_BASE_URL = '/api/v1'

// 获取认证头
const getAuthHeaders = () => {
  const token = localStorage.getItem('dmh_token')
  return token ? { 'Authorization': `Bearer ${token}` } : {}
}

// 通用请求方法
const request = async (url, options = {}) => {
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
      ...options.headers
    },
    ...options
  }

  try {
    const fullUrl = `${API_BASE_URL}${url}`
    console.log('API请求:', fullUrl, config)
    
    const response = await fetch(fullUrl, config)
    
    console.log('API响应:', response.status, response.statusText)
    
    if (!response.ok) {
      // 尝试获取错误详情
      let errorMessage = `HTTP error! status: ${response.status}`
      try {
        const errorData = await response.text()
        if (errorData) {
          errorMessage += ` - ${errorData}`
        }
      } catch (e) {
        // 忽略解析错误
      }
      throw new Error(errorMessage)
    }
    
    const data = await response.json()
    console.log('API数据:', data)
    return data
  } catch (error) {
    console.error('API request failed:', error)
    throw error
  }
}

export const api = {
  // GET 请求
  get: (url, params = {}) => {
    const queryString = new URLSearchParams(params).toString()
    const fullUrl = queryString ? `${url}?${queryString}` : url
    return request(fullUrl, { method: 'GET' })
  },

  // POST 请求
  post: (url, data = {}) => {
    return request(url, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // PUT 请求
  put: (url, data = {}) => {
    return request(url, {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  },

  // DELETE 请求
  delete: (url) => {
    return request(url, { method: 'DELETE' })
  }
}