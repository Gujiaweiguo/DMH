import { describe, expect, it, vi, beforeEach } from 'vitest';
import { distributorApi } from '../../services/distributorApi';

const mockFetch = vi.fn();
global.fetch = mockFetch;

const mockLocalStorage = {
  getItem: vi.fn().mockReturnValue('test-token'),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
Object.defineProperty(global, 'localStorage', { value: mockLocalStorage });

describe('distributorApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockFetch.mockReset();
  });

  describe('getApplications', () => {
    it('should call correct endpoint with brandId', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
      });

      await distributorApi.getApplications(1, { status: 'pending' });

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/brands/1/distributor/applications?status=pending',
        expect.objectContaining({
          headers: expect.objectContaining({
            'Authorization': 'Bearer test-token',
          }),
        })
      );
    });
  });

  describe('getApplication', () => {
    it('should call correct endpoint for single application', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: {} }),
      });

      await distributorApi.getApplication(1, 123);

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/brands/1/distributor/applications/123',
        expect.objectContaining({
          headers: expect.objectContaining({
            'Authorization': 'Bearer test-token',
          }),
        })
      );
    });
  });

  describe('approveApplication', () => {
    it('should send POST request to approve endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ code: 200 }),
      });

      await distributorApi.approveApplication(1, 123, { action: 'approve' });

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/brands/1/distributor/approve/123',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ action: 'approve' }),
        })
      );
    });
  });

  describe('getDistributors', () => {
    it('should call correct endpoint with params', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
      });

      await distributorApi.getDistributors(1, { status: 'active' });

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/brands/1/distributors?status=active',
        expect.any(Object)
      );
    });
  });

  describe('getDistributor', () => {
    it('should call correct endpoint for single distributor', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: {} }),
      });

      await distributorApi.getDistributor(1, 123);

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/brands/1/distributors/123',
        expect.any(Object)
      );
    });
  });

  describe('updateDistributorLevel', () => {
    it('should send PUT request to level endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ code: 200 }),
      });

      await distributorApi.updateDistributorLevel(123, 2);

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/brands/distributors/123/level',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ level: 2 }),
        })
      );
    });
  });

  describe('updateDistributorStatus', () => {
    it('should send PUT request to status endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ code: 200 }),
      });

      await distributorApi.updateDistributorStatus(123, 'suspended', 'Test reason');

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/brands/distributors/123/status',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ status: 'suspended', reason: 'Test reason' }),
        })
      );
    });
  });

  describe('getLevelRewards', () => {
    it('should call correct endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
      });

      await distributorApi.getLevelRewards(1);

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/brands/1/distributor/level-rewards',
        expect.any(Object)
      );
    });
  });

  describe('setLevelRewards', () => {
    it('should send PUT request with rewards', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ code: 200 }),
      });

      const rewards = [{ level: 1, rate: 10 }];
      await distributorApi.setLevelRewards(1, rewards as any);

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/brands/1/distributor/level-rewards',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ rewards }),
        })
      );
    });
  });

  describe('getWithdrawals', () => {
    it('should call correct endpoint with params', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
      });

      await distributorApi.getWithdrawals(1, 'pending', 1, 20);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/platform/withdrawals'),
        expect.any(Object)
      );
    });
  });

  describe('approveWithdrawal', () => {
    it('should send PUT request to approve endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ code: 200 }),
      });

      await distributorApi.approveWithdrawal(123, 'Approved');

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/platform/withdrawals/123/approve',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ notes: 'Approved' }),
        })
      );
    });
  });

  describe('rejectWithdrawal', () => {
    it('should send PUT request to reject endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ code: 200 }),
      });

      await distributorApi.rejectWithdrawal(123, 'Test rejection reason');

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/platform/withdrawals/123/reject',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ reason: 'Test rejection reason' }),
        })
      );
    });
  });

  describe('getGlobalDistributors', () => {
    it('should call correct endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
      });

      await distributorApi.getGlobalDistributors(1, 'active', 1, 20);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/platform/distributors'),
        expect.any(Object)
      );
    });
  });

  describe('getGlobalRewards', () => {
    it('should call correct endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
      });

      await distributorApi.getGlobalRewards(1, 1, 20);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/platform/rewards'),
        expect.any(Object)
      );
    });
  });

  describe('error handling', () => {
    it('should throw error on non-ok response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        text: () => Promise.resolve('Not Found'),
      });

      await expect(distributorApi.getApplications(1)).rejects.toThrow('HTTP error!');
    });
  });
});
