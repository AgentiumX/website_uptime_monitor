<template>
  <span class="status-dot" :class="statusClass" />
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: 'online' | 'offline' | 'warning'
}>()

const statusClass = computed(() => `status-dot--${props.status}`)
</script>

<style lang="scss" scoped>
.status-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  position: relative;

  &--online {
    background-color: var(--color-success);
    box-shadow: 0 0 6px rgba(52, 199, 89, 0.4);
  }

  &--offline {
    background-color: var(--color-danger);
    box-shadow: 0 0 6px rgba(255, 59, 48, 0.4);
  }

  &--warning {
    background-color: var(--color-warning);
    box-shadow: 0 0 6px rgba(255, 149, 0, 0.4);
  }
}
</style>
