import { describe, it, expect } from 'vitest'
import {
  truncate,
  capitalize,
  capitalizeWords,
  trim,
  isEmpty,
  isNotEmpty,
  maskPhone,
  maskEmail,
  maskIdCard,
  formatNumber,
  formatCurrency,
  formatPercent,
  formatFileSize,
  slugify,
  escapeHtml,
  unescapeHtml
} from '../../src/utils/string.logic.js'

describe('string.logic', () => {
  describe('truncate', () => {
    it('should truncate long string', () => {
      expect(truncate('1234567890', 5)).toBe('12...')
    })

    it('should not truncate short string', () => {
      expect(truncate('123', 5)).toBe('123')
    })

    it('should return empty for null', () => {
      expect(truncate(null, 5)).toBe('')
    })
  })

  describe('capitalize', () => {
    it('should capitalize first letter', () => {
      expect(capitalize('hello')).toBe('Hello')
    })

    it('should lowercase rest', () => {
      expect(capitalize('HELLO')).toBe('Hello')
    })

    it('should return empty for null', () => {
      expect(capitalize(null)).toBe('')
    })
  })

  describe('capitalizeWords', () => {
    it('should capitalize each word', () => {
      expect(capitalizeWords('hello world')).toBe('Hello World')
    })
  })

  describe('trim', () => {
    it('should trim whitespace', () => {
      expect(trim('  hello  ')).toBe('hello')
    })

    it('should return empty for null', () => {
      expect(trim(null)).toBe('')
    })
  })

  describe('isEmpty', () => {
    it('should return true for empty string', () => {
      expect(isEmpty('')).toBe(true)
      expect(isEmpty('   ')).toBe(true)
    })

    it('should return false for non-empty string', () => {
      expect(isEmpty('hello')).toBe(false)
    })
  })

  describe('isNotEmpty', () => {
    it('should return false for empty string', () => {
      expect(isNotEmpty('')).toBe(false)
    })

    it('should return true for non-empty string', () => {
      expect(isNotEmpty('hello')).toBe(true)
    })
  })

  describe('maskPhone', () => {
    it('should mask phone number', () => {
      expect(maskPhone('13800138000')).toBe('138****8000')
    })

    it('should return original for invalid length', () => {
      expect(maskPhone('123')).toBe('123')
    })
  })

  describe('maskEmail', () => {
    it('should mask email', () => {
      expect(maskEmail('test@example.com')).toBe('t***t@example.com')
    })

    it('should return original for invalid email', () => {
      expect(maskEmail('invalid')).toBe('invalid')
    })
  })

  describe('maskIdCard', () => {
    it('should mask id card', () => {
      expect(maskIdCard('123456789012345678')).toBe('1234**********5678')
    })
  })

  describe('formatNumber', () => {
    it('should format number with decimals', () => {
      expect(formatNumber(123.456, 2)).toBe('123.46')
    })

    it('should return 0 for invalid', () => {
      expect(formatNumber('abc')).toBe('0')
    })
  })

  describe('formatCurrency', () => {
    it('should format currency', () => {
      expect(formatCurrency(123.4)).toBe('¥123.40')
    })

    it('should use custom symbol', () => {
      expect(formatCurrency(100, '$')).toBe('$100.00')
    })
  })

  describe('formatPercent', () => {
    it('should format percent', () => {
      expect(formatPercent(0.123)).toBe('12.3%')
    })

    it('should use custom decimals', () => {
      expect(formatPercent(0.123, 2)).toBe('12.30%')
    })
  })

  describe('formatFileSize', () => {
    it('should format bytes', () => {
      expect(formatFileSize(500)).toBe('500.00 B')
    })

    it('should format KB', () => {
      expect(formatFileSize(1024)).toBe('1.00 KB')
    })

    it('should format MB', () => {
      expect(formatFileSize(1024 * 1024)).toBe('1.00 MB')
    })

    it('should return 0 B for null', () => {
      expect(formatFileSize(null)).toBe('0 B')
    })
  })

  describe('slugify', () => {
    it('should create slug', () => {
      expect(slugify('Hello World')).toBe('hello-world')
    })

    it('should remove special chars', () => {
      expect(slugify('Hello!@#$World')).toBe('helloworld')
    })
  })

  describe('escapeHtml', () => {
    it('should escape html chars', () => {
      expect(escapeHtml('<div>"test"</div>')).toBe('&lt;div&gt;&quot;test&quot;&lt;/div&gt;')
    })

    it('should return empty for null', () => {
      expect(escapeHtml(null)).toBe('')
    })
  })

  describe('unescapeHtml', () => {
    it('should unescape html chars', () => {
      expect(unescapeHtml('&lt;div&gt;')).toBe('<div>')
    })
  })
})
