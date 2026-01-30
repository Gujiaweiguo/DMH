// 分销商管理视图
import { ref, onMounted, h, computed, defineComponent } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { distributorApi } from '../services/distributorApi';

export const DistributorManagementView = defineComponent({
  props: {
    brandId: Number,
    readOnly: Boolean,
    isPlatformAdmin: Boolean
  },
  setup(props) {
    // props 本身就是响应式的，不需要用 ref() 包裹
    const readOnly = props.readOnly;
    const isPlatformAdmin = props.isPlatformAdmin;
    const distributors = ref<any[]>([]);
    const loading = ref(false);
    const searchKeyword = ref('');
    const statusFilter = ref('');
    const levelFilter = ref(0);
    const brandFilter = ref(0); // 品牌筛选（平台管理员使用）
    const page = ref(1);
    const pageSize = ref(50);
    const total = ref(0);

    // 加载分销商列表
    const loadDistributors = async () => {
      try {
        loading.value = true;
        let response: any = null;
        
        console.log('[DistributorManagement] 开始加载数据...');
        
        // 临时：总是使用全局 API
        response = await distributorApi.getGlobalDistributors(
          undefined,
          statusFilter.value || undefined,
          page.value,
          pageSize.value
        );
        
        console.log('[DistributorManagement] API 响应:', response);
        
        // 处理响应
        distributors.value = response.distributors || [];
        total.value = response.total || 0;
        
        console.log('[DistributorManagement] 分销商数量:', distributors.value.length);
      } catch (error) {
        console.error('[DistributorManagement] 加载失败:', error);
        distributors.value = [];
        total.value = 0;
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

    // 获取状态样式和标签
    const getStatusInfo = (status: string) => {
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
      return {
        style: styles[status as keyof typeof styles] || 'bg-slate-100 text-slate-800',
        label: labels[status as keyof typeof labels] || status
      };
    };

    // 计算统计数据
    const stats = computed(() => {
      const all = distributors.value;
      return {
        total: all.length,
        active: all.filter(d => d.status === 'active').length,
        suspended: all.filter(d => d.status === 'suspended').length,
        pending: all.filter(d => d.status === 'pending').length
      };
    });
  
    // 统计卡片渲染函数
    const renderStatCards = (viewModel: any) => {
      return h('div', { class: 'grid grid-cols-1 md:grid-cols-4 gap-6 mb-6' }, [
        h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
          h('div', { class: 'flex items-center gap-4 mb-4' }, [
            h('div', { class: 'w-12 h-12 bg-indigo-600 text-white rounded-2xl flex items-center justify-center' },
              h((LucideIcons as any).Users, { size: 24 })
            ),
          ]),
          h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, '总数'),
          h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, String(viewModel.total.value || 0))
        ]),
        h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
          h('div', { class: 'flex items-center gap-4 mb-4' }, [
            h('div', { class: 'w-12 h-12 bg-emerald-600 text-white rounded-2xl flex items-center justify-center' },
              h((LucideIcons as any).UserCheck, { size: 24 })
            ),
          ]),
          h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, '正常'),
          h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, String(viewModel.stats.value?.active || 0))
        ]),
        h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
          h('div', { class: 'flex items-center gap-4 mb-4' }, [
            h('div', { class: 'w-12 h-12 bg-rose-600 text-white rounded-2xl flex items-center justify-center' },
              h((LucideIcons as any).UserX, { size: 24 })
            ),
          ]),
          h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, '暂停'),
          h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, String(viewModel.stats.value?.suspended || 0))
        ]),
        h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
          h('div', { class: 'flex items-center gap-4 mb-4' }, [
            h('div', { class: 'w-12 h-12 bg-amber-600 text-white rounded-2xl flex items-center justify-center' },
              h((LucideIcons as any).Clock, { size: 24 })
            ),
          ]),
          h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, '待审核'),
          h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, String(viewModel.stats.value?.pending || 0))
        ])
      ]);
    };
  
    // 计算过滤后的分销商列表
    const filteredDistributors = computed(() => {
      let result = distributors.value;
      
      console.log('[DistributorManagement] 开始过滤 - 原始数据:', result.length);
      console.log('[DistributorManagement] 筛选条件:', {
        statusFilter: statusFilter.value,
        levelFilter: levelFilter.value,
        searchKeyword: searchKeyword.value
      });

      // 状态筛选
      if (statusFilter.value) {
        const before = result.length;
        result = result.filter(d => d.status === statusFilter.value);
        console.log(`[DistributorManagement] 状态筛选: ${before} -> ${result.length}`);
      }

      // 级别筛选
      if (levelFilter.value > 0) {
        const before = result.length;
        result = result.filter(d => d.level === levelFilter.value);
        console.log(`[DistributorManagement] 级别筛选: ${before} -> ${result.length}`);
      }

      // 搜索筛选
      if (searchKeyword.value.trim()) {
        const before = result.length;
        const keyword = searchKeyword.value.trim().toLowerCase();
        result = result.filter(d => 
          (d.username && d.username.toLowerCase().includes(keyword)) ||
          (d.brandName && d.brandName.toLowerCase().includes(keyword))
        );
        console.log(`[DistributorManagement] 搜索筛选: ${before} -> ${result.length}`);
      }

      console.log('[DistributorManagement] 最终结果:', result.length);
      return result;
    });

    onMounted(() => {
      console.log('[DistributorManagement] 组件已挂载，开始加载数据');
      loadDistributors();
    });

    return () => renderDistributorManagementView({
      readOnly,
      isPlatformAdmin,
      distributors,
      loading,
      searchKeyword,
      statusFilter,
      levelFilter,
      brandFilter,
      page,
      pageSize,
      total,
      stats,
      filteredDistributors,
      loadDistributors,
      goPrev,
      goNext,
      getStatusInfo,
      renderStatCards
    });
  }
});

