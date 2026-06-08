<template>
  <v-chart class="health-donut" :option="chartOption" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { TooltipComponent, LegendComponent, GraphicComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([PieChart, TooltipComponent, LegendComponent, GraphicComponent, CanvasRenderer])

const props = defineProps<{
  healthy: number
  warning: number
  danger: number
}>()

const total = computed(() => props.healthy + props.warning + props.danger)

const chartOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: '{b}: {c} ({d}%)',
  },
  graphic: {
    type: 'text',
    left: 'center',
    top: 'center',
    style: {
      text: `${total.value}\n总数`,
      textAlign: 'center',
      fill: '#1D1D1F',
      fontSize: 20,
      fontWeight: 700,
    },
  },
  series: [
    {
      type: 'pie',
      radius: ['55%', '75%'],
      avoidLabelOverlap: false,
      label: { show: false },
      emphasis: {
        label: { show: false },
      },
      data: [
        { value: props.healthy, name: '健康', itemStyle: { color: '#34C759' } },
        { value: props.warning, name: '警告', itemStyle: { color: '#FF9500' } },
        { value: props.danger, name: '异常', itemStyle: { color: '#FF3B30' } },
      ],
    },
  ],
}))
</script>

<style scoped>
.health-donut {
  width: 100%;
  height: 280px;
}
</style>
