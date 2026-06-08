<template>
  <div class="alert-channel-list-view">
    <div class="page-header">
      <h1 class="page-title">告警通道</h1>
      <el-button type="primary" @click="openDialog()">
        <el-icon><Plus /></el-icon>
        新增通道
      </el-button>
    </div>

    <el-card shadow="never" class="list-card">
      <el-table :data="list" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="name" label="名称" min-width="140" />
        <el-table-column prop="type" label="类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagType(row.type)" size="small">{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="webhook_url" label="Webhook URL" min-width="280" show-overflow-tooltip />
        <el-table-column label="启用" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center">
          <template #default="{ row }">
            <el-button size="small" link type="primary" @click="openDialog(row)">编辑</el-button>
            <el-button size="small" link type="warning" @click="handleTest(row)" :loading="row._testing">测试</el-button>
            <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
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

    <AlertChannelForm
      v-model:visible="dialogVisible"
      :channel="editingChannel"
      @saved="onSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listAlertChannels,
  deleteAlertChannel,
  testAlertChannel,
} from '@/api/alert'
import AlertChannelForm from './AlertChannelForm.vue'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const list = ref<Record<string, unknown>[]>([])
const dialogVisible = ref(false)
const editingChannel = ref<Record<string, unknown> | null>(null)

const typeTagType = (type: string) => {
  const map: Record<string, string> = {
    dingtalk: '',
    wechat_work: 'success',
    feishu: 'warning',
    webhook: 'info',
  }
  return (map[type] || 'info') as 'success' | 'warning' | 'danger' | 'info' | ''
}

const typeLabel = (type: string) => {
  const map: Record<string, string> = {
    dingtalk: '钉钉',
    wechat_work: '企微',
    feishu: '飞书',
    webhook: 'Webhook',
  }
  return map[type] || type
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await listAlertChannels({ page: page.value, page_size: pageSize.value })
    list.value = (res.data?.list || []).map((c: Record<string, unknown>) => ({
      ...c,
      _testing: false,
    }))
    total.value = res.data?.total || 0
  } catch {
    /* ignore */
  } finally {
    loading.value = false
  }
}

const openDialog = (channel?: Record<string, unknown>) => {
  editingChannel.value = channel || null
  dialogVisible.value = true
}

const onSaved = () => {
  dialogVisible.value = false
  fetchList()
}

const handleTest = async (row: Record<string, unknown>) => {
  row._testing = true
  try {
    await testAlertChannel(row.id as number)
    ElMessage.success('测试消息已发送')
  } catch {
    /* error shown by interceptor */
  } finally {
    row._testing = false
  }
}

const handleDelete = async (row: Record<string, unknown>) => {
  try {
    await ElMessageBox.confirm(`确定删除通道 "${row.name}" 吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
    await deleteAlertChannel(row.id as number)
    ElMessage.success('删除成功')
    fetchList()
  } catch {
    /* cancelled */
  }
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.alert-channel-list-view {
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

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
