import axios from 'axios'

const adminApi = axios.create({
  baseURL: import.meta.env.VITE_ADMIN_API_BASE_URL,
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' }
})

adminApi.interceptors.response.use(
  res => res,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('trbb_admin_token')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export const thirdApi = axios.create({
  baseURL: import.meta.env.VITE_THIRD_BASE_URL,
  timeout: 30000,
})

export default adminApi
