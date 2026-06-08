<template>
  <div class="app-layout">
    <Sidebar />
    <div class="main-area">
      <header class="app-header">
        <div class="header-left">
          <span class="page-title">{{ pageTitle }}</span>
        </div>
        <div class="header-right">
          <el-badge :value="0" :hidden="true" class="notification-badge">
            <el-button :icon="Bell" circle size="default" />
          </el-badge>
          <el-dropdown trigger="click" @command="handleCommand">
            <span class="user-info">
              <el-icon><UserFilled /></el-icon>
              <span class="username">Admin</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>
      <main class="app-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Bell } from '@element-plus/icons-vue'
import Sidebar from './Sidebar.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const pageTitles: Record<string, string> = {
  Dashboard: '仪表盘',
  MonitorList: '监控列表',
  MonitorCreate: '新建监控',
  MonitorEdit: '编辑监控',
  MonitorDetail: '监控详情',
  AlertChannels: '告警通道',
  AlertHistory: '告警历史',
  AgentList: '探测点管理',
}

const pageTitle = computed(() => {
  return pageTitles[route.name as string] || '仪表盘'
})

const handleCommand = async (command: string) => {
  if (command === 'logout') {
    await auth.logout()
    router.push('/login')
  }
}
</script>

<style lang="scss" scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}

.main-area {
  flex: 1;
  margin-left: 240px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 32px;
  background: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-left {
  .page-title {
    font-size: 20px;
    font-weight: 600;
    color: var(--color-text-primary);
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  color: var(--color-text-secondary);
  font-size: 14px;

  .username {
    font-weight: 500;
  }

  &:hover {
    color: var(--color-primary);
  }
}

.app-content {
  flex: 1;
  padding: 32px;
}
</style>
