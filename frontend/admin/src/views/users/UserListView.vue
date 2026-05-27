<template>
  <div>
    <div class="page-header flex justify-between items-center">
      <div>
        <h1 class="page-title">會員管理</h1>
        <p class="page-subtitle">查看、審核、管理所有會員帳號</p>
      </div>
      <button class="btn btn-primary" @click="showCreateAdmin = true"
        v-if="store.admin?.role >= 9">
        ＋ 新增管理員
      </button>
    </div>

    <!-- Filters -->
    <div class="card mb-2">
      <div class="card-body" style="padding:1rem">
        <div class="filter-row">
          <input v-model="filters.keyword" placeholder="搜尋姓名 / Email / 手機..."
            @keyup.enter="fetchUsers" style="flex:1;min-width:180px" />
          <select v-model="filters.status" @change="fetchUsers">
            <option value="">全部狀態</option>
            <option value="0">待審核</option>
            <option value="1">已啟用</option>
            <option value="2">已停用</option>
            <option value="3">已拒絕</option>
          </select>
          <select v-model="filters.role" @change="fetchUsers">
            <option value="">全部角色</option>
            <option value="1">一般會員</option>
            <option value="8">管理員</option>
            <option value="9">超級管理員</option>
          </select>
          <button class="btn btn-primary btn-sm" @click="fetchUsers">搜尋</button>
          <button class="btn btn-ghost btn-sm" @click="resetFilters">重設</button>
        </div>
      </div>
    </div>

    <!-- Stats bar -->
    <div class="stats-bar mb-2">
      <div class="stat-chip" v-for="s in quickStats" :key="s.label"
        :class="{ active: filters.status === s.val }" @click="setStatus(s.val)">
        <span class="stat-chip-num">{{ s.count }}</span>
        <span class="stat-chip-label">{{ s.label }}</span>
      </div>
    </div>

    <!-- Table -->
    <div class="card">
      <div class="card-body" style="padding:0">
        <div v-if="loading" class="loading-row">載入中...</div>
        <div v-else-if="!users.length" class="loading-row text-gray">查無資料</div>
        <table v-else class="table">
          <thead>
            <tr>
              <th>會員</th>
              <th>Email / 手機</th>
              <th>角色</th>
              <th>狀態</th>
              <th>申請時間</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id">
              <td>
                <div class="user-cell">
                  <div class="avatar">{{ (u.display_name || u.username)[0] }}</div>
                  <div>
                    <div class="fw-bold">{{ u.display_name }}</div>
                    <div class="text-gray text-xs">@{{ u.username }}</div>
                  </div>
                </div>
              </td>
              <td>
                <div>{{ u.email }}</div>
                <div class="text-gray text-xs">{{ u.phone }}</div>
              </td>
              <td><span class="badge" :class="roleBadge(u.role)">{{ roleLabel(u.role) }}</span></td>
              <td><span class="badge" :class="statusBadge(u.status)">{{ statusLabel(u.status) }}</span></td>
              <td class="text-gray">{{ formatDate(u.created_at) }}</td>
              <td>
                <div class="action-btns">
                  <!-- Approve -->
                  <button v-if="u.status === 0" class="btn btn-sm btn-primary"
                    @click="updateStatus(u, 1)" title="核准">✓ 核准</button>
                  <!-- Reject -->
                  <button v-if="u.status === 0" class="btn btn-sm btn-danger"
                    @click="updateStatus(u, 3)" title="拒絕">✗ 拒絕</button>
                  <!-- Suspend active -->
                  <button v-if="u.status === 1 && u.role < 9" class="btn btn-sm btn-ghost"
                    @click="updateStatus(u, 2)" title="停用">停用</button>
                  <!-- Restore suspended -->
                  <button v-if="u.status === 2" class="btn btn-sm btn-ghost"
                    @click="updateStatus(u, 1)" title="恢復">恢復</button>
                  <!-- Detail -->
                  <button class="btn btn-sm btn-ghost" @click="viewUser(u)">查看</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="pagination" v-if="totalPages > 1">
        <button :disabled="page === 1" @click="goPage(page - 1)" class="btn btn-ghost btn-sm">‹</button>
        <span class="text-gray">{{ page }} / {{ totalPages }}</span>
        <button :disabled="page === totalPages" @click="goPage(page + 1)" class="btn btn-ghost btn-sm">›</button>
      </div>
    </div>

    <!-- Detail Modal -->
    <div class="modal-overlay" v-if="selectedUser" @click.self="selectedUser = null">
      <div class="modal">
        <div class="modal-header">
          <h3>會員詳細資料</h3>
          <button @click="selectedUser = null">✕</button>
        </div>
        <div class="modal-body">
          <div class="detail-grid">
            <div class="detail-row"><span>姓名</span><strong>{{ selectedUser.display_name }}</strong></div>
            <div class="detail-row"><span>會員 ID</span><strong>@{{ selectedUser.username }}</strong></div>
            <div class="detail-row"><span>Email</span><strong>{{ selectedUser.email }}</strong></div>
            <div class="detail-row"><span>手機</span><strong>{{ selectedUser.phone }}</strong></div>
            <div class="detail-row"><span>角色</span>
              <span class="badge" :class="roleBadge(selectedUser.role)">{{ roleLabel(selectedUser.role) }}</span>
            </div>
            <div class="detail-row"><span>狀態</span>
              <span class="badge" :class="statusBadge(selectedUser.status)">{{ statusLabel(selectedUser.status) }}</span>
            </div>
            <div class="detail-row"><span>申請時間</span><strong>{{ formatDate(selectedUser.created_at) }}</strong></div>
          </div>
          <div class="modal-actions" v-if="selectedUser.status === 0">
            <button class="btn btn-primary" @click="updateStatus(selectedUser, 1); selectedUser = null">✓ 核准申請</button>
            <button class="btn btn-danger" @click="updateStatus(selectedUser, 3); selectedUser = null">✗ 拒絕申請</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Admin Modal -->
    <div class="modal-overlay" v-if="showCreateAdmin" @click.self="showCreateAdmin = false">
      <div class="modal">
        <div class="modal-header">
          <h3>新增管理員帳號</h3>
          <button @click="showCreateAdmin = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>帳號 ID</label>
            <input v-model="adminForm.username" placeholder="英數字 3~50 字" />
          </div>
          <div class="form-group">
            <label>姓名</label>
            <input v-model="adminForm.display_name" placeholder="顯示名稱" />
          </div>
          <div class="form-group">
            <label>Email</label>
            <input v-model="adminForm.email" type="email" placeholder="admin@example.com" />
          </div>
          <div class="form-group">
            <label>初始密碼</label>
            <input v-model="adminForm.password" type="password" placeholder="至少 8 字元" />
          </div>
          <div class="form-group">
            <label>角色</label>
            <select v-model="adminForm.role">
              <option :value="8">一般管理員</option>
              <option :value="9">超級管理員</option>
            </select>
          </div>
          <div v-if="adminError" class="auth-error" style="margin-top:0.5rem">{{ adminError }}</div>
          <div class="modal-actions">
            <button class="btn btn-primary" @click="createAdmin" :disabled="adminLoading">
              {{ adminLoading ? '建立中...' : '建立帳號' }}
            </button>
            <button class="btn btn-ghost" @click="showCreateAdmin = false">取消</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useAdminStore } from '@/stores/admin'
