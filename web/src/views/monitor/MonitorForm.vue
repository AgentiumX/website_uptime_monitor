<template>
  <div class="monitor-form-view">
    <div class="page-header">
      <h1 class="page-title">{{ isEdit ? '编辑监控' : '新建监控' }}</h1>
    </div>

    <el-card shadow="never" class="form-card">
      <el-steps :active="step" finish-status="success" align-center style="margin-bottom: 32px">
        <el-step title="基本信息" />
        <el-step title="请求配置" />
        <el-step title="告警规则" />
        <el-step title="通知与探测点" />
      </el-steps>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="140px"
        label-position="right"
        class="step-form"
      >
        <!-- Step 1: Basic Info -->
        <div v-show="step === 0" class="step-content">
          <el-form-item label="任务名称" prop="name">
            <el-input v-model="form.name" placeholder="请输入监控任务名称" />
          </el-form-item>
          <el-form-item label="监控地址" prop="url">
            <el-input v-model="form.url" placeholder="https://example.com" />
          </el-form-item>
        </div>

        <!-- Step 2: Request Config -->
        <div v-show="step === 1" class="step-content">
          <el-form-item label="请求方法">
            <el-radio-group v-model="form.method">
              <el-radio-button value="GET">GET</el-radio-button>
              <el-radio-button value="POST">POST</el-radio-button>
              <el-radio-button value="PUT">PUT</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="HTTP 请求头">
            <el-input
              v-model="form.headers"
              type="textarea"
              :rows="4"
              placeholder='{"Content-Type": "application/json"}'
            />
          </el-form-item>
          <el-form-item label="Cookie">
            <el-input v-model="form.cookie" placeholder="key=value; key2=value2" />
          </el-form-item>
          <el-form-item label="HTTP 验证用户名">
            <el-input v-model="form.auth_user" placeholder="可选" />
          </el-form-item>
          <el-form-item label="HTTP 验证密码">
            <el-input v-model="form.auth_pass" type="password" show-password placeholder="可选" />
          </el-form-item>
          <el-form-item label="代理服务器">
            <el-input v-model="form.proxy" placeholder="http://proxy:8080" />
          </el-form-item>
          <el-form-item label="验证证书">
            <el-switch v-model="form.verify_ssl" />
          </el-form-item>
        </div>

        <!-- Step 3: Alert Rules -->
        <div v-show="step === 2" class="step-content">
          <el-form-item label="匹配类型">
            <el-select v-model="form.match_type" style="width: 240px">
              <el-option label="不匹配" value="none" />
              <el-option label="包含匹配" value="contains" />
              <el-option label="不包含匹配" value="not_contains" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="form.match_type !== 'none'" label="匹配内容">
            <el-input v-model="form.match_content" placeholder="请输入匹配关键字" />
          </el-form-item>
          <el-form-item label="状态码阈值">
            <el-input-number v-model="form.status_threshold" :min="200" :max="599" />
            <span class="form-hint">>= 此值视为告警</span>
          </el-form-item>
          <el-form-item label="响应时间阈值">
            <el-input-number v-model="form.latency_threshold" :min="100" :max="30000" :step="500" />
            <span class="form-hint">ms，>= 此值视为告警</span>
          </el-form-item>
          <el-form-item label="连续失败次数">
            <el-radio-group v-model="form.fail_count">
              <el-radio-button v-for="n in [1, 2, 3, 4, 5]" :key="n" :value="n">{{ n }}</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="监控频率">
            <el-select v-model="form.frequency" style="width: 240px">
              <el-option label="30 秒" :value="30" />
              <el-option label="60 秒" :value="60" />
              <el-option label="5 分钟" :value="300" />
              <el-option label="10 分钟" :value="600" />
              <el-option label="30 分钟" :value="1800" />
            </el-select>
          </el-form-item>
        </div>

        <!-- Step 4: Notification & Agents -->
        <div v-show="step === 3" class="step-content">
          <el-form-item label="告警通道">
            <el-checkbox-group v-model="form.alert_channel_ids">
              <el-checkbox
                v-for="ch in channels"
                :key="ch.id"
                :value="ch.id"
              >
                {{ ch.name }}
              </el-checkbox>
            </el-checkbox-group>
            <span v-if="channels.length === 0" class="form-hint">暂无告警通道，请先在告警管理中创建</span>
          </el-form-item>
          <el-form-item label="探测点">
            <el-checkbox-group v-model="form.agent_ids">
              <el-checkbox
                v-for="ag in agents"
                :key="ag.id"
                :value="ag.id"
              >
                {{ ag.name }} ({{ ag.location }})
              </el-checkbox>
            </el-checkbox-group>
            <span v-if="agents.length === 0" class="form-hint">暂无可用探测点</span>
          </el-form-item>
          <el-form-item label="是否启用">
            <el-switch v-model="form.enabled" />
          </el-form-item>
        </div>
      </el-form>

      <div class="form-actions">
        <el-button v-if="step > 0" @click="step--">上一步</el-button>
        <el-button v-if="step < 3" type="primary" @click="nextStep">下一步</el-button>
        <el-button v-if="step === 3" type="primary" :loading="submitting" @click="handleSubmit">
          {{ isEdit ? '保存修改' : '创建监控' }}
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { getMonitor, createMonitor, updateMonitor } from '@/api/monitor'
import { listAlertChannels } from '@/api/alert'
import { listAgents } from '@/api/agent'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const step = ref(0)

