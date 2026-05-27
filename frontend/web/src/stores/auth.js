import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('trbb_token') || null)
  const user  = ref(JSON.parse(localStorage.getItem('trbb_user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin    = computed(() => user.value?.role === 9)

  function setAuth(t, u) {
    token.value = t
    user.value  = u
    localStorage.setItem('trbb_token', t)
    localStorage.setItem('trbb_user', JSON.stringify(u))
    api.defaults.headers.common['Authorization'] = `Bearer ${t}`
  }

  function clearAuth() {
    token.value = null
    user.value  = null
    localStorage.removeItem('trbb_token')
    localStorage.removeItem('trbb_user')
    delete api.defaults.headers.common['Authorization']
  }

  async function login(email, password) {
    const { data } = await api.post('/auth/login', { email, password })
    setAuth(data.token, data.user)
    return data
  }

  async function register(payload) {
    const { data } = await api.post('/auth/register', payload)
    return data
  }

  async function logout() {
    try { await api.post('/auth/logout') } catch {}
    clearAuth()
  }

  async function fetchProfile() {
    const { data } = await api.get('/me')
    user.value = data
    localStorage.setItem('trbb_user', JSON.stringify(data))
    return data
  }

  // Restore token on startup
  if (token.value) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return { token, user, isLoggedIn, isAdmin, login, register, logout, fetchProfile, setAuth, clearAuth }
})
