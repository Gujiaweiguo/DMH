// 分销商管理视图
import { ref, onMounted } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { distributorApi } from '../services/distributorApi';

export const DistributorManagementView = (props: { brandId?: number; readOnly?: boolean }) => {
  const readOnly = ref(!!props.readOnly);
  const distributors = ref<any[]>([]);
  const loading = ref(false);
  const searchKeyword = ref('');
  const statusFilter = ref('');
  const levelFilter = ref(0);
  const page = ref(1);
  const pageSize = ref(50);
  const total = ref(0);

  // 加载分销商列表
  const loadDistributors = async () => {
    try {
      loading.value = true;
      const response: any = await distributorApi.getDistributors(props.brandId || 1, {
        keyword: searchKeyword.value,
        status: statusFilter.value,
        level: levelFilter.value,
        page: page.value,
        pageSize: pageSize.value
      });
      if (response.code === 200) {
        distributors.value = response.data.distributors || [];
        total.value = response.data.total || 0;
      }
    } catch (error) {
      console.error('加载分销商列表失败:', error);
    } finally {
      loading.value = false;
    }
  };

  const goPrev = async () => {
    if (page.value <= 1) return;
    page.value -= 1;
    await loadDistributors();
  };

  const goNext = async () => {
    const maxPage = Math.max(1, Math.ceil((total.value || 0) / (pageSize.value || 1)));
    if (page.value >= maxPage) return;
    page.value += 1;
    await loadDistributors();
  };

  // 获取状态标签
  const getStatusBadge = (status: string) => {
    const styles = {
      active: 'bg-emerald-100 text-emerald-800',
      suspended: 'bg-rose-100 text-rose-800',
      pending: 'bg-amber-100 text-amber-800'
    };
    const labels = {
      active: '正常',
      suspended: '已暂停',
      pending: '待审核'
    };
    const style = styles[status as keyof typeof styles] || 'bg-slate-100 text-slate-800';
    const label = labels[status as keyof typeof labels] || status;
    return h('span', { class: `px-2 py-1 rounded-lg text-xs font-bold ${style}` }, label);
  };

  onMounted(() => {
    loadDistributors();
  });

  return {
    readOnly,
    distributors,
    loading,
    searchKeyword,
    statusFilter,
    levelFilter,
    page,
    pageSize,
    total,
    loadDistributors,
    goPrev,
    goNext,
    getStatusBadge
  };
};

