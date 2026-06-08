<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? '编辑告警渠道' : '新增告警渠道'"
    width="600px"
    @update:model-value="$emit('update:visible', $event)"
    @closed="resetForm"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入渠道名称" />
      </el-form-item>

      <el-form-item label="类型" prop="type">
        <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%">
          <el-option label="钉钉" value="dingtalk" />
          <el-option label="企业微信" value="wechat_work" />
          <el-option label="飞书" value="feishu" />
          <el-option label="自定义 Webhook" value="webhook" />
        </el-select>
      </el-form-item>

      <el-form-item label="Webhook URL" prop="webhook_url">
        <el-input v-model="form.webhook_url" placeholder="https://oapi.dingtalk.com/robot/send?access_token=xxx" />
      </el-form-item>

      <el-form-item label="签名密钥">
        <el-input v-model="form.secret" placeholder="可选，用于签名验证" />
      </el-form-item>

      <el-form-item label="启用">
        <el-switch v-model="form.enabled" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { createAlertChannel, updateAlertChannel } from '@/api/alert'

const props = defineProps<{
  visible: boolean
  channel: Record<string, unknown> | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'saved'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const isEdit = computed(() => !!props.channel)

const form = ref({
  name: '',
  type: '',
  webhook_url: '',
  secret: '',
  enabled: true,
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  webhook_url: [
    { required: true, message: '请输入 Webhook URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的 URL', trigger: 'blur' },
  ],
}

watch(
  () => props.visible,
  (val) => {
    if (val && props.channel) {
      form.value = {
        name: props.channel.name as string,
        type: props.channel.type as string,
        webhook_url: props.channel.webhook_url as string,
        secret: (props.channel.secret as string) || '',
        enabled: props.channel.enabled as boolean,
      }
    } else if (val) {
      resetForm()
    }
  }
)

const resetForm = () => {
  form.value = {
    name: '',
    type: '',
    webhook_url: '',
    secret: '',
    enabled: true,
  }
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        await updateAlertChannel(props.channel!.id as number, form.value)
        ElMessage.success('更新成功')
      } else {
        await createAlertChannel(form.value)
        ElMessage.success('创建成功')
      }
      emit('update:visible', false)
      emit('saved')
    } catch {
      /* error shown by interceptor */
    } finally {
      submitting.value = false
    }
  })
}
</script>
