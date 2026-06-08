<template>
  <div class="agent-list-view">
    <h1 class="page-title">探测点管理</h1>

    <el-card shadow="never" class="list-card">
      <el-table :data="list" v-loading="loading" stripe style="width: 100%">
        <el-table-column label="状态" width="64" align="center">
          <template #default="{ row }">
            <StatusDot :status="row.status === 'online' ? 'online' : 'offline'" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="140" />
        <el-table-column prop="location" label="位置" width="120" />
        <el-table-column label="最后心跳" width="140">
          <template #default="{ row }">
            {{ formatRelativeTime(row.last_heartbeat_at) }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template #default="{ row }">
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
import { ElMessage, ElMessageBox } from 'element-plus'
import StatusDot from '@/components/common/StatusDot.vue'
import { listAgents, deleteAgent } from '@/api/agent'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const list = ref<Record<string, unknown>[]>([])

const formatRelativeTime = (t: string) => {
  if (!t) return '-'
  const d = new Date(t)
  const now = new Date()
  const diff = Math.floor((now.getTime() - d.getTime()) / 1000)

  if (diff < 60) return `${diff} 秒前`
  if (diff < 3600) return `${Math.floor(diff / 60)} 分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)} 小时前`
  return `${Math.floor(diff / 86400)} 天前`
}

const formatTime = (t: string) => {
  if (!t) return '-'
  const d = new Date(t)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await listAgents({ page: page.value, page_size: pageSize.value })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    /* ignore */
  } finally {
    loading.value = false
  }
}

const handleDelete = async (row: Record<string, unknown>) => {
  try {
    await ElMessageBox.confirm(`确定删除探测点 "${row.name}" 吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
    await deleteAgent(row.id as number)
    ElMessage.success('删除成功')
    fetchList()
  } catch {
    /* cancelled */
  }
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.agent-list-view {
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

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