// 渲染管理视图
export const renderDistributorManagementView = (viewModel: ReturnType<typeof DistributorManagementView>) => {
  const { h } = (window as any).Vue || { h: () => null };

  const tableHeaders = ['ID', '分销用户', '品牌', '级别', '上级', '累计收益', '下级数', '状态', '加入时间', '审批时间'];

  return h('div', { class: 'space-y-6' }, [
    // 头部
    h('div', { class: 'flex justify-between items-center' }, [
      h('div', [
        h('h2', { class: 'text-2xl font-black text-slate-900' }, '分销监控'),
        h('p', { class: 'text-slate-400 text-sm mt-1' }, '查看拥有分销权限的普通用户分销情况（只读）')
      ]),
      h('div', { class: 'flex items-center gap-3' }, [
        h('span', { class: 'text-xs font-black text-slate-400 uppercase tracking-widest' }, `总数 ${viewModel.total.value || 0}`),
        h('button', {
          onClick: viewModel.loadDistributors,
          class: 'bg-indigo-600 text-white px-6 py-3 rounded-2xl font-bold hover:bg-indigo-700 transition-colors flex items-center gap-2'
        }, [
          h(LucideIcons.RefreshCw, { size: 18 }),
          '刷新'
        ])
      ])
    ]),

    // 筛选条件
    h('div', { class: 'flex gap-4 items-center' }, [
      h('input', {
        value: viewModel.searchKeyword.value,
        onInput: (e: any) => { viewModel.searchKeyword.value = e.target.value },
        placeholder: '搜索用户名...',
        class: 'flex-1 max-w-xs px-4 py-2 border border-slate-200 rounded-xl'
      }),
      h('select', {
        value: viewModel.statusFilter.value,
        onChange: (e: any) => { viewModel.statusFilter.value = e.target.value },
        class: 'px-4 py-2 border border-slate-200 rounded-xl'
      }, [
        h('option', { value: '' }, '全部状态'),
        h('option', { value: 'active' }, '正常'),
        h('option', { value: 'suspended' }, '已暂停'),
        h('option', { value: 'pending' }, '待审核')
      ]),
      h('select', {
        value: viewModel.levelFilter.value,
        onChange: (e: any) => { viewModel.levelFilter.value = parseInt(e.target.value) },
        class: 'px-4 py-2 border border-slate-200 rounded-xl'
      }, [
        h('option', { value: 0 }, '全部级别'),
        h('option', { value: 1 }, '一级'),
        h('option', { value: 2 }, '二级'),
        h('option', { value: 3 }, '三级')
      ]),
      h('button', {
        onClick: () => { viewModel.page.value = 1; viewModel.loadDistributors(); },
        class: 'px-6 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700'
      }, '搜索')
    ]),

    // 分销商列表
    h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden shadow-sm' }, [
      viewModel.loading.value ?
        h('div', { class: 'p-12 text-center' }, ['加载中...']) :
        viewModel.distributors.value.length === 0 ?
          h('div', { class: 'p-12 text-center' }, [
            h(LucideIcons.Users, { size: 48, class: 'mx-auto text-slate-300 mb-4' }),
            h('p', { class: 'text-slate-500 text-lg' }, '暂无分销商')
          ]) :
          h('table', { class: 'w-full text-left' }, [
            h('thead', { class: 'bg-slate-50' }, [
              h('tr', tableHeaders.map(th => h('th', { class: 'px-6 py-4 text-xs font-black text-slate-400 uppercase tracking-widest' }, th)))
            ]),
            h('tbody', viewModel.distributors.value.map(d => h('tr', { class: 'border-b border-slate-50 last:border-0 hover:bg-slate-50/40' }, [
              h('td', { class: 'px-6 py-4 text-sm text-slate-400 font-mono' }, String(d.id)),
              h('td', { class: 'px-6 py-4' }, [
                h('div', { class: 'text-sm font-bold text-slate-900' }, d.username || `用户${d.userId}`)
              ]),
              h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, d.brandName || `品牌${d.brandId}`),
              h('td', { class: 'px-6 py-4 text-sm' }, [
                h('span', {
                  class: `px-2 py-1 rounded-lg text-xs font-bold ${
                    d.level === 1 ? 'bg-purple-100 text-purple-800' :
                    d.level === 2 ? 'bg-blue-100 text-blue-800' :
                    'bg-cyan-100 text-cyan-800'
                  }`
                }, `${d.level}级`)
              ]),
              h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, d.parentName || '-'),
              h('td', { class: 'px-6 py-4 text-sm font-medium text-emerald-600' }, `¥${(d.totalEarnings ?? 0).toFixed(2)}`),
              h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, String(d.subordinatesCount)),
              h('td', { class: 'px-6 py-4' }, [viewModel.getStatusBadge(d.status)]),
              h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, d.createdAt || '-'),
              h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, d.approvedAt || '-')
            ])))
          ])
    ]),

    // 分页
    h('div', { class: 'flex items-center justify-between' }, [
      h('div', { class: 'text-sm text-slate-500' }, `第 ${viewModel.page.value} 页 · 每页 ${viewModel.pageSize.value} 条`),
      h('div', { class: 'flex gap-2' }, [
        h('button', {
          onClick: viewModel.goPrev,
          class: 'px-4 py-2 text-slate-700 bg-slate-100 rounded-lg hover:bg-slate-200 disabled:opacity-50',
          disabled: viewModel.page.value <= 1
        }, '上一页'),
        h('button', {
          onClick: viewModel.goNext,
          class: 'px-4 py-2 text-slate-700 bg-slate-100 rounded-lg hover:bg-slate-200 disabled:opacity-50',
          disabled: viewModel.total.value <= viewModel.page.value * viewModel.pageSize.value
        }, '下一页')
      ])
    ])
  ]);
};
