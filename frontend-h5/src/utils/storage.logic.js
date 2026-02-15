/**
 * Storage utilities with safe access
 */

export const safeGetItem = (key) => {
  try {
    return localStorage.getItem(key)
  } catch {
    return null
  }
}

export const safeSetItem = (key, value) => {
  try {
    localStorage.setItem(key, value)
    return true
  } catch {
    return false
  }
}

export const safeRemoveItem = (key) => {
  try {
    localStorage.removeItem(key)
    return true
  } catch {
    return false
  }
}

export const safeGetJson = (key, defaultValue = null) => {
  try {
    const value = localStorage.getItem(key)
    return value ? JSON.parse(value) : defaultValue
  } catch {
    return defaultValue
  }
}

export const safeSetJson = (key, value) => {
  try {
    localStorage.setItem(key, JSON.stringify(value))
    return true
  } catch {
    return false
  }
}

export const clearAll = () => {
  try {
    localStorage.clear()
    return true
  } catch {
    return false
  }
}

export const getStorageKeys = () => {
  try {
    const keys = []
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key) keys.push(key)
    }
    return keys
  } catch {
    return []
  }
}

export const getStorageSize = () => {
  try {
    let size = 0
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key) {
        const value = localStorage.getItem(key)
        size += (value?.length || 0) + key.length
      }
    }
    return size
  } catch {
    return 0
  }
}

export const hasKey = (key) => {
  try {
    return localStorage.getItem(key) !== null
  } catch {
    return false
  }
}

export const STORAGE_KEYS = {
  TOKEN: 'dmh_token',
  USER_ROLE: 'dmh_user_role',
  USER_INFO: 'dmh_user_info',
  BRAND_ID: 'dmh_current_brand_id',
  SOURCE: 'dmh_source',
  MY_PHONE: 'dmh_my_phone'
}

export const getToken = () => safeGetItem(STORAGE_KEYS.TOKEN)

export const setToken = (token) => safeSetItem(STORAGE_KEYS.TOKEN, token)

export const removeToken = () => safeRemoveItem(STORAGE_KEYS.TOKEN)

export const getUserRole = () => safeGetItem(STORAGE_KEYS.USER_ROLE)

export const setUserRole = (role) => safeSetItem(STORAGE_KEYS.USER_ROLE, role)

export const getUserInfo = () => safeGetJson(STORAGE_KEYS.USER_INFO)

export const setUserInfo = (info) => safeSetJson(STORAGE_KEYS.USER_INFO, info)

export const getBrandId = () => safeGetItem(STORAGE_KEYS.BRAND_ID)

export const setBrandId = (id) => safeSetItem(STORAGE_KEYS.BRAND_ID, String(id))

export const isLoggedIn = () => !!getToken()