const isEdit = computed(() => !!route.params.id)

interface MonitorForm {
  name: string
  url: string
  method: string
  headers: string
  cookie: string
  auth_user: string
  auth_pass: string
  proxy: string
  verify_ssl: boolean
  match_type: string
  match_content: string
  status_threshold: number
  latency_threshold: number
  fail_count: number
  frequency: number
  alert_channel_ids: number[]
  agent_ids: number[]
  enabled: boolean
}

const form = reactive<MonitorForm>({
  name: '',
  url: '',
  method: 'GET',
  headers: '',
  cookie: '',
  auth_user: '',
  auth_pass: '',
  proxy: '',
  verify_ssl: true,
  match_type: 'none',
  match_content: '',
  status_threshold: 400,
  latency_threshold: 3000,
  fail_count: 1,
  frequency: 60,
  alert_channel_ids: [],
  agent_ids: [],
  enabled: true,
})

const rules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  url: [
    { required: true, message: '请输入监控地址', trigger: 'blur' },
    { type: 'url' as const, message: '请输入有效的 URL 地址', trigger: 'blur' },
  ],
}

interface Channel {
  id: number
  name: string
}

interface Agent {
  id: number
  name: string
  location: string
}

const channels = ref<Channel[]>([])
const agents = ref<Agent[]>([])

const stepFields: string[][] = [
  ['name', 'url'],
  [],
  [],
  [],
]

const nextStep = async () => {
  if (!formRef.value) return
  const fields = stepFields[step.value]
  if (fields.length === 0) {
    step.value++
    return
  }
  try {
    await formRef.value.validateField(fields)
    step.value++
  } catch {
    /* validation failed */
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    // Jump to the step with errors
    step.value = 0
    return
  }

  submitting.value = true
  try {
    const data: Record<string, unknown> = {
      name: form.name,
      url: form.url,
      method: form.method,
      headers: form.headers ? JSON.parse(form.headers) : {},
      cookie: form.cookie,
      auth_user: form.auth_user,
      auth_pass: form.auth_pass,
      proxy: form.proxy,
      verify_ssl: form.verify_ssl,
      match_type: form.match_type,
      match_content: form.match_content,
      status_threshold: form.status_threshold,
      latency_threshold: form.latency_threshold,
      fail_count: form.fail_count,
      frequency: form.frequency,
      alert_channel_ids: form.alert_channel_ids,
      agent_ids: form.agent_ids,
      enabled: form.enabled,
    }

    if (isEdit.value) {
      await updateMonitor(Number(route.params.id), data)
      ElMessage.success('修改成功')
    } else {
      await createMonitor(data)
      ElMessage.success('创建成功')
    }
    router.push('/monitors')
  } catch {
    ElMessage.error(isEdit.value ? '修改失败' : '创建失败')
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  // Load channels and agents
  try {
    const chRes = await listAlertChannels({ page: 1, page_size: 100 })
    channels.value = chRes.data?.list || []
  } catch {
    /* ignore */
  }
  try {
    const agRes = await listAgents({ page: 1, page_size: 100 })
    agents.value = agRes.data?.list || []
  } catch {
    /* ignore */
  }

  // Load existing monitor for edit mode
  if (isEdit.value) {
    try {
      const res = await getMonitor(Number(route.params.id))
      const d = res.data
      form.name = d.name || ''
      form.url = d.url || ''
      form.method = d.method || 'GET'
      form.headers = d.headers ? JSON.stringify(d.headers, null, 2) : ''
      form.cookie = d.cookie || ''
      form.auth_user = d.auth_user || ''
      form.auth_pass = d.auth_pass || ''
      form.proxy = d.proxy || ''
      form.verify_ssl = d.verify_ssl ?? true
      form.match_type = d.match_type || 'none'
      form.match_content = d.match_content || ''
      form.status_threshold = d.status_threshold || 400
      form.latency_threshold = d.latency_threshold || 3000
      form.fail_count = d.fail_count || 1
      form.frequency = d.frequency || 60
      form.alert_channel_ids = d.alert_channel_ids || []
      form.agent_ids = d.agent_ids || []
      form.enabled = d.enabled ?? true
    } catch {
      ElMessage.error('加载监控数据失败')
    }
  }
})
</script>

<style lang="scss" scoped>
.monitor-form-view {
  max-width: 800px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.form-card {
  :deep(.el-card__body) {
    padding: 32px;
  }
}

.step-form {
  max-width: 560px;
  margin: 0 auto;
}

.step-content {
  min-height: 280px;
}

.form-hint {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-left: 8px;
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}
</style>
