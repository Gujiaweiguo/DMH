/**
 * Object utilities
 */

export const pick = (obj, keys) => {
  if (!obj || typeof obj !== 'object') return {}
  return keys.reduce((result, key) => {
    if (key in obj) result[key] = obj[key]
    return result
  }, {})
}

export const omit = (obj, keys) => {
  if (!obj || typeof obj !== 'object') return {}
  const result = { ...obj }
  for (const key of keys) {
    delete result[key]
  }
  return result
}

export const merge = (target, source) => {
  if (!target || typeof target !== 'object') return source
  if (!source || typeof source !== 'object') return target
  return { ...target, ...source }
}

export const deepMerge = (target, source) => {
  if (!target || typeof target !== 'object') return source
  if (!source || typeof source !== 'object') return target
  
  const result = { ...target }
  for (const key of Object.keys(source)) {
    if (typeof source[key] === 'object' && source[key] !== null && !Array.isArray(source[key])) {
      result[key] = deepMerge(result[key], source[key])
    } else {
      result[key] = source[key]
    }
  }
  return result
}

export const isEmpty = (obj) => {
  if (!obj) return true
  if (typeof obj !== 'object') return true
  return Object.keys(obj).length === 0
}

export const isNotEmpty = (obj) => !isEmpty(obj)

export const hasKey = (obj, key) => {
  return obj && typeof obj === 'object' && key in obj
}

export const getKeys = (obj) => {
  if (!obj || typeof obj !== 'object') return []
  return Object.keys(obj)
}

export const getValues = (obj) => {
  if (!obj || typeof obj !== 'object') return []
  return Object.values(obj)
}

export const getEntries = (obj) => {
  if (!obj || typeof obj !== 'object') return []
  return Object.entries(obj)
}

export const fromEntries = (entries) => {
  if (!Array.isArray(entries)) return {}
  return Object.fromEntries(entries)
}

export const mapKeys = (obj, fn) => {
  if (!obj || typeof obj !== 'object') return {}
  const result = {}
  for (const [key, value] of Object.entries(obj)) {
    result[fn(value, key)] = value
  }
  return result
}

export const mapValues = (obj, fn) => {
  if (!obj || typeof obj !== 'object') return {}
  const result = {}
  for (const [key, value] of Object.entries(obj)) {
    result[key] = fn(value, key)
  }
  return result
}

export const invert = (obj) => {
  if (!obj || typeof obj !== 'object') return {}
  const result = {}
  for (const [key, value] of Object.entries(obj)) {
    result[value] = key
  }
  return result
}

export const clone = (obj) => {
  if (!obj || typeof obj !== 'object') return obj
  return JSON.parse(JSON.stringify(obj))
}

export const isEqual = (a, b) => {
  return JSON.stringify(a) === JSON.stringify(b)
}

export const getPath = (obj, path, defaultValue = undefined) => {
  if (!obj || typeof obj !== 'object') return defaultValue
  const keys = path.split('.')
  let result = obj
  for (const key of keys) {
    if (result && typeof result === 'object' && key in result) {
      result = result[key]
    } else {
      return defaultValue
    }
  }
  return result
}

export const setPath = (obj, path, value) => {
  if (!obj || typeof obj !== 'object') return obj
  const keys = path.split('.')
  const result = { ...obj }
  let current = result
  for (let i = 0; i < keys.length - 1; i++) {
    current[keys[i]] = { ...current[keys[i]] }
    current = current[keys[i]]
  }
  current[keys[keys.length - 1]] = value
  return result
}
