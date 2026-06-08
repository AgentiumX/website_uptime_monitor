<template>
  <div class="monitor-list-view">
    <div class="page-header">
      <h1 class="page-title">监控列表</h1>
      <el-button type="primary" @click="router.push('/monitors/create')">
        <el-icon><Plus /></el-icon>
        新增监控
      </el-button>
    </div>

    <el-card shadow="never" class="list-card">
      <div class="search-bar">
        <el-input
          v-model="keyword"
          placeholder="搜索监控名称或 URL..."
          :prefix-icon="Search"
          clearable
          style="width: 320px"
          @keyup.enter="fetchList"
          @clear="fetchList"
        />
      </div>

      <el-table :data="list" v-loading="loading" stripe style="width: 100%">
        <el-table-column label="状态" width="64" align="center">
          <template #default="{ row }">
            <StatusDot :status="getStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="140" />
        <el-table-column prop="url" label="URL" min-width="240" show-overflow-tooltip />
        <el-table-column prop="method" label="方法" width="90" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="methodTagType(row.method)">{{ row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="frequency" label="频率" width="90" align="center">
          <template #default="{ row }">
            {{ row.frequency }}s
          </template>
        </el-table-column>
        <el-table-column label="启用" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.enabled"
              @change="() => handleToggle(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" align="center">
          <template #default="{ row }">
            <el-button size="small" link type="primary" @click="router.push(`/monitors/${row.id}/edit`)">
              编辑
            </el-button>
            <el-button size="small" link type="danger" @click="handleDelete(row)">
              删除
            </el-button>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import StatusDot from '@/components/common/StatusDot.vue'
import { listMonitors, deleteMonitor, toggleMonitor } from '@/api/monitor'

const router = useRouter()
const loading = ref(false)
const keyword = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const list = ref<Record<string, unknown>[]>([])

const getStatus = (row: Record<string, unknown>): 'online' | 'offline' | 'warning' => {
  if (row.last_status === 'up') return 'online'
  if (row.last_status === 'down') return 'offline'
  return 'warning'
}

const methodTagType = (method: string) => {
  const map: Record<string, string> = { GET: '', POST: 'success', PUT: 'warning', DELETE: 'danger' }
  return (map[method] || 'info') as 'success' | 'warning' | 'danger' | 'info' | ''
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await listMonitors({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value || undefined,
    })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    /* ignore */
  } finally {
    loading.value = false
  }
}

const handleToggle = async (row: Record<string, unknown>) => {
  const enabled = row.enabled as boolean
  try {
    await toggleMonitor(row.id as number, enabled)
    ElMessage.success(enabled ? '已启用' : '已禁用')
  } catch {
    row.enabled = !enabled
  }
}

const handleDelete = async (row: Record<string, unknown>) => {
  try {
    await ElMessageBox.confirm(`确定删除监控 "${row.name}" 吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
    await deleteMonitor(row.id as number)
    ElMessage.success('删除成功')
    fetchList()
  } catch {
    /* cancelled */
  }
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.monitor-list-view {
  max-width: 1200px;
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

.list-card {
  :deep(.el-card__body) {
    padding: 20px;
  }
}

.search-bar {
  margin-bottom: 16px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
