import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(localStorage.getItem('trbb_admin_token') || null)
  const admin  = ref(JSON.parse(localStorage.getItem('trbb_admin_user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)

  function setAuth(t, u) {
    token.value = t
    admin.value = u
    localStorage.setItem('trbb_admin_token', t)
    localStorage.setItem('trbb_admin_user', JSON.stringify(u))
    api.defaults.headers.common['Authorization'] = `Bearer ${t}`
  }

  function clearAuth() {
    token.value = null
    admin.value = null
    localStorage.removeItem('trbb_admin_token')
    localStorage.removeItem('trbb_admin_user')
    delete api.defaults.headers.common['Authorization']
  }

  async function login(email, password) {
    const { data } = await api.post('/auth/login', { email, password })
    setAuth(data.token, data.user)
    return data
  }

  async function logout() {
    clearAuth()
  }

  if (token.value) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return { token, admin, isLoggedIn, login, logout, setAuth, clearAuth }
})
