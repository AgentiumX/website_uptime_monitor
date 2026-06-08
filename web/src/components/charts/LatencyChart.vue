<template>
  <v-chart class="latency-chart" :option="chartOption" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer])

defineProps<{
  monitorId?: number
}>()

const generateMockData = () => {
  const labels: string[] = []
  const values: number[] = []
  for (let i = 23; i >= 0; i--) {
    labels.push(`${i}:00`)
    values.push(Math.floor(80 + Math.random() * 200))
  }
  return { labels, values }
}

const chartOption = computed(() => {
  const { labels, values } = generateMockData()
  return {
    tooltip: {
      trigger: 'axis',
      formatter: '{b}<br/>延迟: {c} ms',
    },
    grid: {
      top: 16,
      right: 16,
      bottom: 24,
      left: 48,
    },
    xAxis: {
      type: 'category',
      data: labels,
      boundaryGap: false,
      axisLine: { lineStyle: { color: '#E5E5E7' } },
      axisLabel: { color: '#86868B', fontSize: 11 },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#86868B',
        fontSize: 11,
        formatter: '{value} ms',
      },
      splitLine: { lineStyle: { color: '#F0F0F2' } },
    },
    series: [
      {
        type: 'line',
        data: values,
        smooth: true,
        symbol: 'none',
        lineStyle: { color: '#0071E3', width: 2 },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(0,113,227,0.1)' },
              { offset: 1, color: 'rgba(0,113,227,0)' },
            ],
          },
        },
      },
    ],
  }
})
</script>

<style scoped>
.latency-chart {
  width: 100%;
  height: 280px;
}
</style>
