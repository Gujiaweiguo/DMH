const API_BASE_URL = '/api/v1';

export interface PasswordPolicy {
  id: number;
  minLength: number;
  requireUppercase: boolean;
  requireLowercase: boolean;
  requireNumbers: boolean;
  requireSpecialChars: boolean;
  maxAge: number;
  historyCount: number;
  maxLoginAttempts: number;
  lockoutDuration: number;
  sessionTimeout: number;
  maxConcurrentSessions: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface UserSession {
  id: string;
  userId: number;
  clientIp: string;
  userAgent: string;
  loginAt: string;
  lastActiveAt: string;
  expiresAt: string;
  status: string;
  createdAt: string;
}

export interface SecurityEvent {
  id: number;
  eventType: string;
  severity: string;
  userId?: number;
  username: string;
  clientIp: string;
  userAgent: string;
  description: string;
  details: string;
  handled: boolean;
  handledBy?: number;
  handledAt?: string;
  createdAt: string;
}

interface UserSessionListResp {
  total: number;
  sessions: UserSession[];
}

interface SecurityEventListResp {
  total: number;
  events: SecurityEvent[];
}

class SecurityApiService {
  private getToken() {
    return localStorage.getItem('dmh_token');
  }

  private async request<T>(url: string, options: RequestInit = {}): Promise<T> {
    const token = this.getToken();
    if (!token) {
      throw new Error('请先登录');
    }

    const response = await fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
        ...options.headers,
      },
    });

    if (response.status === 401) {
      localStorage.removeItem('dmh_token');
      throw new Error('登录已过期，请重新登录');
    }

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: '请求失败' }));
      throw new Error(error.message || '请求失败');
    }

    return response.json() as Promise<T>;
  }

  async getPasswordPolicy(): Promise<PasswordPolicy> {
    return this.request<PasswordPolicy>('/security/password-policy', {
      method: 'GET',
    });
  }

  async updatePasswordPolicy(payload: PasswordPolicy): Promise<PasswordPolicy> {
    return this.request<PasswordPolicy>('/security/password-policy', {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
  }

  async getUserSessions(page = 1, pageSize = 20): Promise<UserSessionListResp> {
    return this.request<UserSessionListResp>(`/security/sessions?page=${page}&pageSize=${pageSize}`, {
      method: 'GET',
    });
  }

  async revokeSession(sessionId: string): Promise<void> {
    await this.request<{ message: string }>(`/security/sessions/${sessionId}`, {
      method: 'DELETE',
    });
  }

  async forceLogoutUser(userId: number, reason?: string): Promise<void> {
    await this.request<{ message: string }>(`/security/force-logout/${userId}`, {
      method: 'POST',
      body: JSON.stringify({ reason: reason || '管理员操作' }),
    });
  }

  async getSecurityEvents(page = 1, pageSize = 20): Promise<SecurityEventListResp> {
    return this.request<SecurityEventListResp>(`/security/events?page=${page}&pageSize=${pageSize}`, {
      method: 'GET',
    });
  }

  async handleSecurityEvent(eventId: number, note?: string): Promise<void> {
    await this.request<{ message: string }>(`/security/events/${eventId}/handle`, {
      method: 'POST',
      body: JSON.stringify({ note: note || '' }),
    });
  }
}

export const securityApi = new SecurityApiService();
