<template>
  <div class="monitor-detail-view">
    <!-- Breadcrumb -->
    <el-breadcrumb separator="/" class="breadcrumb">
      <el-breadcrumb-item :to="{ path: '/monitors' }">监控列表</el-breadcrumb-item>
      <el-breadcrumb-item>{{ monitor.name || '监控详情' }}</el-breadcrumb-item>
    </el-breadcrumb>

    <div class="page-header">
      <h1 class="page-title">{{ monitor.name }}</h1>
      <el-switch
        v-model="monitor.enabled"
        active-text="启用"
        inactive-text="禁用"
        @change="handleToggle"
      />
    </div>

    <!-- Metrics cards -->
    <div class="metrics-row">
      <div class="metric-card">
        <div class="metric-label">最近状态码</div>
        <div class="metric-value" :class="{ 'text-danger': monitor.last_status_code >= 400 }">
          {{ monitor.last_status_code || '-' }}
        </div>
      </div>
      <div class="metric-card">
        <div class="metric-label">最近延迟</div>
        <div class="metric-value">{{ monitor.last_latency_ms || '-' }}<span v-if="monitor.last_latency_ms">ms</span></div>
      </div>
      <div class="metric-card">
        <div class="metric-label">SSL 过期天数</div>
        <div class="metric-value" :class="{ 'text-warning': monitor.ssl_expiry_days !== null && monitor.ssl_expiry_days <= 30 }">
          {{ monitor.ssl_expiry_days ?? '-' }}<span v-if="monitor.ssl_expiry_days != null">天</span>
        </div>
      </div>
      <div class="metric-card">
        <div class="metric-label">内容匹配</div>
        <div class="metric-value">
          <el-tag :type="monitor.last_content_matched ? 'success' : 'danger'" size="small">
            {{ monitor.last_content_matched ? '匹配' : '未匹配' }}
          </el-tag>
        </div>
      </div>
    </div>

    <!-- Charts -->
    <div class="charts-row">
      <el-card shadow="never" class="chart-card">
        <template #header>
          <span class="card-title">可用率趋势</span>
        </template>
        <UptimeChart range="7d" />
      </el-card>
      <el-card shadow="never" class="chart-card">
        <template #header>
          <span class="card-title">延迟趋势</span>
        </template>
        <LatencyChart :monitor-id="monitorId" />
      </el-card>
    </div>

    <!-- Alert timeline -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <span class="card-title">告警时间线</span>
      </template>
      <el-timeline v-if="alerts.length > 0">
        <el-timeline-item
          v-for="alert in alerts"
          :key="alert.id"
          :timestamp="alert.triggered_at"
          :color="alert.status === 'firing' ? '#FF3B30' : '#34C759'"
          placement="top"
        >
          <div class="alert-timeline-item">
            <el-tag :type="alert.status === 'firing' ? 'danger' : 'success'" size="small">
              {{ alert.status === 'firing' ? '告警中' : '已恢复' }}
            </el-tag>
            <span class="alert-type">{{ alert.alert_type }}</span>
            <span class="alert-detail">{{ alert.detail }}</span>
          </div>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else description="暂无告警记录" :image-size="60" />
    </el-card>

    <!-- History table -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <span class="card-title">历史记录</span>
      </template>
      <el-table :data="history" stripe style="width: 100%">
        <el-table-column prop="timestamp" label="时间" width="180" />
        <el-table-column prop="agent_name" label="探测点" width="120" />
        <el-table-column prop="status_code" label="状态码" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status_code >= 400 ? 'danger' : 'success'" size="small">
              {{ row.status_code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration_ms" label="延迟" width="90" align="center">
          <template #default="{ row }">{{ row.duration_ms }}ms</template>
        </el-table-column>
        <el-table-column prop="content_matched" label="内容匹配" width="90" align="center">
          <template #default="{ row }">
            <el-icon :color="row.content_matched ? '#34C759' : '#FF3B30'">
              <CircleCheck v-if="row.content_matched" />
              <CircleClose v-else />
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="success" label="结果" width="80" align="center">
          <template #default="{ row }">
            <StatusDot :status="row.success ? 'online' : 'offline'" />
          </template>
        </el-table-column>
        <el-table-column prop="error_msg" label="错误信息" show-overflow-tooltip />
      </el-table>
      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="historyPage"
          v-model:page-size="historyPageSize"
          :total="historyTotal"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          background
          @size-change="fetchHistory"
          @current-change="fetchHistory"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import StatusDot from '@/components/common/StatusDot.vue'
import UptimeChart from '@/components/charts/UptimeChart.vue'
import LatencyChart from '@/components/charts/LatencyChart.vue'
import { getMonitor, toggleMonitor } from '@/api/monitor'
import { listAlertHistory } from '@/api/alert'

const route = useRoute()
const monitorId = computed(() => Number(route.params.id))

interface MonitorData {
  id: number
  name: string
  url: string
  enabled: boolean
  last_status_code: number
  last_latency_ms: number
  ssl_expiry_days: number | null
  last_content_matched: boolean
}

const monitor = reactive<MonitorData>({
  id: 0,
  name: '',
  url: '',
  enabled: true,
  last_status_code: 0,
  last_latency_ms: 0,
  ssl_expiry_days: null,
  last_content_matched: false,
})

interface AlertItem {
  id: number
  status: string
  alert_type: string
  detail: string
  triggered_at: string
}

const alerts = ref<AlertItem[]>([])

const history = ref<Record<string, unknown>[]>([])
const historyPage = ref(1)
const historyPageSize = ref(20)
const historyTotal = ref(0)

const handleToggle = async () => {
  try {
    await toggleMonitor(monitorId.value, monitor.enabled)
    ElMessage.success(monitor.enabled ? '已启用' : '已禁用')
  } catch {
    monitor.enabled = !monitor.enabled
  }
}

const fetchHistory = async () => {
  try {
    const res = await getMonitor(monitorId.value)
    // History might be nested or returned separately; adapt as needed
    history.value = res.data?.history || []
    historyTotal.value = res.data?.history_total || 0
  } catch {
    /* ignore */
  }
}

onMounted(async () => {
  // Load monitor details
  try {
    const res = await getMonitor(monitorId.value)
    const d = res.data
    monitor.id = d.id
    monitor.name = d.name
    monitor.url = d.url
    monitor.enabled = d.enabled ?? true
    monitor.last_status_code = d.last_status_code || 0
    monitor.last_latency_ms = d.last_latency_ms || 0
    monitor.ssl_expiry_days = d.ssl_expiry_days ?? null
    monitor.last_content_matched = d.last_content_matched ?? false
  } catch {
    ElMessage.error('加载监控数据失败')
  }

  // Load alerts
  try {
    const res = await listAlertHistory({ monitor_id: monitorId.value, page: 1, page_size: 10 })
    alerts.value = res.data?.list || []
  } catch {
    /* ignore */
  }

  fetchHistory()
})
</script>

<style lang="scss" scoped>
.monitor-detail-view {
  max-width: 1200px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.metrics-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.metric-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  padding: 20px 24px;
}

.metric-label {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
  margin-bottom: 8px;
}

.metric-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);

  span {
    font-size: 14px;
    font-weight: 500;
    color: var(--color-text-secondary);
    margin-left: 4px;
  }
}

.text-danger {
  color: var(--color-danger) !important;
}

.text-warning {
  color: var(--color-warning) !important;
}

.charts-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 24px;
}

.chart-card,
.section-card {
  margin-bottom: 0;

  :deep(.el-card__header) {
    padding: 16px 20px;
    border-bottom: 1px solid var(--color-border);
  }

  :deep(.el-card__body) {
    padding: 20px;
  }
}

.section-card {
  margin-bottom: 24px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.alert-timeline-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.alert-type {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.alert-detail {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
