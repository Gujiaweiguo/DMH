import { createRouter, createWebHistory } from 'vue-router'

// 用户端页面
import CampaignList from '../views/CampaignList.vue'
import CampaignDetail from '../views/CampaignDetail.vue'
import CampaignForm from '../views/CampaignForm.vue'
import MyOrders from '../views/MyOrders.vue'
import Success from '../views/Success.vue'

// 品牌人员页面
import BrandLogin from '../views/brand/Login.vue'
import BrandDashboard from '../views/brand/Dashboard.vue'
import BrandCampaigns from '../views/brand/Campaigns.vue'
import BrandCampaignEditor from '../views/brand/CampaignEditorVant.vue'
import BrandCampaignPageDesigner from '../views/brand/CampaignPageDesigner.vue'
import BrandMaterials from '../views/brand/Materials.vue'
import BrandOrders from '../views/brand/Orders.vue'
import BrandPromoters from '../views/brand/Promoters.vue'
import BrandAnalytics from '../views/brand/Analytics.vue'
import BrandSettings from '../views/brand/Settings.vue'

const routes = [
  // 用户端路由
  {
    path: '/',
    name: 'Home',
    component: CampaignList
  },
  {
    path: '/campaign/:id',
    name: 'CampaignDetail',
    component: CampaignDetail
  },
  {
    path: '/campaign/:id/form',
    name: 'CampaignForm',
    component: CampaignForm
  },
  {
    path: '/orders',
    name: 'MyOrders',
    component: MyOrders
  },
  {
    path: '/success',
    name: 'Success',
    component: Success
  },

  // 品牌人员路由
  {
    path: '/brand/login',
    name: 'BrandLogin',
    component: BrandLogin
  },
  {
    path: '/brand',
    redirect: '/brand/dashboard'
  },
  {
    path: '/brand/dashboard',
    name: 'BrandDashboard',
    component: BrandDashboard,
    meta: { requiresAuth: true, role: 'brand_admin' }
  },
  {
    path: '/brand/campaigns',
    name: 'BrandCampaigns',
    component: BrandCampaigns,
    meta: { requiresAuth: true, role: 'brand_admin' }
  },
  {
    path: '/brand/campaigns/create',
    name: 'BrandCampaignCreate',
    component: BrandCampaignEditor,
    meta: { requiresAuth: true, role: 'brand_admin' }
  },
  {
    path: '/brand/campaigns/edit/:id',
    name: 'BrandCampaignEdit',
    component: BrandCampaignEditor,
    meta: { requiresAuth: true, role: 'brand_admin' }
  },
  {
    path: '/brand/campaigns/:id/page-design',
    name: 'BrandCampaignPageDesigner',
    component: BrandCampaignPageDesigner,
    meta: { requiresAuth: true, role: 'brand_admin' }
  },
  {
    path: '/brand/materials',
    name: 'BrandMaterials',
    component: BrandMaterials,
    meta: { requiresAuth: true, role: 'brand_admin' }
  },
  {
    path: '/brand/orders',
    name: 'BrandOrders',
    component: BrandOrders,
    meta: { requiresAuth: true, role: 'brand_admin' }
  },
  {
    path: '/brand/promoters',
    name: 'BrandPromoters',
    component: BrandPromoters,
    meta: { requiresAuth: true, role: 'brand_admin' }
  },
  {
    path: '/brand/analytics',
    name: 'BrandAnalytics',
    component: BrandAnalytics,
    meta: { requiresAuth: true, role: 'brand_admin' }
  },
  {
    path: '/brand/settings',
    name: 'BrandSettings',
    component: BrandSettings,
    meta: { requiresAuth: true, role: 'brand_admin' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('dmh_token')
  const userRole = localStorage.getItem('dmh_user_role')
  
  if (to.meta.requiresAuth) {
    if (!token) {
      next('/brand/login')
      return
    }
    
    if (to.meta.role && userRole !== to.meta.role) {
      next('/')
      return
    }
  }
  
  next()
})

export default router