import api from '@/services/api'


const store = useAdminStore()

const users       = ref([])
const loading     = ref(false)
const page        = ref(1)
const totalPages  = ref(1)
const total       = ref(0)
const selectedUser = ref(null)
const showCreateAdmin = ref(false)
const adminLoading = ref(false)
const adminError   = ref('')
const adminForm   = reactive({ username:'', display_name:'', email:'', password:'', role: 8 })

const filters = reactive({ keyword: '', status: '', role: '' })

const quickStats = ref([
  { label: '全部',   val: '',  count: 0 },
  { label: '待審核', val: '0', count: 0 },
  { label: '已啟用', val: '1', count: 0 },
  { label: '已停用', val: '2', count: 0 },
])

async function fetchUsers() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: 20 }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status !== '') params.status = filters.status
    if (filters.role !== '') params.role = filters.role

    const { data } = await api.get('/users', { params })
    users.value      = data.users || []
    totalPages.value  = data.pages || 1
    total.value       = data.total || 0

    // Update quickStats counts (fetch separately without filters)
    updateQuickStats()
  } catch(e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function updateQuickStats() {
  try {
    const all = await api.get('/users', { params: { page_size: 1 } })
    quickStats.value[0].count = all.data.total
    for (const s of [0, 1, 2]) {
      const r = await api.get('/users', { params: { page_size: 1, status: s } })
      quickStats.value[s + 1].count = r.data.total
    }
  } catch {}
}

async function updateStatus(user, status) {
  try {
    await api.put(`/users/${user.id}/status`, { status })
    user.status = status
    await fetchUsers()
  } catch(e) {
    alert(e.response?.data?.error || '操作失敗')
  }
}

async function createAdmin() {
  adminError.value = ''
  if (!adminForm.username || !adminForm.display_name || !adminForm.email || !adminForm.password) {
    adminError.value = '請填寫所有欄位'; return
  }
  adminLoading.value = true
  try {
    await api.post('/users/admins', adminForm)
    showCreateAdmin.value = false
    Object.assign(adminForm, { username:'', display_name:'', email:'', password:'', role: 8 })
    await fetchUsers()
  } catch(e) {
    adminError.value = e.response?.data?.error || '建立失敗'
  } finally {
    adminLoading.value = false
  }
}

function viewUser(u) { selectedUser.value = { ...u } }
function resetFilters() { filters.keyword = ''; filters.status = ''; filters.role = ''; page.value = 1; fetchUsers() }
function setStatus(val) { filters.status = val; page.value = 1; fetchUsers() }
function goPage(p) { page.value = p; fetchUsers() }
function formatDate(d) { return d ? new Date(d).toLocaleString('zh-TW', { year:'numeric', month:'2-digit', day:'2-digit', hour:'2-digit', minute:'2-digit' }) : '-' }

function roleLabel(r) { return { 1: '一般會員', 2: '教練', 8: '管理員', 9: '超級管理員' }[r] || '未知' }
function roleBadge(r) { return { 1: 'badge-gray', 2: 'badge-primary', 8: 'badge-warning', 9: 'badge-danger' }[r] || 'badge-gray' }
function statusLabel(s) { return { 0: '待審核', 1: '已啟用', 2: '已停用', 3: '已拒絕' }[s] || '未知' }
function statusBadge(s) { return { 0: 'badge-warning', 1: 'badge-success', 2: 'badge-gray', 3: 'badge-danger' }[s] || 'badge-gray' }

onMounted(fetchUsers)
</script>

<style scoped>
.filter-row { display: flex; gap: 0.75rem; flex-wrap: wrap; align-items: center; }
.filter-row input, .filter-row select { height: 36px; font-size: 0.85rem; }

.stats-bar { display: flex; gap: 0.75rem; flex-wrap: wrap; }
.stat-chip {
  background: var(--bg-card); border: 1px solid var(--border);
  border-radius: 4px; padding: 0.5rem 1rem;
  display: flex; gap: 0.5rem; align-items: center;
  cursor: pointer; transition: all 0.15s;
}
.stat-chip:hover, .stat-chip.active { border-color: var(--primary); }
.stat-chip.active .stat-chip-num { color: var(--primary); }
.stat-chip-num { font-family: var(--font-c); font-size: 1.2rem; font-weight: 700; }
.stat-chip-label { font-size: 0.78rem; color: var(--gray-2); }

.user-cell { display: flex; align-items: center; gap: 0.75rem; }
.avatar {
  width: 36px; height: 36px; border-radius: 50%; flex-shrink: 0;
  background: var(--primary); display: flex; align-items: center; justify-content: center;
  font-weight: 700; font-size: 0.9rem;
}
.fw-bold { font-weight: 600; font-size: 0.9rem; }
.text-xs { font-size: 0.75rem; }

.action-btns { display: flex; gap: 0.4rem; flex-wrap: wrap; }
.btn-danger { background: var(--danger); color: #fff; }
.btn-danger:hover { opacity: 0.85; }

.pagination { display: flex; align-items: center; justify-content: center; gap: 1rem; padding: 1rem; border-top: 1px solid var(--border); }

.loading-row { padding: 3rem; text-align: center; color: var(--gray-2); }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.7); z-index: 100; display: flex; align-items: center; justify-content: center; padding: 1rem; }
.modal { background: var(--bg-card); border: 1px solid var(--border); border-radius: 8px; width: 100%; max-width: 520px; }
.modal-header { display: flex; align-items: center; justify-content: space-between; padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--border); }
.modal-header h3 { font-family: var(--font-c); font-size: 1.1rem; font-weight: 700; }
.modal-header button { background: none; border: none; color: var(--gray-2); font-size: 1.2rem; cursor: pointer; }
.modal-body { padding: 1.5rem; }
.detail-grid { display: flex; flex-direction: column; gap: 0.1rem; margin-bottom: 1.5rem; }
.detail-row { display: flex; justify-content: space-between; padding: 0.6rem 0; border-bottom: 1px solid var(--border); font-size: 0.88rem; }
.detail-row span { color: var(--gray-2); }
.modal-actions { display: flex; gap: 0.75rem; margin-top: 1rem; }

.form-group { margin-bottom: 1rem; }
.form-group label { display: block; font-size: 0.78rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.08em; color: var(--gray-1); margin-bottom: 0.4rem; }
.form-group input, .form-group select { width: 100%; }
.auth-error { background: rgba(239,68,68,0.1); border: 1px solid rgba(239,68,68,0.3); border-radius: 4px; color: #fca5a5; font-size: 0.83rem; padding: 0.5rem 0.75rem; }
</style>
