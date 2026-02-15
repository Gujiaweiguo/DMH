import { describe, it, expect } from 'vitest'
import {
  parseQuery,
  stringifyQuery,
  buildUrl,
  getBaseUrl,
  getQueryFromUrl,
  isValidUrl,
  isAbsoluteUrl,
  isRelativeUrl,
  joinPaths,
  getFileExtension,
  getFileName,
  isImage,
  isVideo
} from '../../src/utils/url.logic.js'

describe('url.logic', () => {
  describe('parseQuery', () => {
    it('should parse query string', () => {
      expect(parseQuery('a=1&b=2')).toEqual({ a: '1', b: '2' })
    })

    it('should handle leading question mark', () => {
      expect(parseQuery('?a=1')).toEqual({ a: '1' })
    })

    it('should return empty for null', () => {
      expect(parseQuery(null)).toEqual({})
    })
  })

  describe('stringifyQuery', () => {
    it('should stringify params', () => {
      expect(stringifyQuery({ a: 1, b: 2 })).toBe('a=1&b=2')
    })

    it('should skip null values', () => {
      expect(stringifyQuery({ a: 1, b: null })).toBe('a=1')
    })

    it('should return empty for null', () => {
      expect(stringifyQuery(null)).toBe('')
    })
  })

  describe('buildUrl', () => {
    it('should build url with params', () => {
      expect(buildUrl('/api', { a: 1 })).toBe('/api?a=1')
    })

    it('should append to existing query', () => {
      expect(buildUrl('/api?x=1', { a: 1 })).toBe('/api?x=1&a=1')
    })

    it('should return base url for empty params', () => {
      expect(buildUrl('/api', {})).toBe('/api')
    })
  })

  describe('getBaseUrl', () => {
    it('should return base url without query', () => {
      expect(getBaseUrl('/api?a=1')).toBe('/api')
    })

    it('should return full url if no query', () => {
      expect(getBaseUrl('/api')).toBe('/api')
    })
  })

  describe('getQueryFromUrl', () => {
    it('should extract query params', () => {
      expect(getQueryFromUrl('/api?a=1&b=2')).toEqual({ a: '1', b: '2' })
    })

    it('should return empty for no query', () => {
      expect(getQueryFromUrl('/api')).toEqual({})
    })
  })

  describe('isValidUrl', () => {
    it('should return true for valid url', () => {
      expect(isValidUrl('https://example.com')).toBe(true)
    })

    it('should return false for invalid url', () => {
      expect(isValidUrl('not a url')).toBe(false)
    })
  })

  describe('isAbsoluteUrl', () => {
    it('should return true for absolute url', () => {
      expect(isAbsoluteUrl('https://example.com')).toBe(true)
      expect(isAbsoluteUrl('http://example.com')).toBe(true)
    })

    it('should return false for relative url', () => {
      expect(isAbsoluteUrl('/api/test')).toBe(false)
    })
  })

  describe('isRelativeUrl', () => {
    it('should return true for relative url', () => {
      expect(isRelativeUrl('/api/test')).toBe(true)
    })
  })

  describe('joinPaths', () => {
    it('should join paths', () => {
      expect(joinPaths('/api', 'users', '1')).toBe('/api/users/1')
    })

    it('should handle trailing slashes', () => {
      expect(joinPaths('/api/', '/users/', '/1')).toBe('/api/users/1')
    })
  })

  describe('getFileExtension', () => {
    it('should return extension', () => {
      expect(getFileExtension('https://example.com/image.png')).toBe('png')
    })

    it('should return empty for no extension', () => {
      expect(getFileExtension('https://example.com/image')).toBe('')
    })
  })

  describe('getFileName', () => {
    it('should return file name', () => {
      expect(getFileName('https://example.com/path/image.png')).toBe('image.png')
    })
  })

  describe('isImage', () => {
    it('should return true for image urls', () => {
      expect(isImage('image.png')).toBe(true)
      expect(isImage('image.jpg')).toBe(true)
    })

    it('should return false for non-image urls', () => {
      expect(isImage('video.mp4')).toBe(false)
    })
  })

  describe('isVideo', () => {
    it('should return true for video urls', () => {
      expect(isVideo('video.mp4')).toBe(true)
      expect(isVideo('video.webm')).toBe(true)
    })
  })
})
