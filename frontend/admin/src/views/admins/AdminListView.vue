<template>
  <div>
    <div class="page-header flex justify-between items-center">
      <div>
        <h1 class="page-title">管理員列表</h1>
        <p class="page-subtitle">超級管理員唯一且不可刪除；一般管理員由超級管理員新增</p>
      </div>
      <button class="btn btn-primary" @click="openCreate" v-if="store.admin?.role >= 9">
        ＋ 新增管理員
      </button>
    </div>

    <div class="card">
      <div class="card-body" style="padding:0">
        <div v-if="loading" class="loading-row">載入中...</div>
        <div v-else-if="!admins.length" class="loading-row text-gray">尚無管理員資料</div>
        <table v-else class="table">
          <thead>
            <tr>
              <th>帳號</th>
              <th>Email</th>
              <th>角色</th>
              <th>狀態</th>
              <th>建立時間</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in admins" :key="u.id">
              <td>
                <div class="user-cell">
                  <div class="avatar" :class="u.role >= 9 ? 'super' : ''">
                    {{ (u.display_name || u.username)[0] }}
                  </div>
                  <div>
                    <div class="fw-bold">{{ u.display_name }}</div>
                    <div class="text-gray text-xs">@{{ u.username }}</div>
                  </div>
                </div>
              </td>
              <td>{{ u.email }}</td>
              <td>
                <span class="badge" :class="u.role >= 9 ? 'badge-danger' : 'badge-warning'">
                  {{ u.role >= 9 ? '超級管理員' : '一般管理員' }}
                </span>
              </td>
              <td><span class="badge" :class="statusBadge(u.status)">{{ statusLabel(u.status) }}</span></td>
              <td class="text-gray text-xs">{{ fmt(u.created_at) }}</td>
              <td>
                <div class="action-btns">
                  <!-- 不能編輯超級管理員（除非自己是超級） -->
                  <button v-if="store.admin?.role >= 9 || u.role < 9"
                    class="btn btn-sm btn-ghost" @click="openEdit(u)">編輯</button>
                  <!-- 不能刪除超級管理員；不能刪除自己 -->
                  <button v-if="store.admin?.role >= 9 && u.role < 9 && u.id !== store.admin?.id"
                    class="btn btn-sm btn-ghost" style="color:var(--danger)"
                    @click="confirmDelete(u)">刪除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <div class="modal-overlay" v-if="editingAdmin || showCreate" @click.self="closeModal">
      <div class="edit-modal">
        <div class="modal-header">
          <h3>{{ showCreate ? '新增管理員' : '編輯管理員' }}</h3>
          <button @click="closeModal">✕</button>
        </div>
        <div class="modal-body">
          <template v-if="showCreate">
            <div class="form-row">
              <div class="form-group">
                <label>帳號 ID <span class="req">*</span></label>
                <input v-model="editForm.username" placeholder="英數字 3~50" />
              </div>
              <div class="form-group">
                <label>Email <span class="req">*</span></label>
                <input v-model="editForm.email" type="email" />
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>顯示名稱 <span class="req">*</span></label>
                <input v-model="editForm.display_name" />
              </div>
              <div class="form-group">
                <label>角色</label>
                <select v-model.number="editForm.role">
                  <option :value="8">一般管理員</option>
                  <option :value="9" v-if="store.admin?.role >= 9">超級管理員</option>
                </select>
              </div>
            </div>
            <div class="form-group mb-1">
              <label>初始密碼 <span class="req">*</span></label>
              <input v-model="editForm.password" type="password" placeholder="至少 8 字元" />
            </div>
          </template>

          <template v-else>
            <div class="form-row">
              <div class="form-group">
                <label>顯示名稱</label>
                <input v-model="editForm.display_name" />
              </div>
              <div class="form-group">
                <label>帳號 ID <span class="readonly-note">不可修改</span></label>
                <input :value="editingAdmin?.username" disabled class="disabled-input" />
              </div>
            </div>
            <div class="form-group mb-1">
              <label>Email <span class="readonly-note">不可修改</span></label>
              <input :value="editingAdmin?.email" disabled class="disabled-input" />
            </div>
            <div class="form-divider">修改密碼（留空不修改）</div>
            <div class="form-group mb-1">
              <label>新密碼</label>
              <input v-model="editForm.new_password" type="password" placeholder="至少 8 字元" />
            </div>
          </template>

          <div v-if="modalError" class="form-error mt-1">{{ modalError }}</div>
          <div class="modal-footer">
            <button class="btn btn-primary" @click="saveModal" :disabled="modalLoading">
              {{ modalLoading ? '儲存中...' : (showCreate ? '建立' : '儲存') }}
            </button>
            <button class="btn btn-ghost" @click="closeModal">取消</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete confirm -->
    <div class="modal-overlay" v-if="deletingAdmin" @click.self="deletingAdmin=null">
      <div class="confirm-modal">
        <h3>確認刪除管理員？</h3>
        <p class="text-gray">「{{ deletingAdmin?.display_name }}」(@{{ deletingAdmin?.username }}) 將被移除管理員權限。</p>
        <div class="confirm-actions">
          <button class="btn btn-danger" @click="doDelete">確認刪除</button>
          <button class="btn btn-ghost" @click="deletingAdmin=null">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import api from '@/services/api'

const store = useAdminStore()
const admins = ref([])
const loading = ref(false)
const editingAdmin = ref(null)
const showCreate = ref(false)
const deletingAdmin = ref(null)
const modalLoading = ref(false)
const modalError = ref('')

