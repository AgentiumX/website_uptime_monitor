<template>
  <div class="stat-card">
    <div class="stat-label">{{ label }}</div>
    <div class="stat-value" :style="{ color: valueColor }">
      {{ value }}<span v-if="suffix" class="stat-suffix">{{ suffix }}</span>
    </div>
    <div v-if="trend" class="stat-trend" :class="trendClass">
      <el-icon><ArrowUp v-if="trendUp" /><ArrowDown v-else /></el-icon>
      {{ trend }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  label: string
  value: string | number
  suffix?: string
  valueColor?: string
  trend?: string
  trendUp?: boolean
}>()

const trendClass = computed(() => {
  if (!props.trend) return ''
  return props.trendUp ? 'trend-up' : 'trend-down'
})
</script>

<style lang="scss" scoped>
.stat-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  padding: 24px;
  transition: box-shadow 0.3s ease;

  &:hover {
    box-shadow: var(--shadow-card-hover);
  }
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.stat-suffix {
  font-size: 16px;
  font-weight: 500;
  margin-left: 2px;
  color: var(--color-text-secondary);
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 500;
  margin-top: 8px;

  &.trend-up {
    color: var(--color-success);
  }

  &.trend-down {
    color: var(--color-danger);
  }
}
</style>
