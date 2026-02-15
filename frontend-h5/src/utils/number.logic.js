/**
 * Number utilities
 */

export const clamp = (value, min, max) => {
  return Math.min(Math.max(value, min), max)
}

export const randomInt = (min, max) => {
  return Math.floor(Math.random() * (max - min + 1)) + min
}

export const randomFloat = (min, max, decimals = 2) => {
  const value = Math.random() * (max - min) + min
  return Number(value.toFixed(decimals))
}

export const isNumber = (value) => {
  return typeof value === 'number' && !isNaN(value)
}

export const isInteger = (value) => {
  return isNumber(value) && Number.isInteger(value)
}

export const isPositive = (value) => {
  return isNumber(value) && value > 0
}

export const isNegative = (value) => {
  return isNumber(value) && value < 0
}

export const isInRange = (value, min, max) => {
  return isNumber(value) && value >= min && value <= max
}

export const toFixed = (value, decimals = 2) => {
  const num = Number(value)
  if (isNaN(num)) return '0'
  return num.toFixed(decimals)
}

export const toPercent = (value, decimals = 1) => {
  const num = Number(value)
  if (isNaN(num)) return '0%'
  return (num * 100).toFixed(decimals) + '%'
}

export const parseNumber = (value, defaultValue = 0) => {
  const num = Number(value)
  return isNaN(num) ? defaultValue : num
}

export const round = (value, decimals = 0) => {
  const num = Number(value)
  if (isNaN(num)) return 0
  const factor = Math.pow(10, decimals)
  return Math.round(num * factor) / factor
}

export const floor = (value) => {
  const num = Number(value)
  return isNaN(num) ? 0 : Math.floor(num)
}

export const ceil = (value) => {
  const num = Number(value)
  return isNaN(num) ? 0 : Math.ceil(num)
}

export const formatWithCommas = (value) => {
  const num = Number(value)
  if (isNaN(num)) return '0'
  return num.toLocaleString()
}

export const formatCurrency = (value, symbol = '¥', decimals = 2) => {
  const num = Number(value)
  if (isNaN(num)) return symbol + '0'
  return symbol + num.toFixed(decimals)
}

export const formatCompact = (value) => {
  const num = Number(value)
  if (isNaN(num)) return '0'
  
  const units = ['', 'K', 'M', 'B', 'T']
  const absNum = Math.abs(num)
  const unitIndex = Math.floor(Math.log10(absNum) / 3)
  
  if (unitIndex === 0) return String(num)
  
  const scaled = num / Math.pow(1000, unitIndex)
  return scaled.toFixed(1) + units[unitIndex]
}

export const degToRad = (degrees) => degrees * (Math.PI / 180)

export const radToDeg = (radians) => radians * (180 / Math.PI)

export const lerp = (start, end, t) => {
  return start + (end - start) * clamp(t, 0, 1)
}
