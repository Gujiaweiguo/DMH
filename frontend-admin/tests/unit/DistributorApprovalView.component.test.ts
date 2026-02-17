import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { DistributorApprovalView } from '../../views/DistributorApprovalView';

vi.mock('../../services/distributorApi', () => ({
  distributorApi: {
    getApplications: vi.fn().mockResolvedValue({
      code: 200,
      data: { applications: [] }
    }),
    approveApplication: vi.fn().mockResolvedValue({ code: 200 }),
  },
}));

describe('DistributorApprovalView', () => {
  let warnSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    vi.clearAllMocks();
    warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
  });

  afterEach(() => {
    warnSpy.mockRestore();
  });

  it('should return view model with required properties', () => {
    const viewModel = DistributorApprovalView({ brandId: 1 });
    
    expect(viewModel.applications).toBeDefined();
    expect(viewModel.loading).toBeDefined();
    expect(viewModel.showApproveModal).toBeDefined();
    expect(viewModel.currentApplication).toBeDefined();
    expect(viewModel.approveForm).toBeDefined();
    expect(viewModel.processing).toBeDefined();
  });

  it('should have initial values set correctly', () => {
    const viewModel = DistributorApprovalView({ brandId: 1 });
    
    expect(viewModel.loading.value).toBe(false);
    expect(viewModel.showApproveModal.value).toBe(false);
    expect(viewModel.processing.value).toBe(false);
    expect(viewModel.currentApplication.value).toBeNull();
  });

  it('should have approveForm initialized correctly', () => {
    const viewModel = DistributorApprovalView({ brandId: 1 });
    
    expect(viewModel.approveForm.value.action).toBe('approve');
    expect(viewModel.approveForm.value.level).toBe(1);
    expect(viewModel.approveForm.value.reason).toBe('');
  });

  it('should have loadApplications method', () => {
    const viewModel = DistributorApprovalView({ brandId: 1 });
    expect(typeof viewModel.loadApplications).toBe('function');
  });

  it('should have openApproveModal method', () => {
    const viewModel = DistributorApprovalView({ brandId: 1 });
    expect(typeof viewModel.openApproveModal).toBe('function');
  });

  it('should have submitApproval method', () => {
    const viewModel = DistributorApprovalView({ brandId: 1 });
    expect(typeof viewModel.submitApproval).toBe('function');
  });

  it('should open approve modal with application data', () => {
    const viewModel = DistributorApprovalView({ brandId: 1 });
    const mockApp = { id: 1, username: 'test', status: 'pending' };
    
    viewModel.openApproveModal(mockApp);
    
    expect(viewModel.currentApplication.value).toEqual(mockApp);
    expect(viewModel.showApproveModal.value).toBe(true);
    expect(viewModel.approveForm.value.action).toBe('approve');
    expect(viewModel.approveForm.value.level).toBe(1);
    expect(viewModel.approveForm.value.reason).toBe('');
  });
});
