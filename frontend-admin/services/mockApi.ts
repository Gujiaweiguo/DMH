
import { ContentLibraryItem, Order, UserProfile, Industry, Brand, MarketingTask, Withdrawal, TaskType } from '../types';

const INITIAL_INDUSTRIES: Industry[] = [
  { id: 'ind-1', name: '教育培训', icon: 'GraduationCap', brandCount: 2 },
  { id: 'ind-2', name: '职业考证', icon: 'FileCheck', brandCount: 1 }
];

const INITIAL_BRANDS: Brand[] = [
  { id: 'br-1', name: '博学教育', industryId: 'ind-1', logo: 'https://api.dicebear.com/7.x/initials/svg?seed=BX', campaignCount: 1 }
];

const INITIAL_CONTENT: ContentLibraryItem[] = [
  {
    id: 'c-1',
    brandId: 'br-1',
    title: '秋季编程入门营',
    description: '针对零基础学员的 Python 入门课程。',
    imageUrl: 'https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=800',
    contentBlocks: [
      { id: 'cb1', type: 'heading', value: '为什么要学习编程？' },
      { id: 'cb2', type: 'text', value: '在 AI 时代，编程是与机器沟通的基本技能。' }
    ],
    formSchema: [
      { id: 'name', label: '姓名', type: 'string', required: true },
      { id: 'phone', label: '手机号', type: 'number', required: true }
    ]
  }
];

const INITIAL_TASKS: MarketingTask[] = [
  { 
    id: 'tsk-1', 
    contentId: 'c-1',
    title: 'Python训练营裂变招募', 
    type: TaskType.PAYMENT, 
    rewardAmount: 20.0, 
    entryFee: 9.9,
    status: 'ACTIVE', 
    description: '推荐成交奖励20元。',
    socialPost: '我发现一个很棒的 Python 课程，只要 9.9 元，一起来学吧！',
    createdAt: new Date().toISOString()
  }
];

const MOCK_USER: UserProfile = {
  id: 'ref-999',
  name: '张经理',
  email: 'zhang@dmh.com',
  phone: '13812345678',
  avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Admin',
  accumulatedRewards: 450.0,
  referralCount: 42,
  totalShares: 156,
  totalSpent: 1200,
  history: [],
  joinDate: '2023-01-01'
};

class MockApiService {
  private contents: ContentLibraryItem[] = [];
  private tasks: MarketingTask[] = [];
  private orders: Order[] = [];
  private profiles: UserProfile[] = [];
  private industries: Industry[] = [];
  private brands: Brand[] = [];
  private withdrawals: Withdrawal[] = [];

  constructor() {
    this.loadFromStorage();
  }

  private loadFromStorage() {
    const data = localStorage.getItem('dmh_data_final_v1');
    if (data) {
      const parsed = JSON.parse(data);
      this.contents = parsed.contents || INITIAL_CONTENT;
      this.tasks = parsed.tasks || INITIAL_TASKS;
      this.orders = parsed.orders || [];
      this.profiles = parsed.profiles || [MOCK_USER];
      this.industries = parsed.industries || INITIAL_INDUSTRIES;
      this.brands = parsed.brands || INITIAL_BRANDS;
      this.withdrawals = parsed.withdrawals || [];
    } else {
      this.contents = INITIAL_CONTENT;
      this.tasks = INITIAL_TASKS;
      this.profiles = [MOCK_USER];
      this.industries = INITIAL_INDUSTRIES;
      this.brands = INITIAL_BRANDS;
      this.saveToStorage();
    }
  }

  private saveToStorage() {
    localStorage.setItem('dmh_data_final_v1', JSON.stringify({
      contents: this.contents,
      tasks: this.tasks,
      orders: this.orders,
      profiles: this.profiles,
      industries: this.industries,
      brands: this.brands,
      withdrawals: this.withdrawals
    }));
  }

  async listIndustries() { return this.industries; }
  async listBrands() { return this.brands; }
  async listContentLibrary() { return this.contents; }
  async listTasks() { 
    return this.tasks.map(t => ({
      ...t,
      content: this.contents.find(c => c.id === t.contentId)
    }));
  }
  async listOrders() { 
    return this.orders.map(o => ({
      ...o,
      taskTitle: this.tasks.find(t => t.id === o.taskId)?.title || '未知任务'
    }));
  }
  async listWithdrawals() { return this.withdrawals; }
  async listProfiles() { return this.profiles; }

  async saveContent(item: ContentLibraryItem) {
    const idx = this.contents.findIndex(c => c.id === item.id);
    if (idx >= 0) this.contents[idx] = item;
    else { item.id = `c-${Date.now()}`; this.contents.push(item); }
    this.saveToStorage();
    return item;
  }

  async saveTask(task: MarketingTask) {
    const idx = this.tasks.findIndex(t => t.id === task.id);
    if (idx >= 0) this.tasks[idx] = task;
    else { task.id = `tsk-${Date.now()}`; task.createdAt = new Date().toISOString(); this.tasks.push(task); }
    this.saveToStorage();
    return task;
  }

  async approveWithdrawal(id: string, status: 'APPROVED' | 'REJECTED') {
    const wd = this.withdrawals.find(w => w.id === id);
    if (wd) { wd.status = status; this.saveToStorage(); }
    return wd;
  }

  async createOrder(taskId: string, formData: any, referrerId?: string) {
    const task = this.tasks.find(t => t.id === taskId);
    if (!task) throw new Error('Task not found');
    const order: Order = {
      id: `ord-${Date.now()}`,
      taskId,
      userId: `u-${Math.random().toString(36).substr(2, 6)}`,
      referrerId,
      status: 'PAID',
      amount: task.entryFee,
      formData,
      createdAt: new Date().toISOString()
    };
    this.orders.unshift(order);
    this.saveToStorage();
    return order;
  }
}

export const api = new MockApiService();
