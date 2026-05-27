<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">儀表板</h1>
      <p class="page-subtitle">歡迎回來，{{ store.admin?.display_name }}。以下是今日摘要。</p>
    </div>

    <!-- Stats -->
    <div class="grid-4 mb-2">
      <div class="card stat-card" v-for="s in stats" :key="s.label">
        <div class="stat-label">{{ s.label }}</div>
        <div class="stat-value" :style="s.color ? `color:${s.color}` : ''">{{ s.value }}</div>
        <div class="stat-change" :class="{ neg: s.trend < 0 }">
          {{ s.trend > 0 ? '▲' : '▼' }} {{ Math.abs(s.trend) }}% 較上月
        </div>
      </div>
    </div>

    <div class="grid-2 mt-4">
      <!-- Recent Registrations -->
      <div class="card">
        <div class="card-header">
          <span style="font-weight:600">最新報名</span>
          <RouterLink to="/events" class="btn btn-ghost btn-sm">查看全部</RouterLink>
        </div>
        <div class="card-body" style="padding:0">
          <table class="table">
            <thead>
              <tr><th>姓名</th><th>賽事</th><th>狀態</th><th>時間</th></tr>
            </thead>
            <tbody>
              <tr v-for="r in mockRegistrations" :key="r.id">
                <td>{{ r.name }}</td>
                <td>{{ r.event }}</td>
                <td><span class="badge" :class="r.badgeClass">{{ r.status }}</span></td>
                <td class="text-gray">{{ r.time }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pending Reviews -->
      <div class="card">
        <div class="card-header">
          <span style="font-weight:600">待審核申請</span>
          <RouterLink to="/membership" class="btn btn-primary btn-sm">審核</RouterLink>
        </div>
        <div class="card-body" style="padding:0">
          <table class="table">
            <thead>
              <tr><th>申請人</th><th>類型</th><th>申請時間</th></tr>
            </thead>
            <tbody>
              <tr v-for="m in mockPending" :key="m.id">
                <td>{{ m.name }}</td>
                <td><span class="badge badge-warning">{{ m.type }}</span></td>
                <td class="text-gray">{{ m.time }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
import { useAdminStore } from '@/stores/admin'

const store = useAdminStore()

const stats = [
  { label: '總會員數',   value: '523',     trend: 5.2,   color: null },
  { label: '本月報名',   value: '87',      trend: 12.1,  color: 'var(--primary)' },
  { label: '本月營收',   value: 'NT$98,400', trend: 8.3, color: 'var(--success)' },
  { label: '待審核申請', value: '12',       trend: -3.1,  color: 'var(--warning)' },
]

const mockRegistrations = [
  { id:1, name:'陳大明', event:'台北鐵人三項', status:'已確認', badgeClass:'badge-success', time:'10分鐘前' },
  { id:2, name:'林小花', event:'淡水半程超鐵', status:'待付款', badgeClass:'badge-warning', time:'25分鐘前' },
  { id:3, name:'王志偉', event:'秋季騎乘挑戰', status:'已確認', badgeClass:'badge-success', time:'1小時前' },
  { id:4, name:'張美玲', event:'台北鐵人三項', status:'已取消', badgeClass:'badge-danger',  time:'2小時前' },
]

const mockPending = [
  { id:1, name:'黃建國', type:'一般會員', time:'今天 09:15' },
  { id:2, name:'吳淑芬', type:'進階會員', time:'今天 08:32' },
  { id:3, name:'李俊宏', type:'一般會員', time:'昨天 22:10' },
]
</script>
