// 分销商 API 服务
const API_BASE = '/api/v1';

// 获取 token
const getToken = () => localStorage.getItem('dmh_token');

// 通用请求封装
async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`,
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
  }

  return response.json();
}

export const distributorApi = {
  // 获取分销商申请列表
  getApplications: async (brandId: number, params?: any) => {
    const query = params ? '?' + new URLSearchParams(params).toString() : '';
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/applications${query}`);
  },

  // 获取申请详情
  getApplication: async (brandId: number, applicationId: number) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/applications/${applicationId}`);
  },

  // 审批申请
  approveApplication: async (brandId: number, applicationId: number, data: any) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/approve/${applicationId}`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 获取分销商列表
  getDistributors: async (brandId: number, params?: any) => {
    const query = params ? '?' + new URLSearchParams(params).toString() : '';
    return request<any>(`${API_BASE}/brands/${brandId}/distributors${query}`);
  },

  // 获取分销商详情
  getDistributor: async (brandId: number, distributorId: number) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributors/${distributorId}`);
  },

  // 更新分销商级别
  updateDistributorLevel: async (distributorId: number, level: number) => {
    return request<any>(`${API_BASE}/brands/distributors/${distributorId}/level`, {
      method: 'PUT',
      body: JSON.stringify({ level }),
    });
  },

  // 更新分销商状态
  updateDistributorStatus: async (distributorId: number, status: string, reason?: string) => {
    return request<any>(`${API_BASE}/brands/distributors/${distributorId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status, reason }),
    });
  },

  // 获取级别奖励配置
  getLevelRewards: async (brandId: number) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/level-rewards`);
  },

  // 设置级别奖励配置
  setLevelRewards: async (brandId: number, rewards: any[]) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/level-rewards`, {
      method: 'PUT',
      body: JSON.stringify({ rewards }),
    });
  }
};
