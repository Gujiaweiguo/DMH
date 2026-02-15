/**
 * PosterGenerator business logic
 */

export const isWeixinEnvironment = (userAgent) => {
  return /micromessenger/i.test(userAgent || '')
}

export const canZoomIn = (scale, maxScale = 3) => scale < maxScale

export const canZoomOut = (scale, minScale = 0.5) => scale > minScale

export const zoomIn = (scale, step = 0.25, maxScale = 3) => {
  return canZoomIn(scale, maxScale) ? scale + step : scale
}

export const zoomOut = (scale, step = 0.25, minScale = 0.5) => {
  return canZoomOut(scale, minScale) ? scale - step : scale
}

export const resetZoom = () => 1

export const rotateLeft = (rotation, step = 90) => rotation - step

export const rotateRight = (rotation, step = 90) => rotation + step

export const resetRotation = () => 0

export const buildPosterFileName = (timestamp) => `poster_${timestamp}.png`

export const getDefaultState = () => ({
  loading: {
    templates: false,
    preview: false,
    generate: false
  },
  templates: [],
  selectedTemplateId: null,
  previewUrl: '',
  posterUrl: '',
  qrcodeUrl: '',
  scale: 1,
  rotation: 0
})

export const canGeneratePoster = (selectedTemplateId, previewUrl) => {
  return !!(selectedTemplateId && previewUrl)
}

export const canDownloadPoster = (posterUrl) => !!posterUrl

export const canSharePoster = (posterUrl) => !!posterUrl

export const fallbackCopyText = (text, document) => {
  try {
    const textArea = document.createElement('textarea')
    textArea.value = text
    textArea.style.position = 'fixed'
    textArea.style.left = '-9999px'
    textArea.style.top = '0'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()
    const successful = document.execCommand('copy')
    document.body.removeChild(textArea)
    return { success: successful }
  } catch (err) {
    return { success: false, error: err }
  }
}
