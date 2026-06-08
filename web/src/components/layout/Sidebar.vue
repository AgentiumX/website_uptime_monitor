<template>
  <aside class="sidebar">
    <div class="sidebar-logo" @click="router.push('/dashboard')">
      <span class="logo-uptime">Uptime</span>
      <span class="logo-monitor">Monitor</span>
    </div>
    <el-menu
      :default-active="activeMenu"
      :router="true"
      class="sidebar-menu"
    >
      <el-menu-item index="/dashboard">
        <el-icon><Odometer /></el-icon>
        <span>仪表盘</span>
      </el-menu-item>
      <el-menu-item index="/monitors">
        <el-icon><Monitor /></el-icon>
        <span>监控列表</span>
      </el-menu-item>
      <el-sub-menu index="alerts">
        <template #title>
          <el-icon><BellFilled /></el-icon>
          <span>告警管理</span>
        </template>
        <el-menu-item index="/alerts/channels">告警通道</el-menu-item>
        <el-menu-item index="/alerts/history">告警历史</el-menu-item>
      </el-sub-menu>
      <el-menu-item index="/agents">
        <el-icon><Position /></el-icon>
        <span>探测点</span>
      </el-menu-item>
    </el-menu>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const activeMenu = computed(() => {
  const path = route.path
  // Map child routes to parent menu items
  if (path.startsWith('/monitors')) return '/monitors'
  if (path.startsWith('/alerts')) return path
  if (path.startsWith('/agents')) return '/agents'
  return '/dashboard'
})
</script>

<style lang="scss" scoped>
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 240px;
  background: var(--color-bg-card);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  z-index: 20;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  height: 64px;
  padding: 0 24px;
  cursor: pointer;
  border-bottom: 1px solid var(--color-border);
  gap: 6px;
}

.logo-uptime {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.logo-monitor {
  font-size: 20px;
  font-weight: 300;
  color: var(--color-primary);
}

.sidebar-menu {
  flex: 1;
  padding-top: 8px;

  .el-menu-item {
    margin: 2px 8px;
    border-radius: 10px;
    height: 44px;
    line-height: 44px;
  }

  .el-menu-item.is-active {
    background: rgba(0, 113, 227, 0.08);
    color: var(--color-primary);
  }
}
</style>
