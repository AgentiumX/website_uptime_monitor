import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/dashboard/DashboardView.vue') },
      { path: 'monitors', name: 'MonitorList', component: () => import('@/views/monitor/MonitorList.vue') },
      { path: 'monitors/create', name: 'MonitorCreate', component: () => import('@/views/monitor/MonitorForm.vue') },
      { path: 'monitors/:id/edit', name: 'MonitorEdit', component: () => import('@/views/monitor/MonitorForm.vue') },
      { path: 'monitors/:id', name: 'MonitorDetail', component: () => import('@/views/monitor/MonitorDetail.vue') },
      { path: 'alerts/channels', name: 'AlertChannels', component: () => import('@/views/alert/AlertChannelList.vue') },
      { path: 'alerts/history', name: 'AlertHistory', component: () => import('@/views/alert/AlertHistory.vue') },
      { path: 'agents', name: 'AgentList', component: () => import('@/views/agent/AgentList.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    next('/login')
  } else if (to.meta.guest && auth.isLoggedIn) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
