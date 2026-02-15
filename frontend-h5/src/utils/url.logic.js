/**
 * URL utilities
 */

export const parseQuery = (queryString) => {
  if (!queryString) return {}
  const str = queryString.startsWith('?') ? queryString.slice(1) : queryString
  if (!str) return {}
  
  return str.split('&').reduce((result, pair) => {
    const [key, value] = pair.split('=')
    if (key) {
      result[decodeURIComponent(key)] = value ? decodeURIComponent(value) : ''
    }
    return result
  }, {})
}

export const stringifyQuery = (params) => {
  if (!params || typeof params !== 'object') return ''
  
  return Object.entries(params)
    .filter(([_, v]) => v !== null && v !== undefined && v !== '')
    .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
    .join('&')
}

export const buildUrl = (baseUrl, params) => {
  if (!baseUrl) return ''
  const query = stringifyQuery(params)
  if (!query) return baseUrl
  const separator = baseUrl.includes('?') ? '&' : '?'
  return baseUrl + separator + query
}

export const getBaseUrl = (url) => {
  if (!url) return ''
  const index = url.indexOf('?')
  return index > 0 ? url.substring(0, index) : url
}

export const getQueryFromUrl = (url) => {
  if (!url) return {}
  const index = url.indexOf('?')
  return index > 0 ? parseQuery(url.substring(index + 1)) : {}
}

export const isValidUrl = (url) => {
  if (!url) return false
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

export const isAbsoluteUrl = (url) => {
  if (!url) return false
  return /^[a-z][a-z0-9+.-]*:/.test(url)
}

export const isRelativeUrl = (url) => !isAbsoluteUrl(url)

export const joinPaths = (...paths) => {
  return paths
    .map((path, i) => {
      if (i === 0) return path.replace(/\/+$/, '')
      if (i === paths.length - 1) return path.replace(/^\/+/, '')
      return path.replace(/^\/+|\/+$/g, '')
    })
    .filter(Boolean)
    .join('/')
}

export const getFileExtension = (url) => {
  if (!url) return ''
  const pathname = url.split('?')[0].split('#')[0]
  const lastSlash = pathname.lastIndexOf('/')
  const filename = lastSlash >= 0 ? pathname.substring(lastSlash + 1) : pathname
  const dotIndex = filename.lastIndexOf('.')
  return dotIndex > 0 ? filename.substring(dotIndex + 1).toLowerCase() : ''
}

export const getFileName = (url) => {
  if (!url) return ''
  const pathname = url.split('?')[0].split('#')[0]
  const index = pathname.lastIndexOf('/')
  return index >= 0 ? pathname.substring(index + 1) : pathname
}

export const isImage = (url) => {
  const ext = getFileExtension(url)
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp'].includes(ext)
}

export const isVideo = (url) => {
  const ext = getFileExtension(url)
  return ['mp4', 'webm', 'ogg', 'avi', 'mov'].includes(ext)
}
