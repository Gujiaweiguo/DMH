import { beforeEach, describe, expect, it, vi } from 'vitest';
import { posterApi } from '../../services/posterApi';

describe('posterApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getPosterRecords sends auth header and returns payload', async () => {
    localStorage.setItem('dmh_token', 'token-poster');
    const payload = { total: 1, records: [] };
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => payload,
    } as Response);

    const result = await posterApi.getPosterRecords();

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(url).toBe('/api/v1/poster/records');
    expect(headers.Authorization).toBe('Bearer token-poster');
    expect(result).toEqual(payload);
  });

  it('getPosterTemplates works without token', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ total: 0, templates: [] }),
    } as Response);

    await posterApi.getPosterTemplates();

    const [, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(headers.Authorization).toBeUndefined();
  });

  it('getPosterTemplates throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      text: async () => 'template service down',
    } as Response);

    await expect(posterApi.getPosterTemplates()).rejects.toThrow(
      'Failed to fetch poster templates: 500 template service down',
    );
  });

  it('generateCampaignPoster sends POST with templateId body', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ posterUrl: 'https://cdn/p1.png' }),
    } as Response);

    await posterApi.generateCampaignPoster(101, 7);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/campaigns/101/poster');
    expect((options as RequestInit).method).toBe('POST');
    expect((options as RequestInit).body).toBe(JSON.stringify({ templateId: 7 }));
  });

  it('generateCampaignPoster throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 400,
      text: async () => 'campaign not found',
    } as Response);

    await expect(posterApi.generateCampaignPoster(404)).rejects.toThrow(
      'Failed to generate poster: 400 campaign not found',
    );
  });

  it('generateDistributorPoster sends distributor endpoint and body', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ posterUrl: 'https://cdn/d1.png' }),
    } as Response);

    await posterApi.generateDistributorPoster(88, 3);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/distributors/88/poster');
    expect((options as RequestInit).method).toBe('POST');
    expect((options as RequestInit).body).toBe(JSON.stringify({ templateId: 3 }));
  });

  it('generateDistributorPoster throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 403,
      text: async () => 'forbidden',
    } as Response);

    await expect(posterApi.generateDistributorPoster(9)).rejects.toThrow(
      'Failed to generate distributor poster: 403 forbidden',
    );
  });
});
