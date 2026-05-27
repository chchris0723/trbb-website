<template>
  <div class="auth-page">
    <div class="auth-bg"><div class="auth-grid"></div></div>
    <div class="auth-box">
      <RouterLink to="/" class="auth-logo">
        <span class="tr">TR</span><span class="bb">BB</span>
      </RouterLink>
      <h1 class="auth-title">會員登入</h1>

      <div class="form-group">
        <label>Email</label>
        <input v-model="form.email" type="email" placeholder="your@email.com" @keyup.enter="handleLogin" />
      </div>
      <div class="form-group">
        <label>密碼</label>
        <div class="input-wrap">
          <input v-model="form.password" :type="showPwd ? 'text' : 'password'"
            placeholder="••••••••" @keyup.enter="handleLogin" />
          <button class="eye-btn" @click="showPwd = !showPwd" type="button">
            {{ showPwd ? '🙈' : '👁' }}
          </button>
        </div>
      </div>

      <div v-if="error" class="auth-error">{{ error }}</div>

      <button class="btn btn-primary auth-submit" @click="handleLogin" :disabled="loading">
        <span v-if="loading" class="spinner"></span>
        {{ loading ? '登入中...' : '登入' }}
      </button>

      <div class="auth-footer">
        還沒有帳號？<RouterLink to="/register">立即註冊</RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route  = useRoute()
const auth   = useAuthStore()

const form    = ref({ email: '', password: '' })
const loading = ref(false)
const error   = ref('')
const showPwd = ref(false)

async function handleLogin() {
  error.value = ''
  if (!form.value.email || !form.value.password) {
    error.value = '請填寫 Email 及密碼'
    return
  }
  loading.value = true
  try {
    await auth.login(form.value.email, form.value.password)
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (e) {
    error.value = e.response?.data?.error || '登入失敗，請稍後再試'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh; display: flex; align-items: center; justify-content: center;
  position: relative; background: var(--color-bg);
}
.auth-bg { position: absolute; inset: 0; }
.auth-grid {
  position: absolute; inset: 0;
  background-image:
    linear-gradient(rgba(229,25,26,0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(229,25,26,0.04) 1px, transparent 1px);
  background-size: 50px 50px;
}
.auth-box {
  position: relative; z-index: 1;
  background: var(--color-bg-card); border: 1px solid var(--color-border);
  border-radius: 8px; padding: 2.5rem; width: 100%; max-width: 420px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.6);
}
.auth-logo { display: block; font-family: var(--font-display); font-size: 2.8rem; text-align: center; margin-bottom: 0.5rem; }
.auth-logo .tr { color: #fff; }
.auth-logo .bb { color: var(--color-primary); }
.auth-title { font-family: var(--font-cond); font-size: 1.1rem; letter-spacing: 0.1em; text-align: center; color: var(--color-gray-2); text-transform: uppercase; margin-bottom: 2rem; }
.form-group { margin-bottom: 1.1rem; }
.form-group label { display: block; font-size: 0.78rem; font-weight: 600; letter-spacing: 0.08em; text-transform: uppercase; color: var(--color-gray-1); margin-bottom: 0.4rem; }
.form-group input { width: 100%; }
.input-wrap { position: relative; }
.input-wrap input { width: 100%; padding-right: 2.5rem; }
.eye-btn { position: absolute; right: 0.6rem; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; font-size: 1rem; opacity: 0.6; }
.eye-btn:hover { opacity: 1; }
.auth-error { background: rgba(229,25,26,0.1); border: 1px solid rgba(229,25,26,0.3); border-radius: 4px; color: #ff6b6b; font-size: 0.85rem; padding: 0.6rem 0.9rem; margin-bottom: 1rem; }
.auth-submit { width: 100%; margin-top: 0.5rem; padding: 0.85rem; font-size: 1rem; display: flex; align-items: center; justify-content: center; gap: 0.5rem; }
.spinner { width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.auth-footer { text-align: center; margin-top: 1.5rem; font-size: 0.88rem; color: var(--color-gray-2); }
.auth-footer a { color: var(--color-primary); font-weight: 600; }
.auth-footer a:hover { text-decoration: underline; }
</style>
