import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, h } from 'vue'
import { PermissionProvider, PermissionGuard, PermissionButton, usePermission } from '../../../components/PermissionGuard'

describe('PermissionGuard', () => {
  const mockUser = {
    username: 'admin',
    roles: ['platform_admin'] as const,
    token: 'test-token'
  }

  describe('PermissionProvider', () => {
    it('provides permission context to children', () => {
      const Child = defineComponent({
        setup() {
          return () => h('div', 'child')
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(Child)
      ))

      expect(wrapper.exists()).toBe(true)
    })

    it('hasPermission returns true for platform_admin', () => {
      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasPermission('any:permission')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('true')
    })

    it('hasPermission returns false for null user', () => {
      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasPermission('any:permission')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: null },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('false')
    })

    it('hasRole returns correct value', () => {
      const userWithRoles = {
        username: 'test',
        roles: ['participant'] as const,
        token: 'test-token'
      }

      const Child = defineComponent({
        setup() {
          const { hasRole } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasRole('participant')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userWithRoles },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('true')
    })

    it('hasRole returns false for non-matching role', () => {
      const userWithDifferentRole = {
        username: 'test',
        roles: ['brand_admin'] as const,
        token: 'test-token'
      }

      const Child = defineComponent({
        setup() {
          const { hasRole } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(hasRole('participant')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userWithDifferentRole },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('false')
    })

    it('canAccessBrand returns true for platform_admin', () => {
      const Child = defineComponent({
        setup() {
          const { canAccessBrand } = usePermission()
          return () => h('div', { 'data-testid': 'result' }, String(canAccessBrand(123)))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(Child)
      ))

      expect(wrapper.find('[data-testid="result"]').text()).toBe('true')
    })
  })

  describe('PermissionGuard', () => {
    it('renders children when has permission', () => {
      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(
          PermissionGuard,
          { permission: 'campaign:read' },
          () => h('div', { 'data-testid': 'content' }, 'Protected Content')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(true)
    })

    it('does not render children when no permission', () => {
      const userNoPerm = {
        username: 'user',
        roles: ['participant'] as const,
        token: 'test-token'
      }

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userNoPerm },
        () => h(
          PermissionGuard,
          { permission: 'admin:delete' },
          () => h('div', { 'data-testid': 'content' }, 'Protected Content')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(false)
    })

    it('renders fallback when provided and no permission', () => {
      const userNoPerm = {
        username: 'user',
        roles: ['participant'] as const,
        token: 'test-token'
      }

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userNoPerm },
        () => h(
          PermissionGuard,
          { permission: 'admin:delete', fallback: 'No Access' },
          () => h('div', { 'data-testid': 'content' }, 'Protected Content')
        )
      ))

      expect(wrapper.text()).toContain('No Access')
      expect(wrapper.find('[data-testid="content"]').exists()).toBe(false)
    })

    it('checks role correctly', () => {
      const userParticipant = {
        username: 'user',
        roles: ['participant'] as const,
        token: 'test-token'
      }

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionGuard,
          { role: 'participant' },
          () => h('div', { 'data-testid': 'content' }, 'Has Role')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(true)
    })

    it('checks roles array correctly', () => {
      const userParticipant = {
        username: 'user',
        roles: ['participant'] as const,
        token: 'test-token'
      }

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userParticipant },
        () => h(
          PermissionGuard,
          { roles: ['participant', 'brand_admin'] as any },
          () => h('div', { 'data-testid': 'content' }, 'Has Role')
        )
      ))

      expect(wrapper.find('[data-testid="content"]').exists()).toBe(true)
    })
  })

  describe('PermissionButton', () => {
    it('renders button when has permission', () => {
      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(
          PermissionButton,
          { permission: 'campaign:read' },
          () => h('span', 'Click Me')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(true)
      expect(wrapper.text()).toContain('Click Me')
    })

    it('does not render button when no permission', () => {
      const userNoPerm = {
        username: 'user',
        roles: ['participant'] as const,
        token: 'test-token'
      }

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userNoPerm },
        () => h(
          PermissionButton,
          { permission: 'admin:delete' },
          () => h('span', 'Click Me')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(false)
    })

    it('handles click when enabled', async () => {
      const handleClick = vi.fn()

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(
          PermissionButton,
          { permission: 'campaign:read', onClick: handleClick },
          () => h('span', 'Click Me')
        )
      ))

      await wrapper.find('button').trigger('click')
      expect(handleClick).toHaveBeenCalled()
    })

    it('button does not render when no permission', () => {
      const userNoPerm = {
        username: 'user',
        roles: ['participant'] as const,
        token: 'test-token'
      }

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: userNoPerm },
        () => h(
          PermissionButton,
          { permission: 'admin:delete' },
          () => h('span', 'Click Me')
        )
      ))

      expect(wrapper.find('button').exists()).toBe(false)
    })
  })

  describe('usePermission', () => {
    it('returns context when used within provider', () => {
      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', String(hasPermission('test')))
        }
      })

      const wrapper = mount(() => h(
        PermissionProvider,
        { user: mockUser },
        () => h(Child)
      ))

      expect(wrapper.text()).toBe('true')
    })

    it('returns fallback when used outside provider', () => {
      const Child = defineComponent({
        setup() {
          const { hasPermission } = usePermission()
          return () => h('div', String(hasPermission('test')))
        }
      })

      const wrapper = mount(() => h(Child))

      expect(wrapper.text()).toBe('false')
    })
  })
})
