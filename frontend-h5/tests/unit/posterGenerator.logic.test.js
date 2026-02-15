import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  isWeixinEnvironment,
  canZoomIn,
  canZoomOut,
  zoomIn,
  zoomOut,
  resetZoom,
  rotateLeft,
  rotateRight,
  resetRotation,
  buildPosterFileName,
  getDefaultState,
  canGeneratePoster,
  canDownloadPoster,
  canSharePoster,
  fallbackCopyText
} from '../../src/views/poster/posterGenerator.logic.js'

describe('posterGenerator.logic', () => {
  describe('isWeixinEnvironment', () => {
    it('should return true for WeChat user agent', () => {
      expect(isWeixinEnvironment('Mozilla/5.0 MicroMessenger/6.0')).toBe(true)
    })

    it('should return false for non-WeChat user agent', () => {
      expect(isWeixinEnvironment('Mozilla/5.0 Chrome/90.0')).toBe(false)
    })

    it('should return false for null', () => {
      expect(isWeixinEnvironment(null)).toBe(false)
    })
  })

  describe('canZoomIn', () => {
    it('should return true when below max scale', () => {
      expect(canZoomIn(1)).toBe(true)
      expect(canZoomIn(2.5)).toBe(true)
    })

    it('should return false when at max scale', () => {
      expect(canZoomIn(3)).toBe(false)
      expect(canZoomIn(4)).toBe(false)
    })
  })

  describe('canZoomOut', () => {
    it('should return true when above min scale', () => {
      expect(canZoomOut(1)).toBe(true)
      expect(canZoomOut(1)).toBe(true)
    })

    it('should return false when at min scale', () => {
      expect(canZoomOut(0.5)).toBe(false)
      expect(canZoomOut(0.3)).toBe(false)
    })
  })

  describe('zoomIn', () => {
    it('should increase scale by step', () => {
      expect(zoomIn(1)).toBe(1.25)
    })

    it('should not exceed max scale', () => {
      expect(zoomIn(3)).toBe(3)
      expect(zoomIn(2.75)).toBe(3)
    })
  })

  describe('zoomOut', () => {
    it('should decrease scale by step', () => {
      expect(zoomOut(1)).toBe(0.75)
    })

    it('should not go below min scale', () => {
      expect(zoomOut(0.5)).toBe(0.5)
      expect(zoomOut(0.75)).toBe(0.5)
    })
  })

  describe('resetZoom', () => {
    it('should return 1', () => {
      expect(resetZoom()).toBe(1)
    })
  })

  describe('rotateLeft', () => {
    it('should decrease rotation by 90', () => {
      expect(rotateLeft(0)).toBe(-90)
      expect(rotateLeft(90)).toBe(0)
    })
  })

  describe('rotateRight', () => {
    it('should increase rotation by 90', () => {
      expect(rotateRight(0)).toBe(90)
      expect(rotateRight(-90)).toBe(0)
    })
  })

  describe('resetRotation', () => {
    it('should return 0', () => {
      expect(resetRotation()).toBe(0)
    })
  })

  describe('buildPosterFileName', () => {
    it('should build filename with timestamp', () => {
      expect(buildPosterFileName(1234567890)).toBe('poster_1234567890.png')
    })
  })

  describe('getDefaultState', () => {
    it('should return default state object', () => {
      const state = getDefaultState()
      expect(state.loading.templates).toBe(false)
      expect(state.loading.preview).toBe(false)
      expect(state.loading.generate).toBe(false)
      expect(state.templates).toEqual([])
      expect(state.selectedTemplateId).toBeNull()
      expect(state.previewUrl).toBe('')
      expect(state.posterUrl).toBe('')
      expect(state.scale).toBe(1)
      expect(state.rotation).toBe(0)
    })
  })

  describe('canGeneratePoster', () => {
    it('should return true when template selected and preview available', () => {
      expect(canGeneratePoster('template-1', 'preview-url')).toBe(true)
    })

    it('should return false when no template selected', () => {
      expect(canGeneratePoster(null, 'preview-url')).toBe(false)
    })

    it('should return false when no preview', () => {
      expect(canGeneratePoster('template-1', '')).toBe(false)
    })
  })

  describe('canDownloadPoster', () => {
    it('should return true when poster URL exists', () => {
      expect(canDownloadPoster('poster-url')).toBe(true)
    })

    it('should return false when no poster URL', () => {
      expect(canDownloadPoster('')).toBe(false)
      expect(canDownloadPoster(null)).toBe(false)
    })
  })

  describe('canSharePoster', () => {
    it('should return true when poster URL exists', () => {
      expect(canSharePoster('poster-url')).toBe(true)
    })

    it('should return false when no poster URL', () => {
      expect(canSharePoster('')).toBe(false)
      expect(canSharePoster(null)).toBe(false)
    })
  })

  describe('fallbackCopyText', () => {
    let mockDocument
    let mockTextArea

    beforeEach(() => {
      mockTextArea = {
        value: '',
        style: {},
        focus: vi.fn(),
        select: vi.fn()
      }
      mockDocument = {
        createElement: vi.fn(() => mockTextArea),
        body: {
          appendChild: vi.fn(),
          removeChild: vi.fn()
        },
        execCommand: vi.fn(() => true)
      }
    })

    it('should copy text successfully', () => {
      const result = fallbackCopyText('test text', mockDocument)
      expect(result.success).toBe(true)
      expect(mockDocument.body.appendChild).toHaveBeenCalled()
      expect(mockDocument.body.removeChild).toHaveBeenCalled()
    })

    it('should handle copy failure', () => {
      mockDocument.execCommand = vi.fn(() => false)
      const result = fallbackCopyText('test text', mockDocument)
      expect(result.success).toBe(false)
    })

    it('should handle exceptions', () => {
      mockDocument.execCommand = vi.fn(() => {
        throw new Error('Copy failed')
      })
      const result = fallbackCopyText('test text', mockDocument)
      expect(result.success).toBe(false)
      expect(result.error).toBeDefined()
    })
  })
})
