import { describe, expect, it, vi, beforeEach } from 'vitest';
import { WithdrawalApprovalView } from '../../views/WithdrawalApprovalView';

vi.mock('../../services/distributorApi', () => ({
  distributorApi: {
    getWithdrawals: vi.fn().mockResolvedValue({
      code: 200,
      data: { withdrawals: [], total: 0 }
    }),
    approveWithdrawal: vi.fn().mockResolvedValue({ code: 200 }),
    rejectWithdrawal: vi.fn().mockResolvedValue({ code: 200 }),
  },
}));

describe('WithdrawalApprovalView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should return view model with required properties', () => {
    const viewModel = WithdrawalApprovalView();
    
    expect(viewModel.withdrawals).toBeDefined();
    expect(viewModel.loading).toBeDefined();
    expect(viewModel.refreshing).toBeDefined();
    expect(viewModel.page).toBeDefined();
    expect(viewModel.pageSize).toBeDefined();
    expect(viewModel.total).toBeDefined();
    expect(viewModel.selectedBrandId).toBeDefined();
    expect(viewModel.statusFilter).toBeDefined();
    expect(viewModel.showDetailModal).toBeDefined();
    expect(viewModel.currentWithdrawal).toBeDefined();
    expect(viewModel.approvalForm).toBeDefined();
    expect(viewModel.processing).toBeDefined();
    expect(viewModel.finished).toBeDefined();
  });

  it('should have initial values set correctly', () => {
    const viewModel = WithdrawalApprovalView();
    
    expect(viewModel.page.value).toBe(1);
    expect(viewModel.pageSize.value).toBe(20);
    expect(viewModel.total.value).toBe(0);
    expect(viewModel.statusFilter.value).toBe('pending');
    expect(viewModel.showDetailModal.value).toBe(false);
    expect(viewModel.processing.value).toBe(false);
    expect(viewModel.finished.value).toBe(false);
  });

  it('should have approvalForm initialized correctly', () => {
    const viewModel = WithdrawalApprovalView();
    
    expect(viewModel.approvalForm.value.action).toBe('approve');
    expect(viewModel.approvalForm.value.notes).toBe('');
  });

  it('should have loadWithdrawals method', () => {
    const viewModel = WithdrawalApprovalView();
    expect(typeof viewModel.loadWithdrawals).toBe('function');
  });

  it('should have onRefresh method', () => {
    const viewModel = WithdrawalApprovalView();
    expect(typeof viewModel.onRefresh).toBe('function');
  });

  it('should have onLoad method', () => {
    const viewModel = WithdrawalApprovalView();
    expect(typeof viewModel.onLoad).toBe('function');
  });

  it('should have openDetailModal method', () => {
    const viewModel = WithdrawalApprovalView();
    expect(typeof viewModel.openDetailModal).toBe('function');
  });

  it('should have submitApproval method', () => {
    const viewModel = WithdrawalApprovalView();
    expect(typeof viewModel.submitApproval).toBe('function');
  });

  it('should have getStatusType method', () => {
    const viewModel = WithdrawalApprovalView();
    expect(viewModel.getStatusType('pending')).toBe('warning');
    expect(viewModel.getStatusType('approved')).toBe('primary');
    expect(viewModel.getStatusType('rejected')).toBe('danger');
    expect(viewModel.getStatusType('processing')).toBe('info');
    expect(viewModel.getStatusType('completed')).toBe('success');
    expect(viewModel.getStatusType('failed')).toBe('danger');
    expect(viewModel.getStatusType('unknown')).toBe('default');
  });

  it('should have getStatusText method', () => {
    const viewModel = WithdrawalApprovalView();
    expect(viewModel.getStatusText('pending')).toBe('待审核');
    expect(viewModel.getStatusText('approved')).toBe('已批准');
    expect(viewModel.getStatusText('rejected')).toBe('已拒绝');
    expect(viewModel.getStatusText('processing')).toBe('处理中');
    expect(viewModel.getStatusText('completed')).toBe('已完成');
    expect(viewModel.getStatusText('failed')).toBe('失败');
    expect(viewModel.getStatusText('unknown')).toBe('unknown');
  });

  it('should have getPayTypeText method', () => {
    const viewModel = WithdrawalApprovalView();
    expect(viewModel.getPayTypeText('wechat')).toBe('微信');
    expect(viewModel.getPayTypeText('alipay')).toBe('支付宝');
    expect(viewModel.getPayTypeText('bank')).toBe('银行卡');
    expect(viewModel.getPayTypeText('unknown')).toBe('unknown');
  });

  it('should have formatDate method', () => {
    const viewModel = WithdrawalApprovalView();
    expect(viewModel.formatDate('')).toBe('');
    expect(viewModel.formatDate('2024-01-15T10:30:00Z')).toBeTruthy();
  });
});
