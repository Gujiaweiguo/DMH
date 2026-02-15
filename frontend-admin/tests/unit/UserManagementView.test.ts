import { describe, expect, it } from 'vitest';

const getRoleBadgeColor = (roleCode: string) => {
  const colors: Record<string, string> = {
    platform_admin: 'bg-purple-100 text-purple-700 border-purple-200',
    brand_admin: 'bg-blue-100 text-blue-700 border-blue-200',
    participant: 'bg-green-100 text-green-700 border-green-200',
  };
  return colors[roleCode] || 'bg-gray-100 text-gray-700 border-gray-200';
};

const getRoleName = (roleCode: string) => {
  const names: Record<string, string> = {
    platform_admin: '平台管理员',
    participant: '活动参与者',
  };
  return names[roleCode] || roleCode;
};

const isValidPassword = (password: string) => {
  return password.length >= 6;
};

const filterUsers = (
  users: any[],
  searchQuery: string,
  filterRole: string,
  filterStatus: string
) => {
  return users.filter(user => {
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      const matchSearch = 
        user.username.toLowerCase().includes(query) ||
        user.realName?.toLowerCase().includes(query) ||
        user.phone.includes(query) ||
        user.email?.toLowerCase().includes(query);
      if (!matchSearch) return false;
    }
    
    if (filterRole !== 'all') {
      if (!user.roles.includes(filterRole)) return false;
    }
    
    if (filterStatus !== 'all') {
      if (user.status !== filterStatus) return false;
    }
    
    return true;
  });
};

const mockUsers = [
  {
    id: 1,
    username: 'admin',
    realName: '系统管理员',
    phone: '13800000001',
    email: 'admin@dmh.com',
    status: 'active',
    roles: ['platform_admin'],
  },
  {
    id: 2,
    username: 'user001',
    realName: '张三',
    phone: '13800000003',
    email: 'user001@dmh.com',
    status: 'active',
    roles: ['participant'],
  },
  {
    id: 3,
    username: 'user002',
    realName: '李四',
    phone: '13800000004',
    email: 'user002@dmh.com',
    status: 'disabled',
    roles: ['participant'],
  },
];

describe('UserManagementView', () => {
  describe('getRoleBadgeColor', () => {
    it('should return purple color for platform_admin', () => {
      expect(getRoleBadgeColor('platform_admin')).toBe(
        'bg-purple-100 text-purple-700 border-purple-200'
      );
    });

    it('should return blue color for brand_admin', () => {
      expect(getRoleBadgeColor('brand_admin')).toBe(
        'bg-blue-100 text-blue-700 border-blue-200'
      );
    });

    it('should return green color for participant', () => {
      expect(getRoleBadgeColor('participant')).toBe(
        'bg-green-100 text-green-700 border-green-200'
      );
    });

    it('should return gray color for unknown role', () => {
      expect(getRoleBadgeColor('unknown_role')).toBe(
        'bg-gray-100 text-gray-700 border-gray-200'
      );
    });
  });

  describe('getRoleName', () => {
    it('should return Chinese name for platform_admin', () => {
      expect(getRoleName('platform_admin')).toBe('平台管理员');
    });

    it('should return Chinese name for participant', () => {
      expect(getRoleName('participant')).toBe('活动参与者');
    });

    it('should return original code for unknown role', () => {
      expect(getRoleName('unknown_role')).toBe('unknown_role');
    });
  });

  describe('isValidPassword', () => {
    it('should return true for password with 6 or more characters', () => {
      expect(isValidPassword('123456')).toBe(true);
      expect(isValidPassword('1234567')).toBe(true);
      expect(isValidPassword('verylongpassword')).toBe(true);
    });

    it('should return false for password with less than 6 characters', () => {
      expect(isValidPassword('')).toBe(false);
      expect(isValidPassword('12345')).toBe(false);
    });
  });

  describe('filterUsers', () => {
    it('should return all users when no filters applied', () => {
      const result = filterUsers(mockUsers, '', 'all', 'all');
      expect(result).toHaveLength(3);
    });

    it('should filter by search query (username)', () => {
      const result = filterUsers(mockUsers, 'admin', 'all', 'all');
      expect(result).toHaveLength(1);
      expect(result[0].username).toBe('admin');
    });

    it('should filter by search query (realName)', () => {
      const result = filterUsers(mockUsers, '张三', 'all', 'all');
      expect(result).toHaveLength(1);
      expect(result[0].realName).toBe('张三');
    });

    it('should filter by search query (phone)', () => {
      const result = filterUsers(mockUsers, '13800000001', 'all', 'all');
      expect(result).toHaveLength(1);
      expect(result[0].phone).toBe('13800000001');
    });

    it('should filter by role', () => {
      const result = filterUsers(mockUsers, '', 'platform_admin', 'all');
      expect(result).toHaveLength(1);
      expect(result[0].roles).toContain('platform_admin');
    });

    it('should filter by status', () => {
      const result = filterUsers(mockUsers, '', 'all', 'disabled');
      expect(result).toHaveLength(1);
      expect(result[0].status).toBe('disabled');
    });

    it('should filter by multiple criteria', () => {
      const result = filterUsers(mockUsers, 'user', 'participant', 'active');
      expect(result).toHaveLength(1);
      expect(result[0].username).toBe('user001');
    });

    it('should return empty array when no matches', () => {
      const result = filterUsers(mockUsers, 'nonexistent', 'all', 'all');
      expect(result).toHaveLength(0);
    });

    it('should be case-insensitive for search', () => {
      const result = filterUsers(mockUsers, 'ADMIN', 'all', 'all');
      expect(result).toHaveLength(1);
      expect(result[0].username).toBe('admin');
    });
  });
});
