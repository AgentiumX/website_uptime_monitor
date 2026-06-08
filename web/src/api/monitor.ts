import request from './request'

export const listMonitors = (params?: Record<string, unknown>) =>
  request.get('/monitors', { params })

export const getMonitor = (id: number) =>
  request.get(`/monitors/${id}`)

export const createMonitor = (data: Record<string, unknown>) =>
  request.post('/monitors', data)

export const updateMonitor = (id: number, data: Record<string, unknown>) =>
  request.put(`/monitors/${id}`, data)

export const deleteMonitor = (id: number) =>
  request.delete(`/monitors/${id}`)

export const toggleMonitor = (id: number, enabled: boolean) =>
  request.patch(`/monitors/${id}/toggle`, { enabled })

export const getMonitorMetrics = (id: number, params?: Record<string, unknown>) =>
  request.get(`/monitors/${id}/metrics`, { params })
