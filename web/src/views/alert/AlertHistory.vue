<template>
  <div class="alert-history-view">
    <h1 class="page-title">告警历史</h1>

    <el-card shadow="never" class="list-card">
      <el-form :inline="true" class="filter-form">
        <el-form-item label="监控">
          <el-select
            v-model="filters.monitor_id"
            placeholder="全部监控"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="m in monitors"
              :key="m.id"
              :label="m.name"
              :value="m.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="告警中" value="firing" />
            <el-option label="已恢复" value="resolved" />
          </el-select>
        </el-form-item>

        <el-form-item label="类型">
          <el-select v-model="filters.alert_type" placeholder="全部" clearable style="width: 120px">
            <el-option label="状态码" value="status_code" />
            <el-option label="响应时间" value="latency" />
            <el-option label="内容匹配" value="content_match" />
          </el-select>
        </el-form-item>

        <el-form-item label="时间范围">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="fetchList">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="monitor_name" label="监控名称" min-width="140" />
        <el-table-column prop="agent_name" label="探测点" width="120" />
        <el-table-column prop="alert_type" label="告警类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ alertTypeLabel(row.alert_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="detail" label="详情" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'firing' ? 'danger' : 'success'" size="small">
              {{ row.status === 'firing' ? '告警中' : '已恢复' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="triggered_at" label="触发时间" width="170">
          <template #default="{ row }">
            {{ formatTime(row.triggered_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="resolved_at" label="恢复时间" width="170">
          <template #default="{ row }">
            {{ row.resolved_at ? formatTime(row.resolved_at) : '-' }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          background
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { listAlertHistory } from '@/api/alert'
import { listMonitors } from '@/api/monitor'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const list = ref<Record<string, unknown>[]>([])
const monitors = ref<Record<string, unknown>[]>([])

const filters = reactive({
  monitor_id: null as number | null,
  status: '',
  alert_type: '',
  dateRange: null as [string, string] | null,
})

const alertTypeLabel = (type: string) => {
  const map: Record<string, string> = {
    status_code: '状态码',
    latency: '响应时间',
    content_match: '内容匹配',
  }
  return map[type] || type
}

const formatTime = (t: string) => {
  if (!t) return '-'
  const d = new Date(t)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const fetchList = async () => {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.monitor_id) params.monitor_id = filters.monitor_id
    if (filters.status) params.status = filters.status
    if (filters.alert_type) params.alert_type = filters.alert_type
    if (filters.dateRange) {
      params.start_date = filters.dateRange[0]
      params.end_date = filters.dateRange[1]
    }

    const res = await listAlertHistory(params)
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    /* ignore */
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filters.monitor_id = null
  filters.status = ''
  filters.alert_type = ''
  filters.dateRange = null
  page.value = 1
  fetchList()
}

const fetchMonitors = async () => {
  try {
    const res = await listMonitors({ page: 1, page_size: 1000 })
    monitors.value = res.data?.list || []
  } catch {
    /* ignore */
  }
}

onMounted(() => {
  fetchMonitors()
  fetchList()
})
</script>

<style lang="scss" scoped>
.alert-history-view {
  max-width: 1200px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 24px;
}

.list-card {
  :deep(.el-card__body) {
    padding: 20px;
  }
}

.filter-form {
  margin-bottom: 20px;

  :deep(.el-form-item) {
    margin-bottom: 12px;
    margin-right: 16px;
  }
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
