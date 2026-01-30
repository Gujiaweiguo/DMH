const API_BASE_URL = '/api/v1';

export interface VerificationRecord {
  id: number;
  orderId: number;
  verificationStatus: string;
  verifiedAt?: string;
  verifiedBy?: number;
  verificationCode: string;
  verificationMethod: string;
  remark: string;
  createdAt: string;
}

export interface VerificationRecordsListResponse {
  total: number;
  records: VerificationRecord[];
}

class OrderApiService {
  private getAuthHeaders() {
    const token = localStorage.getItem('dmh_token');
    return {
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
    };
  }

  // 获取核销记录列表
  async getVerificationRecords(): Promise<VerificationRecordsListResponse> {
    const response = await fetch(`${API_BASE_URL}/orders/verification-records`, {
      headers: this.getAuthHeaders(),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to fetch verification records: ${response.status} ${errorText}`);
    }
    return response.json();
  }
}

export const orderApi = new OrderApiService();
