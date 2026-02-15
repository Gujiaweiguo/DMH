import { describe, it, expect } from 'vitest'
import {
  COMPONENT_TYPES,
  COMPONENT_LIBRARY,
  getComponentType,
  canMoveUp,
  canMoveDown,
  moveUp,
  moveDown,
  removeComponent,
  addComponent,
  generateComponentId,
  createComponent,
  hasComponents,
  getComponentCount
} from '../../src/views/brand/campaignPageDesigner.logic.js'

describe('campaignPageDesigner.logic', () => {
  describe('COMPONENT_TYPES', () => {
    it('should have correct type mappings', () => {
      expect(COMPONENT_TYPES.banner).toBe('BannerComponent')
      expect(COMPONENT_TYPES.text).toBe('TitleComponent')
    })
  })

  describe('COMPONENT_LIBRARY', () => {
    it('should have component options', () => {
      expect(COMPONENT_LIBRARY.length).toBeGreaterThan(0)
    })
  })

  describe('getComponentType', () => {
    it('should return mapped component type', () => {
      expect(getComponentType('banner')).toBe('BannerComponent')
    })

    it('should return original type for unknown', () => {
      expect(getComponentType('unknown')).toBe('unknown')
    })
  })

  describe('canMoveUp', () => {
    it('should return false for first item', () => {
      expect(canMoveUp(0)).toBe(false)
    })

    it('should return true for other items', () => {
      expect(canMoveUp(1)).toBe(true)
      expect(canMoveUp(2)).toBe(true)
    })
  })

  describe('canMoveDown', () => {
    it('should return false for last item', () => {
      expect(canMoveDown(2, 3)).toBe(false)
    })

    it('should return true for other items', () => {
      expect(canMoveDown(0, 3)).toBe(true)
      expect(canMoveDown(1, 3)).toBe(true)
    })
  })

  describe('moveUp', () => {
    it('should swap with previous item', () => {
      const components = [{ id: 'a' }, { id: 'b' }, { id: 'c' }]
      const result = moveUp(components, 1)
      expect(result[0].id).toBe('b')
      expect(result[1].id).toBe('a')
    })

    it('should not change for first item', () => {
      const components = [{ id: 'a' }, { id: 'b' }]
      const result = moveUp(components, 0)
      expect(result[0].id).toBe('a')
    })
  })

  describe('moveDown', () => {
    it('should swap with next item', () => {
      const components = [{ id: 'a' }, { id: 'b' }, { id: 'c' }]
      const result = moveDown(components, 1)
      expect(result[1].id).toBe('c')
      expect(result[2].id).toBe('b')
    })

    it('should not change for last item', () => {
      const components = [{ id: 'a' }, { id: 'b' }]
      const result = moveDown(components, 1)
      expect(result[1].id).toBe('b')
    })
  })

  describe('removeComponent', () => {
    it('should remove component by id', () => {
      const components = [{ id: 'a' }, { id: 'b' }, { id: 'c' }]
      const result = removeComponent(components, 'b')
      expect(result).toHaveLength(2)
      expect(result.find(c => c.id === 'b')).toBeUndefined()
    })
  })

  describe('addComponent', () => {
    it('should add component to end', () => {
      const components = [{ id: 'a' }]
      const newComp = { id: 'b' }
      const result = addComponent(components, newComp)
      expect(result).toHaveLength(2)
      expect(result[1].id).toBe('b')
    })
  })

  describe('generateComponentId', () => {
    it('should generate unique id', () => {
      const id1 = generateComponentId()
      const id2 = generateComponentId()
      expect(id1).not.toBe(id2)
      expect(id1).toMatch(/^comp_/)
    })
  })

  describe('createComponent', () => {
    it('should create component with id and type', () => {
      const comp = createComponent('banner', { title: 'test' })
      expect(comp.id).toBeDefined()
      expect(comp.type).toBe('banner')
      expect(comp.data.title).toBe('test')
      expect(comp.style).toEqual({})
    })
  })

  describe('hasComponents', () => {
    it('should return true for non-empty array', () => {
      expect(hasComponents([{ id: 'a' }])).toBe(true)
    })

    it('should return false for empty array', () => {
      expect(hasComponents([])).toBe(false)
    })

    it('should return false for undefined', () => {
      expect(hasComponents(undefined)).toBe(false)
    })
  })

  describe('getComponentCount', () => {
    it('should return correct count', () => {
      expect(getComponentCount([{ id: 'a' }, { id: 'b' }])).toBe(2)
    })

    it('should return 0 for null', () => {
      expect(getComponentCount(null)).toBe(0)
    })
  })
})
