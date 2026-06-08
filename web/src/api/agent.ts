import request from './request'

export const listAgents = (params?: Record<string, unknown>) =>
  request.get('/agents', { params })

export const deleteAgent = (id: number) =>
  request.delete(`/agents/${id}`)
