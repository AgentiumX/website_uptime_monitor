import request from './request'

export const listAlertChannels = (params?: Record<string, unknown>) =>
  request.get('/alert-channels', { params })

export const createAlertChannel = (data: Record<string, unknown>) =>
  request.post('/alert-channels', data)

export const updateAlertChannel = (id: number, data: Record<string, unknown>) =>
  request.put(`/alert-channels/${id}`, data)

export const deleteAlertChannel = (id: number) =>
  request.delete(`/alert-channels/${id}`)

export const testAlertChannel = (id: number) =>
  request.post(`/alert-channels/${id}/test`)

export const listAlertHistory = (params?: Record<string, unknown>) =>
  request.get('/alert-history', { params })
