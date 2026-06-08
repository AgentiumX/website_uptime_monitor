<template>
  <div class="dashboard-view">
    <h1 class="page-title">仪表盘</h1>

    <!-- Statistics cards -->
    <div class="stat-row">
      <StatCard
        label="监控总数"
        :value="overview.total_monitors ?? 0"
      />
      <StatCard
        label="在线率"
        :value="`${(overview.uptime_rate ?? 0).toFixed(2)}`"
        suffix="%"
        :value-color="(overview.uptime_rate ?? 0) > 99 ? '#34C759' : '#FF9500'"
      />
      <StatCard
        label="活跃告警"
        :value="overview.active_alerts ?? 0"
        :value-color="(overview.active_alerts ?? 0) > 0 ? '#FF3B30' : '#34C759'"
      />
      <StatCard
        label="平均响应时间"
        :value="overview.avg_latency_ms ?? 0"
        suffix="ms"
      />
    </div>

    <!-- Charts row -->
    <div class="charts-row">
      <el-card shadow="never" class="chart-card chart-wide">
        <template #header>
          <div class="card-header">
            <span class="card-title">可用率趋势</span>
            <div class="range-btns">
              <el-button
                v-for="r in ranges"
                :key="r.value"
                :type="range === r.value ? 'primary' : 'default'"
                size="small"
                @click="range = r.value"
              >
                {{ r.label }}
              </el-button>
            </div>
          </div>
        </template>
        <UptimeChart :range="range" />
      </el-card>
      <el-card shadow="never" class="chart-card chart-narrow">
        <template #header>
          <span class="card-title">健康度分布</span>
        </template>
        <HealthDonut
          :healthy="overview.healthy_count ?? 0"
          :warning="overview.warning_count ?? 0"
          :danger="overview.danger_count ?? 0"
        />
      </el-card>
    </div>

    <!-- Bottom row -->
    <div class="bottom-row">
      <el-card shadow="never" class="bottom-card">
        <template #header>
          <span class="card-title">最近告警</span>
        </template>
        <div class="alert-list">
          <div
            v-for="alert in recentAlerts"
            :key="alert.id"
            class="alert-item"
            :class="`alert-item--${alert.status}`"
          >
            <StatusDot :status="alert.status === 'firing' ? 'offline' : 'online'" />
            <div class="alert-info">
              <span class="alert-name">{{ alert.monitor_name }}</span>
              <span class="alert-detail">{{ alert.detail }}</span>
            </div>
            <span class="alert-time">{{ formatTime(alert.triggered_at) }}</span>
          </div>
          <el-empty v-if="recentAlerts.length === 0" description="暂无告警" :image-size="60" />
        </div>
      </el-card>
      <el-card shadow="never" class="bottom-card">
        <template #header>
          <span class="card-title">TOP5 监控站点</span>
        </template>
        <div class="top-list">
          <div
            v-for="item in topMonitors"
            :key="item.id"
            class="top-item"
          >
            <StatusDot :status="item.status === 'up' ? 'online' : 'offline'" />
            <span class="top-name">{{ item.name }}</span>
            <span class="top-uptime">{{ item.uptime_rate.toFixed(2) }}%</span>
            <span class="top-latency">{{ item.avg_latency_ms }}ms</span>
          </div>
          <el-empty v-if="topMonitors.length === 0" description="暂无数据" :image-size="60" />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import StatCard from '@/components/common/StatCard.vue'
import StatusDot from '@/components/common/StatusDot.vue'
import UptimeChart from '@/components/charts/UptimeChart.vue'
import HealthDonut from '@/components/charts/HealthDonut.vue'
import { getDashboardOverview } from '@/api/dashboard'
import { listAlertHistory } from '@/api/alert'
import { listMonitors } from '@/api/monitor'

const range = ref<'24h' | '7d' | '30d'>('24h')
const ranges = [
  { label: '24h', value: '24h' as const },
  { label: '7d', value: '7d' as const },
  { label: '30d', value: '30d' as const },
]

interface Overview {
  total_monitors: number
  uptime_rate: number
  active_alerts: number
  avg_latency_ms: number
  healthy_count: number
  warning_count: number
  danger_count: number
}

const overview = ref<Overview>({
  total_monitors: 0,
  uptime_rate: 0,
  active_alerts: 0,
  avg_latency_ms: 0,
  healthy_count: 0,
  warning_count: 0,
  danger_count: 0,
})

interface AlertItem {
  id: number
  monitor_name: string
  detail: string
  status: string
  triggered_at: string
}

const recentAlerts = ref<AlertItem[]>([])

interface TopMonitor {
  id: number
  name: string
  status: string
  uptime_rate: number
  avg_latency_ms: number
}

const topMonitors = ref<TopMonitor[]>([])

const formatTime = (t: string) => {
  if (!t) return ''
  const d = new Date(t)
  return `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

onMounted(async () => {
  try {
    const res = await getDashboardOverview()
    overview.value = res.data
  } catch {
    /* API not ready */
  }

  try {
    const res = await listAlertHistory({ page: 1, page_size: 5, status: 'firing' })
    recentAlerts.value = res.data?.list || []
  } catch {
    /* API not ready */
  }

  try {
    const res = await listMonitors({ page: 1, page_size: 5 })
    topMonitors.value = (res.data?.list || []).map((m: Record<string, unknown>) => ({
      id: m.id,
      name: m.name as string,
      status: m.last_status === 'up' ? 'up' : 'down',
      uptime_rate: (m.uptime_rate as number) || 0,
      avg_latency_ms: (m.avg_latency_ms as number) || 0,
    }))
  } catch {
    /* API not ready */
  }
})
</script>

<style lang="scss" scoped>
.dashboard-view {
  max-width: 1200px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 24px;
}

.stat-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.charts-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
  margin-bottom: 24px;
}

.bottom-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.chart-card,
.bottom-card {
  :deep(.el-card__header) {
    padding: 16px 20px;
    border-bottom: 1px solid var(--color-border);
  }

  :deep(.el-card__body) {
    padding: 20px;
  }
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.range-btns {
  display: flex;
  gap: 6px;
}

.alert-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.alert-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 10px;
  background: #f9f9fb;

  &--firing {
    background: rgba(255, 59, 48, 0.06);
  }

  &--resolved {
    background: rgba(52, 199, 89, 0.06);
  }
}

.alert-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.alert-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.alert-detail {
  font-size: 12px;
  color: var(--color-text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.alert-time {
  font-size: 12px;
  color: var(--color-text-secondary);
  white-space: nowrap;
}

.top-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.top-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 10px;
  background: #f9f9fb;
}

.top-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.top-uptime {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-success);
}

.top-latency {
  font-size: 13px;
  color: var(--color-text-secondary);
  min-width: 56px;
  text-align: right;
}
</style>
