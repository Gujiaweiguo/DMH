import { beforeEach, describe, expect, it, vi } from 'vitest';
import { feedbackApi } from '../../services/feedbackApi';

describe('feedbackApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getFeedbackList builds query and sends auth header', async () => {
    localStorage.setItem('dmh_token', 'token-feedback');
    const payload = { total: 1, feedbacks: [] };
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => payload,
    } as Response);

    const result = await feedbackApi.getFeedbackList({
      page: 2,
      pageSize: 15,
      category: 'bug',
      status: 'pending',
      priority: 'high',
    });

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(String(url)).toContain(
      '/api/v1/feedback/list?page=2&pageSize=15&category=bug&status=pending&priority=high',
    );
    expect(headers.Authorization).toBe('Bearer token-feedback');
    expect(result).toEqual(payload);
  });

  it('getFeedbackList throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      text: async () => 'internal error',
    } as Response);

    await expect(feedbackApi.getFeedbackList()).rejects.toThrow(
      'Failed to fetch feedback list: 500 internal error',
    );
  });

  it('getFeedbackDetail fetches by id', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ id: 88, title: 'detail' }),
    } as Response);

    const result = await feedbackApi.getFeedbackDetail(88);

    const [url] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/feedback/detail?id=88');
    expect(result).toEqual({ id: 88, title: 'detail' });
  });

  it('updateFeedbackStatus sends PUT and request body', async () => {
    const body = { id: 10, status: 'resolved', response: 'done' };
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ id: 10, status: 'resolved' }),
    } as Response);

    await feedbackApi.updateFeedbackStatus(body);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/feedback/status');
    expect((options as RequestInit).method).toBe('PUT');
    expect((options as RequestInit).body).toBe(JSON.stringify(body));
  });

  it('updateFeedbackStatus throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 403,
      text: async () => 'forbidden',
    } as Response);

    await expect(feedbackApi.updateFeedbackStatus({ id: 1, status: 'closed' })).rejects.toThrow(
      'Failed to update feedback status: 403 forbidden',
    );
  });

  it('getFeedbackStatistics builds query and returns data', async () => {
    const stats = {
      totalFeedbacks: 5,
      byCategory: { bug: 3 },
      byStatus: { pending: 2 },
      byPriority: { high: 1 },
      averageRating: 4,
      resolutionRate: 80,
      avgResolutionTime: 12,
      byRating: { '5': 2 },
    };
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => stats,
    } as Response);

    const result = await feedbackApi.getFeedbackStatistics({
      startDate: '2026-02-01',
      endDate: '2026-02-28',
      category: 'bug',
    });

    const [url] = vi.mocked(fetch).mock.calls[0];
    expect(String(url)).toContain('/api/v1/feedback/statistics?startDate=2026-02-01&endDate=2026-02-28&category=bug');
    expect(result).toEqual(stats);
  });

  it('getFeedbackStatistics throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 400,
      text: async () => 'bad request',
    } as Response);

    await expect(feedbackApi.getFeedbackStatistics()).rejects.toThrow(
      'Failed to fetch feedback statistics: 400 bad request',
    );
  });
});
