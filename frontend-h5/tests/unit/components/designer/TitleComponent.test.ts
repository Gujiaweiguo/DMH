import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TitleComponent from '../../../../src/components/designer/TitleComponent.vue'

describe('TitleComponent', () => {
  it('renders title with default text', () => {
    const wrapper = mount(TitleComponent, {
      props: {
        data: {}
      }
    })

    expect(wrapper.find('.title').text()).toBe('活动标题')
  })

  it('renders title with custom text', () => {
    const wrapper = mount(TitleComponent, {
      props: {
        data: {
          title: 'My Campaign'
        }
      }
    })

    expect(wrapper.find('.title').text()).toBe('My Campaign')
  })

  it('renders subtitle when provided', () => {
    const wrapper = mount(TitleComponent, {
      props: {
        data: {
          title: 'Test',
          subtitle: 'This is subtitle'
        }
      }
    })

    expect(wrapper.find('.subtitle').exists()).toBe(true)
    expect(wrapper.find('.subtitle').text()).toBe('This is subtitle')
  })

  it('does not render subtitle when not provided', () => {
    const wrapper = mount(TitleComponent, {
      props: {
        data: {
          title: 'Test'
        }
      }
    })

    expect(wrapper.find('.subtitle').exists()).toBe(false)
  })

  it('applies custom color', () => {
    const wrapper = mount(TitleComponent, {
      props: {
        data: {
          title: 'Test',
          color: '#ff0000'
        }
      }
    })

    const title = wrapper.find('.title')
    expect(title.attributes('style')).toBeDefined()
  })
})
