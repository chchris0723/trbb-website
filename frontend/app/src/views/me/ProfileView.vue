<template>
  <div class="profile-page">
    <div class="page-header">
      <h1 class="section-title">個人資料</h1>
      <p class="text-gray">完善資料後，報名賽事時可快速帶入</p>
    </div>

    <div v-if="pageLoading" class="loading-box">載入中...</div>

    <form v-else @submit.prevent="handleSave" novalidate>

      <!-- ── 基本帳號 ─────────────────────────────────────── -->
      <div class="form-section">
        <h3 class="form-section-title">帳號資訊 <span class="readonly-badge">不可修改</span></h3>
        <div class="form-row-2">
          <div class="form-group">
            <label>會員 ID</label>
            <input :value="auth.user?.username" disabled class="disabled-input" />
          </div>
          <div class="form-group">
            <label>Email（登入帳號）</label>
            <input :value="auth.user?.email" disabled class="disabled-input" />
          </div>
        </div>
      </div>

      <!-- ── 姓名 ──────────────────────────────────────────── -->
      <div class="form-section">
        <h3 class="form-section-title">姓名資料</h3>
        <div class="form-row-3">
          <div class="form-group">
            <label>暱稱 / 顯示名稱</label>
            <input v-model="form.display_name" placeholder="顯示於社群的名稱" />
          </div>
          <div class="form-group">
            <label>中文姓名 <span class="req">*</span></label>
            <input v-model="form.name_zh" placeholder="真實中文姓名" />
          </div>
          <div class="form-group">
            <label>英文姓名</label>
            <input v-model="form.name_en" placeholder="英文全名（如護照）" />
          </div>
        </div>
      </div>

      <!-- ── 個人資訊 ───────────────────────────────────────── -->
      <div class="form-section">
        <h3 class="form-section-title">個人資訊</h3>
        <div class="form-row-3">
          <div class="form-group">
            <label>性別</label>
            <select v-model="form.gender">
              <option :value="null">請選擇</option>
              <option :value="1">男</option>
              <option :value="2">女</option>
              <option :value="3">其他</option>
            </select>
          </div>
          <div class="form-group">
            <label>出生年月日</label>
            <input v-model="form.birthday" type="date" placeholder="YYYY-MM-DD" />
          </div>
          <div class="form-group">
            <label>手機號碼 <span class="req">*</span></label>
            <input v-model="form.phone" type="tel" placeholder="09xxxxxxxx" />
          </div>
        </div>
        <div class="form-row-3">
          <div class="form-group">
            <label>身份證字號</label>
            <input v-model="form.id_number" placeholder="A123456789" />
          </div>
          <div class="form-group">
            <label>護照號碼</label>
            <input v-model="form.passport_number" placeholder="如有護照請填寫" />
          </div>
          <div class="form-group">
            <label>通訊地址</label>
            <input v-model="form.address" placeholder="縣市 + 完整地址" />
          </div>
        </div>
      </div>

      <!-- ── 偏好設定 ───────────────────────────────────────── -->
      <div class="form-section">
        <h3 class="form-section-title">偏好設定</h3>
        <div class="form-row-2">
          <div class="form-group">
            <label>衣服尺寸</label>
            <select v-model="form.shirt_size">
              <option value="">請選擇</option>
              <option v-for="s in shirtSizes" :key="s" :value="s">{{ s }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>飲食習慣</label>
            <select v-model="form.food_type">
              <option :value="null">請選擇</option>
              <option :value="1">葷食</option>
              <option :value="2">素食</option>
              <option :value="3">全素（純素）</option>
            </select>
          </div>
        </div>
      </div>

      <!-- ── 緊急聯絡人 ─────────────────────────────────────── -->
      <div class="form-section">
        <h3 class="form-section-title">緊急聯絡人 <span class="required-hint">報名賽事必填</span></h3>
        <div class="form-row-3">
          <div class="form-group">
            <label>聯絡人姓名 <span class="req">*</span></label>
            <input v-model="form.emergency_contact" placeholder="緊急聯絡人姓名" />
          </div>
          <div class="form-group">
            <label>聯絡人手機 <span class="req">*</span></label>
            <input v-model="form.emergency_phone" type="tel" placeholder="09xxxxxxxx" />
          </div>
          <div class="form-group">
            <label>與本人關係 <span class="req">*</span></label>
            <select v-model="form.emergency_relation">
              <option value="">請選擇</option>
              <option v-for="r in relations" :key="r" :value="r">{{ r }}</option>
            </select>
          </div>
        </div>
      </div>

      <!-- ── 完整度提示 ─────────────────────────────────────── -->
      <div class="completeness-bar">
        <div class="completeness-label">
          資料完整度
          <span :style="{ color: completenessColor }">{{ completeness }}%</span>
        </div>
        <div class="completeness-track">
          <div class="completeness-fill" :style="{ width: completeness + '%', background: completenessColor }"></div>
        </div>
        <div class="completeness-hint" v-if="completeness < 100">
          還差：{{ missingFields.join('、') }}
        </div>
      </div>

      <!-- ── 錯誤 / 成功 ───────────────────────────────────── -->
      <div v-if="error"   class="form-msg error">{{ error }}</div>
      <div v-if="success" class="form-msg success">✓ {{ success }}</div>

      <!-- ── 提交 ──────────────────────────────────────────── -->
      <div class="form-actions">
        <button type="submit" class="btn btn-primary" :disabled="saving">
          <span v-if="saving" class="spinner"></span>
          {{ saving ? '儲存中...' : '儲存變更' }}
        </button>
        <button type="button" class="btn btn-ghost" @click="resetForm">重設</button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'

const auth   = useAuthStore()
const saving     = ref(false)
const pageLoading = ref(true)
const error   = ref('')
const success = ref('')

const shirtSizes = ['XS', 'S', 'M', 'L', 'XL', '2XL', '3XL']
const relations  = ['配偶', '父親', '母親', '兄弟', '姊妹', '子女', '朋友', '其他']

const form = reactive({
  display_name: '', name_zh: '', name_en: '',
  id_number: '', passport_number: '',
  gender: null, birthday: '', phone: '',
  shirt_size: '', food_type: null, address: '',
  emergency_contact: '', emergency_phone: '', emergency_relation: '',
})

// 資料完整度計算
const requiredForEvent = [
  { key: 'name_zh',            label: '中文姓名' },
  { key: 'phone',              label: '手機號碼' },
  { key: 'id_number',          label: '身份證字號' },
  { key: 'gender',             label: '性別' },
  { key: 'birthday',           label: '出生年月日' },
  { key: 'shirt_size',         label: '衣服尺寸' },
  { key: 'food_type',          label: '飲食習慣' },
  { key: 'address',            label: '通訊地址' },
  { key: 'emergency_contact',  label: '緊急聯絡人姓名' },
  { key: 'emergency_phone',    label: '緊急聯絡人手機' },
  { key: 'emergency_relation', label: '緊急聯絡人關係' },
]
const missingFields = computed(() =>
  requiredForEvent.filter(f => !form[f.key] && form[f.key] !== 0).map(f => f.label)
)
const completeness = computed(() => {
  const filled = requiredForEvent.length - missingFields.value.length
  return Math.round((filled / requiredForEvent.length) * 100)
})
const completenessColor = computed(() => {
  if (completeness.value >= 100) return '#22c55e'
  if (completeness.value >= 60)  return '#f59e0b'
  return '#ef4444'
})

function fillForm(user) {
  Object.keys(form).forEach(k => {
    if (user[k] !== undefined && user[k] !== null) form[k] = user[k]
  })
}

function resetForm() {
  if (auth.user) fillForm(auth.user)
  error.value = ''
  success.value = ''
}

async function handleSave() {
  error.value = ''
  success.value = ''
  if (!form.name_zh) { error.value = '請填寫中文姓名'; return }
  if (!form.phone)   { error.value = '請填寫手機號碼'; return }

  saving.value = true
  try {
    const { data } = await api.put('/me', form)
    // 更新 store
    auth.user = { ...auth.user, ...data.user }
    localStorage.setItem('trbb_user', JSON.stringify(auth.user))
    success.value = '個人資料已成功儲存'
    setTimeout(() => { success.value = '' }, 4000)
  } catch(e) {
    error.value = e.response?.data?.error || '儲存失敗，請稍後再試'
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const { data } = await api.get('/me')
    auth.user = { ...auth.user, ...data }
    fillForm(data)
  } catch {}
  finally {
    pageLoading.value = false
  }
})
</script>

