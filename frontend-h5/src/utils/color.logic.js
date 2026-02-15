/**
 * Color utilities
 */

export const hexToRgb = (hex) => {
  if (!hex) return null
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  return result ? {
    r: parseInt(result[1], 16),
    g: parseInt(result[2], 16),
    b: parseInt(result[3], 16)
  } : null
}

export const rgbToHex = (r, g, b) => {
  return '#' + [r, g, b].map(x => {
    const hex = Math.round(clamp(x, 0, 255)).toString(16)
    return hex.length === 1 ? '0' + hex : hex
  }).join('')
}

export const hexToRgba = (hex, alpha) => {
  const rgb = hexToRgb(hex)
  if (!rgb) return ''
  return `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, ${clamp(alpha, 0, 1)})`
}

export const lighten = (hex, percent) => {
  const rgb = hexToRgb(hex)
  if (!rgb) return hex
  const amount = Math.round(2.55 * percent)
  return rgbToHex(
    Math.min(255, rgb.r + amount),
    Math.min(255, rgb.g + amount),
    Math.min(255, rgb.b + amount)
  )
}

export const darken = (hex, percent) => {
  const rgb = hexToRgb(hex)
  if (!rgb) return hex
  const amount = Math.round(2.55 * percent)
  return rgbToHex(
    Math.max(0, rgb.r - amount),
    Math.max(0, rgb.g - amount),
    Math.max(0, rgb.b - amount)
  )
}

export const isValidHex = (hex) => {
  return /^#?([a-f\d]{3}|[a-f\d]{6})$/i.test(hex)
}

export const clamp = (value, min, max) => Math.min(Math.max(value, min), max)

export const getContrastColor = (hex) => {
  const rgb = hexToRgb(hex)
  if (!rgb) return '#000000'
  const brightness = (rgb.r * 299 + rgb.g * 587 + rgb.b * 114) / 1000
  return brightness > 128 ? '#000000' : '#ffffff'
}

export const getRandomColor = () => {
  return '#' + Math.floor(Math.random() * 16777215).toString(16).padStart(6, '0')
}

export const mixColors = (color1, color2, weight = 0.5) => {
  const rgb1 = hexToRgb(color1)
  const rgb2 = hexToRgb(color2)
  if (!rgb1 || !rgb2) return ''
  
  const w = clamp(weight, 0, 1)
  return rgbToHex(
    Math.round(rgb1.r + (rgb2.r - rgb1.r) * w),
    Math.round(rgb1.g + (rgb2.g - rgb1.g) * w),
    Math.round(rgb1.b + (rgb2.b - rgb1.b) * w)
  )
}

export const PRESET_COLORS = {
  primary: '#4f46e5',
  secondary: '#6366f1',
  success: '#10b981',
  warning: '#f59e0b',
  danger: '#ef4444',
  info: '#3b82f6',
  dark: '#1f2937',
  light: '#f3f4f6',
  white: '#ffffff',
  black: '#000000'
}
