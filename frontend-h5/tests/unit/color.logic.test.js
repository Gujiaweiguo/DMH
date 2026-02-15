import { describe, it, expect } from 'vitest'
import {
  hexToRgb,
  rgbToHex,
  hexToRgba,
  lighten,
  darken,
  isValidHex,
  getContrastColor,
  getRandomColor,
  mixColors,
  PRESET_COLORS
} from '../../src/utils/color.logic.js'

describe('color.logic', () => {
  describe('hexToRgb', () => {
    it('should convert hex to rgb', () => {
      expect(hexToRgb('#ff0000')).toEqual({ r: 255, g: 0, b: 0 })
    })

    it('should handle hex without #', () => {
      expect(hexToRgb('00ff00')).toEqual({ r: 0, g: 255, b: 0 })
    })

    it('should return null for invalid hex', () => {
      expect(hexToRgb('invalid')).toBeNull()
    })
  })

  describe('rgbToHex', () => {
    it('should convert rgb to hex', () => {
      expect(rgbToHex(255, 0, 0)).toBe('#ff0000')
    })
  })

  describe('hexToRgba', () => {
    it('should convert hex to rgba', () => {
      expect(hexToRgba('#ff0000', 0.5)).toBe('rgba(255, 0, 0, 0.5)')
    })
  })

  describe('lighten', () => {
    it('should lighten color', () => {
      const result = lighten('#000000', 50)
      expect(result).not.toBe('#000000')
    })
  })

  describe('darken', () => {
    it('should darken color', () => {
      const result = darken('#ffffff', 50)
      expect(result).not.toBe('#ffffff')
    })
  })

  describe('isValidHex', () => {
    it('should return true for valid hex', () => {
      expect(isValidHex('#ff0000')).toBe(true)
      expect(isValidHex('ff0000')).toBe(true)
    })

    it('should return false for invalid hex', () => {
      expect(isValidHex('invalid')).toBe(false)
    })
  })

  describe('getContrastColor', () => {
    it('should return black for light colors', () => {
      expect(getContrastColor('#ffffff')).toBe('#000000')
    })

    it('should return white for dark colors', () => {
      expect(getContrastColor('#000000')).toBe('#ffffff')
    })
  })

  describe('getRandomColor', () => {
    it('should return valid hex color', () => {
      const color = getRandomColor()
      expect(isValidHex(color)).toBe(true)
    })
  })

  describe('mixColors', () => {
    it('should mix two colors', () => {
      const result = mixColors('#000000', '#ffffff', 0.5)
      expect(result).toBeDefined()
    })
  })

  describe('PRESET_COLORS', () => {
    it('should have primary color', () => {
      expect(PRESET_COLORS.primary).toBe('#4f46e5')
    })
  })
})