<style scoped>
.profile-page { max-width: 800px; }
.page-header { margin-bottom: 2rem; }
.loading-box { padding: 4rem; text-align: center; color: var(--color-gray-2); }

.form-section {
  background: var(--color-bg-card); border: 1px solid var(--color-border);
  border-radius: 8px; padding: 1.5rem; margin-bottom: 1.25rem;
}
.form-section-title {
  font-family: var(--font-cond); font-size: 0.85rem; font-weight: 700;
  letter-spacing: 0.12em; text-transform: uppercase;
  color: var(--color-gray-2); margin-bottom: 1.25rem;
  display: flex; align-items: center; gap: 0.75rem;
}
.readonly-badge {
  font-size: 0.65rem; padding: 0.15rem 0.5rem; border-radius: 3px;
  background: rgba(107,114,128,0.15); color: var(--color-gray-2);
  letter-spacing: 0.05em;
}
.required-hint {
  font-size: 0.72rem; color: var(--color-primary); font-weight: 400;
  letter-spacing: 0.05em; text-transform: none;
}

.form-row-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.form-row-3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 1rem; }
@media (max-width: 640px) {
  .form-row-2, .form-row-3 { grid-template-columns: 1fr; }
}

.form-group { display: flex; flex-direction: column; gap: 0.35rem; }
.form-group label { font-size: 0.78rem; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: var(--color-gray-1); }
.form-group input, .form-group select { width: 100%; }
.disabled-input { opacity: 0.5; cursor: not-allowed; }
.req { color: var(--color-primary); }

/* Completeness */
.completeness-bar {
  background: var(--color-bg-card); border: 1px solid var(--color-border);
  border-radius: 8px; padding: 1.25rem; margin-bottom: 1.25rem;
}
.completeness-label {
  display: flex; justify-content: space-between;
  font-size: 0.85rem; font-weight: 600; margin-bottom: 0.6rem;
}
.completeness-track {
  height: 6px; background: var(--color-border); border-radius: 3px; overflow: hidden;
}
.completeness-fill {
  height: 100%; border-radius: 3px; transition: width 0.5s ease;
}
.completeness-hint { font-size: 0.78rem; color: var(--color-gray-2); margin-top: 0.6rem; }

.form-msg {
  padding: 0.7rem 1rem; border-radius: 6px; font-size: 0.88rem; margin-bottom: 1rem;
}
.form-msg.error   { background: rgba(239,68,68,0.1);  border: 1px solid rgba(239,68,68,0.3);  color: #fca5a5; }
.form-msg.success { background: rgba(34,197,94,0.1);  border: 1px solid rgba(34,197,94,0.3);  color: #86efac; }

.form-actions { display: flex; gap: 1rem; align-items: center; }
.spinner { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
