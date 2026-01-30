const API_BASE_URL = '/api/v1';

export interface PosterRecord {
  id: number;
  recordType: string;
  campaignId: number;
  distributorId: number;
  templateName: string;
  posterUrl: string;
  thumbnailUrl: string;
  fileSize: string;
  generationTime: number;
  downloadCount: number;
  shareCount: number;
  generatedBy?: number;
  status: string;
  createdAt: string;
}

export interface PosterRecordsListResponse {
  total: number;
  records: PosterRecord[];
}

export interface PosterTemplate {
  id: number;
  name: string;
  previewImage: string;
  config: Record<string, any>;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface PosterTemplatesListResponse {
  total: number;
  templates: PosterTemplate[];
}

class PosterApiService {
  private getAuthHeaders() {
    const token = localStorage.getItem('dmh_token');
    return {
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
    };
  }

  // 获取海报生成记录列表
  async getPosterRecords(): Promise<PosterRecordsListResponse> {
    const response = await fetch(`${API_BASE_URL}/poster/records`, {
      headers: this.getAuthHeaders(),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to fetch poster records: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  // 获取海报模板列表
  async getPosterTemplates(): Promise<PosterTemplatesListResponse> {
    const response = await fetch(`${API_BASE_URL}/poster/templates`, {
      headers: this.getAuthHeaders(),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to fetch poster templates: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  // 生成活动海报
  async generateCampaignPoster(campaignId: number, templateId?: number): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/campaigns/${campaignId}/poster`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify({ templateId }),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to generate poster: ${response.status} ${errorText}`);
    }
    return response.json();
  }

  // 生成分销员海报
  async generateDistributorPoster(distributorId: number, templateId?: number): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/distributors/${distributorId}/poster`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify({ templateId }),
    });
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to generate distributor poster: ${response.status} ${errorText}`);
    }
    return response.json();
  }
}

export const posterApi = new PosterApiService();
