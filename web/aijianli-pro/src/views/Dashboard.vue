<template>
  <div class="dashboard-container">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-content">
        <h1 class="welcome-title">欢迎回来，{{ username }}！</h1>
        <p class="welcome-subtitle">AI简历助手为您提供专业的简历制作和面试管理服务</p>
      </div>
      <div class="welcome-illustration">
        <span class="illustration-icon">🎯</span>
      </div>
    </div>

    <!-- 快速操作区域 -->
    <div class="quick-actions">
      <h2 class="section-title">快速操作</h2>
      <div class="action-cards">
        <router-link to="/resume/create" class="action-card create-resume">
          <div class="card-icon">✨</div>
          <h3>创建简历</h3>
          <p>使用AI助手快速生成专业简历</p>
          <span class="card-arrow">→</span>
        </router-link>

        <router-link to="/resume" class="action-card manage-resume">
          <div class="card-icon">📄</div>
          <h3>简历管理</h3>
          <p>查看、编辑和管理您的所有简历</p>
          <span class="card-arrow">→</span>
        </router-link>

        <router-link to="/interview" class="action-card manage-interview">
          <div class="card-icon">💼</div>
          <h3>面试管理</h3>
          <p>安排和管理您的面试日程</p>
          <span class="card-arrow">→</span>
        </router-link>

        <div class="action-card stats">
          <div class="card-icon">📊</div>
          <h3>数据统计</h3>
          <div class="stats-content">
            <div class="stat-item">
              <span class="stat-number">{{ resumeCount }}</span>
              <span class="stat-label">简历数量</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">{{ interviewCount }}</span>
              <span class="stat-label">面试数量</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 最近活动 -->
    <div class="recent-activity">
      <h2 class="section-title">最近活动</h2>
      <div class="activity-list">
        <div v-if="recentActivities.length === 0" class="no-activity">
          <span class="no-activity-icon">📝</span>
          <p>暂无最近活动</p>
          <p class="no-activity-hint">开始创建您的第一份简历吧！</p>
        </div>
        <div v-else v-for="activity in recentActivities" :key="activity.id" class="activity-item">
          <div class="activity-icon">{{ activity.icon }}</div>
          <div class="activity-content">
            <h4>{{ activity.title }}</h4>
            <p>{{ activity.description }}</p>
            <span class="activity-time">{{ activity.time }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 用户信息
const username = ref(localStorage.getItem('username') || '用户')

// 统计数据
const resumeCount = ref(0)
const interviewCount = ref(0)

// 最近活动
const recentActivities = ref<any[]>([])

// 加载统计数据
const loadStats = () => {
  // 这里可以从API获取真实数据
  resumeCount.value = 2
  interviewCount.value = 1
}

// 加载最近活动
const loadRecentActivities = () => {
  // 这里可以从API获取真实数据
  recentActivities.value = [
    {
      id: 1,
      icon: '📄',
      title: '创建了简历',
      description: '张三的软件工程师简历',
      time: '2小时前'
    },
    {
      id: 2,
      icon: '💼',
      title: '安排了面试',
      description: '与腾讯科技的面试',
      time: '1天前'
    }
  ]
}

onMounted(() => {
  loadStats()
  loadRecentActivities()
})
</script>

<style scoped>
.dashboard-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

/* 欢迎区域 */
.welcome-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  padding: 3rem;
  margin-bottom: 3rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: white;
  box-shadow: 0 20px 40px rgba(102, 126, 234, 0.3);
}

.welcome-content {
  flex: 1;
}

.welcome-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0 0 1rem 0;
  line-height: 1.2;
}

.welcome-subtitle {
  font-size: 1.1rem;
  margin: 0;
  opacity: 0.9;
  line-height: 1.6;
}

.welcome-illustration {
  flex-shrink: 0;
}

.illustration-icon {
  font-size: 6rem;
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.2));
}

/* 快速操作区域 */
.quick-actions {
  margin-bottom: 3rem;
}

.section-title {
  font-size: 1.8rem;
  font-weight: 600;
  color: #1F2937;
  margin: 0 0 2rem 0;
  text-align: center;
}

.action-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
}

.action-card {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  text-decoration: none;
  color: inherit;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  border: 2px solid transparent;
}

.action-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
}

.action-card.create-resume:hover {
  border-color: #10B981;
}

.action-card.manage-resume:hover {
  border-color: #3B82F6;
}

.action-card.manage-interview:hover {
  border-color: #F59E0B;
}

.card-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  display: block;
}

.action-card h3 {
  font-size: 1.3rem;
  font-weight: 600;
  color: #1F2937;
  margin: 0 0 0.75rem 0;
}

.action-card p {
  color: #6B7280;
  margin: 0;
  line-height: 1.6;
}

.card-arrow {
  position: absolute;
  top: 1.5rem;
  right: 1.5rem;
  font-size: 1.5rem;
  color: #9CA3AF;
  transition: all 0.3s ease;
}

.action-card:hover .card-arrow {
  color: #EA0202;
  transform: translateX(4px);
}

/* 统计卡片 */
.action-card.stats {
  background: linear-gradient(135deg, #F3F4F6 0%, #E5E7EB 100%);
}

.stats-content {
  display: flex;
  gap: 2rem;
  margin-top: 1rem;
}

.stat-item {
  text-align: center;
  flex: 1;
}

.stat-number {
  display: block;
  font-size: 2rem;
  font-weight: 700;
  color: #EA0202;
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 0.9rem;
  color: #6B7280;
}

/* 最近活动 */
.recent-activity {
  margin-bottom: 3rem;
}

.activity-list {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.no-activity {
  text-align: center;
  padding: 3rem 2rem;
  color: #9CA3AF;
}

.no-activity-icon {
  font-size: 4rem;
  display: block;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.no-activity p {
  margin: 0.5rem 0;
  font-size: 1.1rem;
}

.no-activity-hint {
  font-size: 0.9rem;
  opacity: 0.7;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1.5rem 0;
  border-bottom: 1px solid #F3F4F6;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.activity-content h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  font-weight: 600;
  color: #1F2937;
}

.activity-content p {
  margin: 0 0 0.5rem 0;
  color: #6B7280;
  font-size: 0.9rem;
}

.activity-time {
  font-size: 0.8rem;
  color: #9CA3AF;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 1rem;
  }

  .welcome-section {
    flex-direction: column;
    text-align: center;
    padding: 2rem;
  }

  .welcome-title {
    font-size: 2rem;
  }

  .illustration-icon {
    font-size: 4rem;
    margin-top: 1rem;
  }

  .action-cards {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }

  .stats-content {
    flex-direction: column;
    gap: 1rem;
  }
}
</style>
