/**
 * Array utilities
 */

export const groupBy = (array, key) => {
  if (!Array.isArray(array)) return {}
  return array.reduce((result, item) => {
    const groupKey = typeof key === 'function' ? key(item) : item[key]
    if (!result[groupKey]) result[groupKey] = []
    result[groupKey].push(item)
    return result
  }, {})
}

export const sortBy = (array, key, order = 'asc') => {
  if (!Array.isArray(array)) return []
  return [...array].sort((a, b) => {
    const valA = typeof key === 'function' ? key(a) : a[key]
    const valB = typeof key === 'function' ? key(b) : b[key]
    if (valA < valB) return order === 'asc' ? -1 : 1
    if (valA > valB) return order === 'asc' ? 1 : -1
    return 0
  })
}

export const uniqueBy = (array, key) => {
  if (!Array.isArray(array)) return []
  const seen = new Set()
  return array.filter(item => {
    const k = typeof key === 'function' ? key(item) : item[key]
    if (seen.has(k)) return false
    seen.add(k)
    return true
  })
}

export const unique = (array) => {
  if (!Array.isArray(array)) return []
  return [...new Set(array)]
}

export const chunk = (array, size) => {
  if (!Array.isArray(array)) return []
  if (size <= 0) return [array]
  const result = []
  for (let i = 0; i < array.length; i += size) {
    result.push(array.slice(i, i + size))
  }
  return result
}

export const flatten = (array) => {
  if (!Array.isArray(array)) return []
  return array.flat(Infinity)
}

export const first = (array) => {
  if (!Array.isArray(array) || array.length === 0) return undefined
  return array[0]
}

export const last = (array) => {
  if (!Array.isArray(array) || array.length === 0) return undefined
  return array[array.length - 1]
}

export const isEmptyArray = (array) => {
  return !Array.isArray(array) || array.length === 0
}

export const isNotEmptyArray = (array) => !isEmptyArray(array)

export const sum = (array, key) => {
  if (!Array.isArray(array)) return 0
  return array.reduce((total, item) => {
    const value = key ? item[key] : item
    return total + (Number(value) || 0)
  }, 0)
}

export const average = (array, key) => {
  if (!Array.isArray(array) || array.length === 0) return 0
  return sum(array, key) / array.length
}

export const min = (array, key) => {
  if (!Array.isArray(array) || array.length === 0) return undefined
  const values = key ? array.map(item => item[key]) : array
  return Math.min(...values.map(v => Number(v) || 0))
}

export const max = (array, key) => {
  if (!Array.isArray(array) || array.length === 0) return undefined
  const values = key ? array.map(item => item[key]) : array
  return Math.max(...values.map(v => Number(v) || 0))
}

export const range = (start, end, step = 1) => {
  const result = []
  for (let i = start; i < end; i += step) {
    result.push(i)
  }
  return result
}

export const shuffle = (array) => {
  if (!Array.isArray(array)) return []
  const result = [...array]
  for (let i = result.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[result[i], result[j]] = [result[j], result[i]]
  }
  return result
}
