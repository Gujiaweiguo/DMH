import { describe, expect, it } from 'vitest'
import {
  buildAIGeneratedText,
  createAIGeneratedMaterial,
  createUploadedMaterial,
  filterMaterialsByCategory,
  formatMaterialDate,
  getDefaultAITextForm,
  getDefaultUploadForm,
  getMaterialCategories,
  getMaterialTypeText,
  getMockMaterials,
  validateAITextInput,
  validateUploadInput,
} from '../../src/views/brand/materials.logic.js'

describe('materials logic', () => {
  it('provides default categories and forms', () => {
    expect(getMaterialCategories()).toHaveLength(4)
    expect(getDefaultUploadForm()).toEqual({
      name: '',
      description: '',
      category: 'image',
      file: null,
    })
    expect(getDefaultAITextForm()).toEqual({
      topic: '',
      style: 'professional',
      length: 'medium',
    })
  })

  it('filters materials by category', () => {
    const list = [
      { id: 1, type: 'image' },
      { id: 2, type: 'text' },
      { id: 3, type: 'image' },
    ]
    expect(filterMaterialsByCategory(list, 'all')).toHaveLength(3)
    expect(filterMaterialsByCategory(list, 'image').map((item) => item.id)).toEqual([1, 3])
  })

  it('maps type text and formats date', () => {
    expect(getMaterialTypeText('image')).toBe('图片')
    expect(getMaterialTypeText('text')).toBe('文案')
    expect(getMaterialTypeText('video')).toBe('视频')
    expect(getMaterialTypeText('other')).toBe('other')
    expect(formatMaterialDate('2026-01-01')).toContain('2026')
  })

  it('validates upload and ai text input', () => {
    expect(validateUploadInput({ file: null, name: '' })).toBe('请选择文件并填写素材名称')
    expect(validateUploadInput({ file: { name: 'a.png' }, name: '素材A' })).toBe('')

    expect(validateAITextInput('')).toBe('请输入文案主题')
    expect(validateAITextInput('春节促销')).toBe('')
  })

  it('creates uploaded and ai generated materials', () => {
    const uploaded = createUploadedMaterial({
      id: 1,
      name: '海报',
      description: 'desc',
      category: 'image',
      url: 'blob:test',
      createdAt: '2026-02-13',
    })
    expect(uploaded).toMatchObject({
      id: 1,
      type: 'image',
      url: 'blob:test',
    })

    const text = buildAIGeneratedText('春节活动')
    expect(text).toContain('春节活动')

    const aiMaterial = createAIGeneratedMaterial({
      id: 2,
      topic: '春节活动',
      createdAt: '2026-02-13',
    })
    expect(aiMaterial).toMatchObject({
      id: 2,
      type: 'text',
      name: 'AI生成-春节活动',
    })
    expect(aiMaterial.content).toContain('春节活动')
  })

  it('provides mock materials', () => {
    const mock = getMockMaterials()
    expect(mock).toHaveLength(3)
    expect(mock[0]).toMatchObject({ type: 'image' })
  })
})
