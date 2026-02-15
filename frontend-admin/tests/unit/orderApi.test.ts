import { beforeEach, describe, expect, it, vi } from 'vitest';
import { orderApi } from '../../services/orderApi';

describe('orderApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getVerificationRecords sends auth header and returns payload', async () => {
    localStorage.setItem('dmh_token', 'token-xyz');
    const payload = {
      total: 1,
      records: [
        {
          id: 11,
          orderId: 1001,
          verificationStatus: 'verified',
          verificationCode: 'abc',
          verificationMethod: 'scan',
          remark: '',
          createdAt: '2026-02-13T00:00:00Z',
        },
      ],
    };
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => payload,
    } as Response);

    const result = await orderApi.getVerificationRecords();

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(url).toBe('/api/v1/orders/verification-records');
    expect(headers.Authorization).toBe('Bearer token-xyz');
    expect(result).toEqual(payload);
  });

  it('getVerificationRecords works without token', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ total: 0, records: [] }),
    } as Response);

    await orderApi.getVerificationRecords();

    const [, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(headers.Authorization).toBeUndefined();
  });

  it('getVerificationRecords throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 401,
      text: async () => 'unauthorized',
    } as Response);

    await expect(orderApi.getVerificationRecords()).rejects.toThrow(
      'Failed to fetch verification records: 401 unauthorized',
    );
  });
});
