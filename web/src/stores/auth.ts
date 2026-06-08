import { defineStore } from 'pinia'
import { login, logout } from '@/api/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userId: null as number | null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
  },
  actions: {
    async login(username: string, password: string) {
      const res = await login(username, password)
      this.token = res.data.token
      this.userId = res.data.user_id
      localStorage.setItem('token', this.token)
    },
    async logout() {
      try { await logout() } catch { /* ignore */ }
      this.token = ''
      this.userId = null
      localStorage.removeItem('token')
    },
  },
})
