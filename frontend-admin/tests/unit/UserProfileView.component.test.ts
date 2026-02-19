import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';

const mockProfileApi = vi.hoisted(() => ({
  getUserInfo: vi.fn(),
  updateProfile: vi.fn(),
  sendPhoneCode: vi.fn(),
  sendEmailCode: vi.fn(),
  bindPhone: vi.fn(),
  bindEmail: vi.fn(),
  changePassword: vi.fn(),
}));

vi.mock('../../services/profileApi', () => ({
  profileApi: mockProfileApi,
}));

import { UserProfileView } from '../../views/UserProfileView';

describe('UserProfileView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockProfileApi.getUserInfo.mockResolvedValue({
      id: 1,
      username: 'demo_user',
      realName: '张三',
      phone: '13800138000',
      email: 'demo@example.com',
      avatar: 'https://example.com/avatar.png',
      status: 'active',
      roles: ['participant'],
      createdAt: '2025-01-01 10:00:00',
    });
  });

  it('should mount without errors', () => {
    const wrapper = mount(UserProfileView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(UserProfileView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.html()).toBeDefined();
  });

  it('should have component instance', () => {
    const wrapper = mount(UserProfileView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.vm).toBeDefined();
  });

  it('loads user info on mounted', () => {
    mount(UserProfileView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });

    expect(mockProfileApi.getUserInfo).toHaveBeenCalledTimes(1);
  });
});