// 渲染管理视图
const renderDistributorManagementView = (viewModel: any) => {
  // 确保 viewModel 包含 stats
  if (!viewModel.stats) {
    console.error('[renderDistributorManagementView] viewModel.stats 未定义');
  }
  
  const tableHeaders = ['ID', '分销用户', '品牌', '级别', '上级', '累计收益', '下级数', '状态', '加入时间', '审批时间'];

  return h('div', { class: 'space-y-6' }, [
    // 统计卡片
    renderStatCards(viewModel),
    
    // 头部
    h('div', { class: 'flex justify-between items-center' }, [
      h('div', [
        h('h2', { class: 'text-2xl font-black text-slate-900' }, '分销监控'),
        h('p', { class: 'text-slate-400 text-sm mt-1' }, '查看拥有分销权限的普通用户分销情况（只读）')
      ]),
      h('button', {
        onClick: viewModel.loadDistributors,
        class: 'bg-indigo-600 text-white px-6 py-3 rounded-2xl font-bold hover:bg-indigo-700 transition-colors flex items-center gap-2'
      }, [
        h(LucideIcons.RefreshCw, { size: 18 }),
        '刷新'
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
        viewModel.filteredDistributors.value.length === 0 ?
          h('div', { class: 'p-12 text-center' }, [
            h(LucideIcons.Users, { size: 48, class: 'mx-auto text-slate-300 mb-4' }),
            h('p', { class: 'text-slate-500 text-lg' }, '暂无分销商')
          ]) :
          h('table', { class: 'w-full text-left' }, [
            h('thead', { class: 'bg-slate-50' }, [
              h('tr', tableHeaders.map(th => h('th', { class: 'px-6 py-4 text-xs font-black text-slate-400 uppercase tracking-widest' }, th)))
            ]),
            h('tbody', viewModel.filteredDistributors.value.map(d => h('tr', { class: 'border-b border-slate-50 last:border-0 hover:bg-slate-50/40' }, [
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
              h('td', { class: 'px-6 py-4' }, [
                (() => {
                  const statusInfo = viewModel.getStatusInfo(d.status);
                  return h('span', { class: `px-2 py-1 rounded-lg text-xs font-bold ${statusInfo.style}` }, statusInfo.label);
                })()
              ]),
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
