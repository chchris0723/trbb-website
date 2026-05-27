<template>
  <div class="membership-page">
    <h1 class="section-title">會員狀態</h1>
    <div class="status-card card">
      <div class="status-icon" :class="`s${auth.user?.status}`">
        {{ statusIcon(auth.user?.status) }}
      </div>
      <div class="status-body">
        <h2 class="status-title" :style="{ color: statusColor(auth.user?.status) }">
          {{ statusTitle(auth.user?.status) }}
        </h2>
        <p class="status-desc">{{ statusDesc(auth.user?.status) }}</p>
        <div class="status-meta">
          <div class="meta-row"><span>會員 ID</span><strong>{{ auth.user?.username }}</strong></div>
          <div class="meta-row"><span>Email</span><strong>{{ auth.user?.email }}</strong></div>
          <div class="meta-row"><span>申請時間</span><strong>{{ formatDate(auth.user?.created_at) }}</strong></div>
        </div>
      </div>
    </div>

    <div class="tips card" v-if="auth.user?.status === 0">
      <h3>⏳ 審核中</h3>
      <ul>
        <li>管理員將於 1~3 個工作天內完成審核</li>
        <li>審核通過後，您將可以報名賽事、使用商城等功能</li>
        <li>如有疑問請聯繫 <a href="mailto:info@trbbtw.com">info@trbbtw.com</a></li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'


const auth = useAuthStore()

function statusIcon(s) { return { 0:'⏳', 1:'✅', 2:'🚫', 3:'❌' }[s] ?? '?' }
function statusTitle(s) { return { 0:'審核中', 1:'正式會員', 2:'已停用', 3:'申請拒絕' }[s] ?? '未知' }
function statusColor(s) { return { 0:'#f59e0b', 1:'#22c55e', 2:'#9ca3af', 3:'#ef4444' }[s] ?? '#fff' }
function statusDesc(s) {
  return {
    0: '您的申請已收到，請耐心等待管理員審核，審核完成後將以 Email 通知您。',
    1: '您已是 TRBB 正式會員，可享有所有會員功能。',
    2: '您的帳號已被停用，如有疑問請聯繫管理員。',
    3: '您的申請已被拒絕，如有疑問請聯繫管理員說明。',
  }[s] ?? ''
}
function formatDate(d) { return d ? d ? new Date(d).toLocaleDateString('zh-TW', { year:'numeric', month:'2-digit', day:'2-digit' }) : '-' : '-' }
</script>

<style scoped>
.status-card { padding: 2rem; display: flex; gap: 1.5rem; align-items: flex-start; margin-bottom: 1.5rem; }
.status-icon { font-size: 3rem; flex-shrink: 0; }
.status-title { font-size: 1.4rem; margin-bottom: 0.5rem; }
.status-desc { color: var(--color-gray-1); font-size: 0.9rem; line-height: 1.7; margin-bottom: 1.25rem; }
.meta-row { display: flex; justify-content: space-between; padding: 0.5rem 0; border-bottom: 1px solid var(--color-border); font-size: 0.88rem; }
.meta-row span { color: var(--color-gray-2); }
.tips { padding: 1.5rem; }
.tips h3 { margin-bottom: 0.75rem; }
.tips ul { list-style: none; display: flex; flex-direction: column; gap: 0.4rem; }
.tips li { font-size: 0.9rem; color: var(--color-gray-1); padding-left: 0.5rem; }
.tips li::before { content: '• '; color: var(--color-primary); }
.tips a { color: var(--color-primary); }
</style>