const emptyForm = () => ({
  username: '', email: '', display_name: '', password: '', role: 8, new_password: '',
})
const editForm = reactive(emptyForm())

async function fetchAdmins() {
  loading.value = true
  try {
    const { data } = await api.get('/admins', { params: { page_size: 50 } })
    admins.value = data.users || []
  } catch(e) { console.error(e) }
  finally { loading.value = false }
}

function openCreate() {
  Object.assign(editForm, emptyForm())
  modalError.value = ''
  showCreate.value = true
}

function openEdit(u) {
  Object.assign(editForm, {
    username: u.username, email: u.email,
    display_name: u.display_name || '', role: u.role,
    password: '', new_password: '',
  })
  editingAdmin.value = u
  modalError.value = ''
}

function closeModal() { editingAdmin.value = null; showCreate.value = false }

async function saveModal() {
  modalError.value = ''
  modalLoading.value = true
  try {
    if (showCreate.value) {
      if (!editForm.username || !editForm.email || !editForm.display_name || !editForm.password) {
        modalError.value = '請填寫所有必填欄位'; return
      }
      await api.post('/admins', editForm)
    } else {
      await api.put(`/admins/${editingAdmin.value.id}/profile`, {
        display_name: editForm.display_name,
      })
      if (editForm.new_password) {
        if (editForm.new_password.length < 8) { modalError.value = '密碼至少 8 字元'; return }
        await api.put(`/admins/${editingAdmin.value.id}/password`, { password: editForm.new_password })
      }
    }
    closeModal()
    await fetchAdmins()
  } catch(e) {
    modalError.value = e.response?.data?.error || '操作失敗'
  } finally {
    modalLoading.value = false
  }
}

function confirmDelete(u) { deletingAdmin.value = u }
async function doDelete() {
  try {
    await api.delete(`/admins/${deletingAdmin.value.id}`)
    deletingAdmin.value = null
    await fetchAdmins()
  } catch(e) { alert(e.response?.data?.error || '刪除失敗') }
}

function fmt(d) { return d ? new Date(d).toLocaleDateString('zh-TW', { year:'numeric', month:'2-digit', day:'2-digit' }) : '-' }
function statusLabel(s) { return {0:'待審核',1:'啟用',2:'停用',3:'拒絕'}[s]||'未知' }
function statusBadge(s) { return {0:'badge-warning',1:'badge-success',2:'badge-gray',3:'badge-danger'}[s]||'badge-gray' }

onMounted(fetchAdmins)
</script>

<style scoped>
.user-cell { display:flex; align-items:center; gap:.75rem; }
.avatar { width:36px; height:36px; border-radius:50%; background:var(--primary); display:flex; align-items:center; justify-content:center; font-weight:700; font-size:.9rem; flex-shrink:0; }
.avatar.super { background:linear-gradient(135deg,#e5191a,#ff6b35); }
.fw-bold { font-weight:600; font-size:.9rem; }
.text-xs { font-size:.75rem; }
.action-btns { display:flex; gap:.4rem; }
.loading-row { padding:3rem; text-align:center; color:var(--gray-2); }
.modal-overlay { position:fixed; inset:0; background:rgba(0,0,0,.75); z-index:100; display:flex; align-items:center; justify-content:center; padding:1rem; }
.edit-modal { background:var(--bg-card); border:1px solid var(--border); border-radius:8px; width:100%; max-width:520px; max-height:90vh; overflow-y:auto; }
.modal-header { display:flex; align-items:center; justify-content:space-between; padding:1.25rem 1.5rem; border-bottom:1px solid var(--border); position:sticky; top:0; background:var(--bg-card); }
.modal-header h3 { font-family:var(--font-c); font-size:1.1rem; font-weight:700; }
.modal-header button { background:none; border:none; color:var(--gray-2); font-size:1.2rem; cursor:pointer; }
.modal-body { padding:1.5rem; }
.form-row { display:grid; grid-template-columns:1fr 1fr; gap:.75rem; margin-bottom:.75rem; }
.form-group { display:flex; flex-direction:column; gap:.3rem; }
.form-group.mb-1 { margin-bottom:.75rem; }
.form-group label { font-size:.72rem; font-weight:600; text-transform:uppercase; letter-spacing:.06em; color:var(--gray-1); }
.form-group input, .form-group select { width:100%; }
.disabled-input { opacity:.5; cursor:not-allowed; }
.readonly-note { font-size:.65rem; color:var(--gray-2); font-weight:400; text-transform:none; letter-spacing:0; }
.req { color:var(--primary); }
.form-divider { font-size:.72rem; font-weight:600; text-transform:uppercase; letter-spacing:.1em; color:var(--gray-2); border-top:1px solid var(--border); padding-top:.75rem; margin:1rem 0 .75rem; }
.form-error { background:rgba(239,68,68,.1); border:1px solid rgba(239,68,68,.3); border-radius:4px; color:#fca5a5; font-size:.83rem; padding:.5rem .75rem; }
.modal-footer { display:flex; gap:.75rem; margin-top:1.25rem; padding-top:1rem; border-top:1px solid var(--border); }
.mt-1 { margin-top:.5rem; }
.confirm-modal { background:var(--bg-card); border:1px solid var(--border); border-radius:8px; padding:2rem; max-width:400px; width:100%; }
.confirm-modal h3 { font-size:1.1rem; margin-bottom:.75rem; }
.confirm-actions { display:flex; gap:.75rem; margin-top:1.5rem; }
.btn-danger { background:var(--danger); color:#fff; }
.btn-danger:hover { opacity:.85; }
</style>
