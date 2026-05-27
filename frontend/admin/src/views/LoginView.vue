<template>
  <div class="login-page">
    <div class="login-bg"><div class="login-grid"></div></div>
    <div class="login-box">
      <div class="login-logo"><span class="tr">TR</span><span class="bb">BB</span></div>
      <h1 class="login-title">後台管理系統</h1>
      <p class="text-gray" style="margin-bottom:2rem;font-size:0.85rem">管理員專用入口，請輸入帳號密碼</p>
      <div class="form-group">
        <label>Email</label>
        <input v-model="form.email" type="email" placeholder="admin@trbbtw.com" />
      </div>
      <div class="form-group">
        <label>密碼</label>
        <input v-model="form.password" type="password" placeholder="••••••••" />
      </div>
      <button class="btn btn-primary" style="width:100%;margin-top:1.5rem" @click="handleLogin" :disabled="loading">
        {{ loading ? '登入中...' : '登入後台' }}
      </button>
      <p v-if="error" style="color:var(--danger);font-size:0.85rem;margin-top:1rem">{{ error }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'

const router = useRouter()
const store  = useAdminStore()
const loading = ref(false)
const error   = ref('')
const form = ref({ email: '', password: '' })

async function handleLogin() {
  loading.value = true
  error.value   = ''
  try {
    await store.login(form.value.email, form.value.password)
    router.push('/')
  } catch(e) {
    error.value = e.response?.data?.error || '登入失敗，請確認帳號密碼'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh; display: flex; align-items: center; justify-content: center;
  position: relative; background: #000;
}
.login-bg { position: absolute; inset: 0; }
.login-grid {
  position: absolute; inset: 0;
  background-image: linear-gradient(rgba(229,25,26,0.04) 1px, transparent 1px),
                    linear-gradient(90deg, rgba(229,25,26,0.04) 1px, transparent 1px);
  background-size: 40px 40px;
}
.login-box {
  position: relative; z-index: 1;
  background: var(--bg-card); border: 1px solid var(--border);
  border-radius: 8px; padding: 3rem 2.5rem; width: 100%; max-width: 400px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.6);
}
.login-logo { font-family: 'Barlow Condensed', sans-serif; font-size: 3rem; font-weight: 700; margin-bottom: 0.5rem; }
.login-logo .tr { color: #fff; }
.login-logo .bb { color: var(--primary); }
.login-title { font-family: 'Barlow Condensed', sans-serif; font-size: 1.2rem; font-weight: 600; letter-spacing: 0.05em; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; font-size: 0.8rem; font-weight: 600; margin-bottom: 0.4rem; color: var(--gray-1); text-transform: uppercase; letter-spacing: 0.08em; }
.form-group input { width: 100%; }
</style>
