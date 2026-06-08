<template>
  <v-chart class="uptime-chart" :option="chartOption" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps<{
  range: '24h' | '7d' | '30d'
}>()

const generateMockData = () => {
  const points: Record<string, number> = { '24h': 24, '7d': 7, '30d': 30 }
  const count = points[props.range]
  const labels: string[] = []
  const values: number[] = []

  for (let i = 0; i < count; i++) {
    if (props.range === '24h') {
      labels.push(`${i}:00`)
    } else {
      const d = new Date()
      d.setDate(d.getDate() - (count - 1 - i))
      labels.push(`${d.getMonth() + 1}/${d.getDate()}`)
    }
    values.push(+(99 + Math.random() * 1).toFixed(2))
  }
  return { labels, values }
}

const chartOption = computed(() => {
  const { labels, values } = generateMockData()
  return {
    tooltip: {
      trigger: 'axis',
      formatter: '{b}<br/>可用率: {c}%',
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
      min: 98,
      max: 100,
      axisLabel: {
        color: '#86868B',
        fontSize: 11,
        formatter: '{value}%',
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
              { offset: 0, color: 'rgba(0,113,227,0.15)' },
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
.uptime-chart {
  width: 100%;
  height: 280px;
}
</style>
